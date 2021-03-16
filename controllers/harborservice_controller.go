/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"harbor-operator/config"
	"harbor-operator/controllers/internal/sync"
	"harbor-operator/pkg/syncer"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	harborv1 "harbor-operator/api/v1"
)

// HarborServiceReconciler reconciles a HarborService object
type HarborServiceReconciler struct {
	Client     client.Client
	KubeClient *kubernetes.Clientset
	Log        logr.Logger
	Scheme     *runtime.Scheme
	EventsCli  event.Event
	ConfigInfo *config.ConfigFile
}

var HarborServiceList = make(map[string]*harborv1.HarborService)

// +kubebuilder:rbac:groups=harbor.example.com,resources=harborservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=harbor.example.com,resources=harborservices/status,verbs=get;update;patch

func (r *HarborServiceReconciler) Reconcile(contxt context.Context, req ctrl.Request) (ctrl.Result, error) {

	// your logic here
	ctx := context.Background()
	_ = r.Log.WithValues("harbroservice", req.NamespacedName)

	// 获取当前的 CR
	instance := &harborv1.HarborService{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// CR被删除，执行Harbor清理逻辑
			syncers := r.InitSyncers(req)
			r.delete(syncers)

			delete(HarborServiceList, req.Name)
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// 将list到的cr信息初始化到HarborServiceList中
	HarborServiceList[instance.Name] = instance

	// 执行Harbor集群的部署或更新逻辑
	syncers := r.InitSyncers(req)
	if err = r.sync(syncers); err != nil {
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *HarborServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&harborv1.HarborService{}).
		Complete(r)
}

func (r *HarborServiceReconciler) InitSyncers(req ctrl.Request) []syncer.SyncInterface {
	syncers := []syncer.SyncInterface{
		// database同步控制器
		sync.NewDatabaseSyncer(HarborServiceList[req.Name], r.Client, r.ConfigInfo, r.EventsCli),

		// k8s资源(namespace)同步控制器
		sync.NewK8sNsSyncer(HarborServiceList[req.Name], r.Client, r.KubeClient, r.EventsCli),

		// s3同步控制器
		sync.NewMinioSyncer(HarborServiceList[req.Name], r.Client, r.ConfigInfo, r.EventsCli),

		// k8s存储卷同步器
		sync.NewK8sVolumeSyncer(HarborServiceList[req.Name], r.Client, r.KubeClient, r.EventsCli),

		// harbor实例同步控制器，用于删除harbor实例
		sync.NewInstanceSyncer(HarborServiceList[req.Name], r.Client, r.Scheme, req.Namespace, r.ConfigInfo, r.EventsCli),
	}

	return syncers
}

// 执行同步
func (r *HarborServiceReconciler) sync(syncers []syncer.SyncInterface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s); err != nil {
			continue
			//return err
		}
	}
	return nil
}

// 执行删除操作
func (r *HarborServiceReconciler) delete(syncer []syncer.SyncInterface) error {
	for _, s := range syncer {
		s.Delete(context.TODO())
	}
	return nil
}

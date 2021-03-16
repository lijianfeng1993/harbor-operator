package kubernetes

import (
	"context"
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/pkg/syncer"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes"
)

type NamespaceSyncer struct {
	Name              string
	Namespace         string
	Client            client.Client
	KubeClient        *kubernetes.Clientset
	HarborServiceInfo *v1.HarborService
	EventCli          event.Event
}

func NewK8sNsSyncer(name string, hs *v1.HarborService, c client.Client, kubeClient *kubernetes.Clientset, eventCli event.Event) syncer.SyncInterface {
	return &NamespaceSyncer{
		Name:              name,
		Client:            c,
		KubeClient:        kubeClient,
		HarborServiceInfo: hs,
		EventCli:          eventCli,
	}
}

func (is *NamespaceSyncer) GetName() interface{} {
	return is.Name
}

func (is *NamespaceSyncer) Sync(c context.Context) (ctrl.Result, error) {
	switch is.HarborServiceInfo.Status.Condition.Phase {
	case "":
		//为实例初始化单独的namespace用于隔离每个harbor集群
		installInfo := fmt.Sprintf("Start init kubernetes namespace for harborservice: %s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)
		klog.Info(installInfo)
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "initKubernetesResources", "Start init kubernetes namespace for harborservice")

		// 创建namespace
		err := is.createNamespace()
		if err != nil {
			is.EventCli.NewEventAdd(is.HarborServiceInfo, "initKubernetesResourcesFailed", "Init kubernetes namespace for harborservice failed")
			is.HarborServiceInfo.Status.SetFailedStatus("Init kubernetes namespace for harborservice failed")
			v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (is *NamespaceSyncer) Delete(c context.Context) error {

	err := is.KubeClient.CoreV1().Namespaces().Delete(context.TODO(), is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, metav1.DeleteOptions{})
	if err != nil {
		klog.Error(fmt.Sprintf("fail to remove k8s namespace:%s, witherr:%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, err.Error()))
	}

	klog.Info("Remove kubernetes namespace success")
	return nil
}

func (is *NamespaceSyncer) createNamespace() error {
	namespaceName := is.HarborServiceInfo.Spec.InstanceInfo.InstanceName
	ns := k8sv1.Namespace{}
	ns.Name = namespaceName
	_, err := is.KubeClient.CoreV1().Namespaces().Create(context.TODO(), &ns, metav1.CreateOptions{})
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create k8s namespace:%s, witherr: %s", namespaceName, err.Error()))
		return err
	}
	return nil
}

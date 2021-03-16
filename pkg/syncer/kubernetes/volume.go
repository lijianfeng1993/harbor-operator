package kubernetes

import (
	"context"
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/pkg/syncer"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes"
)

type VolumeSyncer struct {
	Name              string
	Namespace         string
	Client            client.Client
	KubeClient        *kubernetes.Clientset
	HarborServiceInfo *v1.HarborService
	EventCli          event.Event
}

func NewK8sVolumeSyncer(name string, hs *v1.HarborService, c client.Client, kubeClient *kubernetes.Clientset, eventCli event.Event) syncer.SyncInterface {
	return &VolumeSyncer{
		Name:              name,
		Client:            c,
		KubeClient:        kubeClient,
		HarborServiceInfo: hs,
		EventCli:          eventCli,
	}
}

func (is *VolumeSyncer) GetName() interface{} {
	return is.Name
}

func (is *VolumeSyncer) Sync(c context.Context) (ctrl.Result, error) {
	switch is.HarborServiceInfo.Status.Condition.Phase {
	case "":
		// 为实例初始化单独的pv，pvc用于存储相关数据
		// pvc没有对接第三方存储，直接用的hostpath，只用于存储jobservice组件的日志
		installInfo := fmt.Sprintf("Start init kubernetes pv, pvc for harborservice: %s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)
		klog.Info(installInfo)
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "initKubernetesResources", "Start init kubernetes  pv, pvc for harborservice")

		// 创建pv
		err := is.createPV()
		if err != nil {
			is.EventCli.NewEventAdd(is.HarborServiceInfo, "initKubernetesResourcesFailed", "Init kubernetes pv for harborservice failed")
			is.HarborServiceInfo.Status.SetFailedStatus("Init kubernetes pv for harborservice failed")
			v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
			return ctrl.Result{}, err
		}

		// 创建pvc
		err = is.createPVC()
		if err != nil {
			is.EventCli.NewEventAdd(is.HarborServiceInfo, "initKubernetesResourcesFailed", "Init kubernetes pvc for harborservice failed")
			is.HarborServiceInfo.Status.SetFailedStatus("Init kubernetes pvc for harborservice failed")
			v1.FlushInstanceStatus(is.Client, is.HarborServiceInfo)
			return ctrl.Result{}, err
		}

		klog.Info("Init k8s namespace pv pvc successd.")
		is.EventCli.NewEventAdd(is.HarborServiceInfo, "initKubernetesResourcesFailedSuccess", "Init kubernetes resource for harborservice success")
	}

	return ctrl.Result{}, nil
}

func (is *VolumeSyncer) Delete(c context.Context) error {
	// 只需删除pv, pvc会被namespace一起删除
	pvName := fmt.Sprintf("pv-%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)

	err := is.KubeClient.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{})
	if err != nil {
		klog.Error(fmt.Sprintf("fail to remove k8s pv:%s, witherr:%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName, err.Error()), "fail to remove pv")
	}
	return nil
}

func (is *VolumeSyncer) createPV() error {
	namespaceName := is.HarborServiceInfo.Spec.InstanceInfo.InstanceName
	pvName := fmt.Sprintf("pv-%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)

	capacity := make(map[k8sv1.ResourceName]resource.Quantity, 0)
	capacity["storage"] = resource.MustParse("1Gi")

	persistentVolumeSource := k8sv1.PersistentVolumeSource{
		HostPath: &k8sv1.HostPathVolumeSource{
			Path: fmt.Sprintf("/apps/logs/%s/jobservice", namespaceName),
		},
	}

	pv := k8sv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvName,
			Namespace: namespaceName,
		},
		Spec: k8sv1.PersistentVolumeSpec{
			StorageClassName:              namespaceName,
			AccessModes:                   []k8sv1.PersistentVolumeAccessMode{k8sv1.ReadWriteMany},
			Capacity:                      capacity,
			PersistentVolumeReclaimPolicy: "Recycle",
			PersistentVolumeSource:        persistentVolumeSource,
		},
	}

	_, err := is.KubeClient.CoreV1().PersistentVolumes().Create(context.TODO(), &pv, metav1.CreateOptions{})
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create k8s pv:%s, witherr: %s", pvName, err.Error()))
		return err
	}
	return nil
}

func (is *VolumeSyncer) createPVC() error {
	namespaceName := is.HarborServiceInfo.Spec.InstanceInfo.InstanceName
	pvcName := fmt.Sprintf("pvc-%s", is.HarborServiceInfo.Spec.InstanceInfo.InstanceName)

	pvc := k8sv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: namespaceName,
		},
		Spec: k8sv1.PersistentVolumeClaimSpec{
			StorageClassName: &namespaceName,
			AccessModes:      []k8sv1.PersistentVolumeAccessMode{k8sv1.ReadWriteMany},
			Resources: k8sv1.ResourceRequirements{
				Requests: k8sv1.ResourceList{
					k8sv1.ResourceName(k8sv1.ResourceStorage): resource.MustParse("1Gi"),
				},
			},
		},
	}

	_, err := is.KubeClient.CoreV1().PersistentVolumeClaims(namespaceName).Create(context.TODO(), &pvc, metav1.CreateOptions{})
	if err != nil {
		klog.Error(fmt.Errorf("fail to create k8s pvc:%s, witherr: %s", pvcName, err.Error()))
		return err
	}
	return nil
}

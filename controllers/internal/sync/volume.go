package sync

import (
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/pkg/syncer"
	syncer_kubernetes "harbor-operator/pkg/syncer/kubernetes"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewK8sNsSyncer returns a new sync.Interface for reconciling k8s namespace
func NewK8sVolumeSyncer(r *v1.HarborService, c client.Client, kubeclient *kubernetes.Clientset, eventCli event.Event) syncer.SyncInterface {

	return syncer_kubernetes.NewK8sVolumeSyncer(fmt.Sprintf("HarborInstance:%s", r.Spec.InstanceInfo.InstanceName), r, c, kubeclient, eventCli)
}

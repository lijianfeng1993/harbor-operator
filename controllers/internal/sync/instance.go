package sync

import (
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/config"
	"harbor-operator/pkg/syncer"
	"harbor-operator/pkg/syncer/instance"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewInstanceSyncer returns a new sync.Interface for reconciling instance
func NewInstanceSyncer(r *v1.HarborService, c client.Client, scheme *runtime.Scheme, namespaceName string, opConfig *config.ConfigFile, eventCli event.Event) syncer.SyncInterface {

	return instance.NewIntanceSyncer(fmt.Sprintf("HarborInstance:%s", r.Spec.InstanceInfo.InstanceName), namespaceName, r, c, scheme, opConfig, eventCli)
}

package sync

import (
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/config"
	"harbor-operator/pkg/syncer"
	"harbor-operator/pkg/syncer/database"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewDatabaseSyncer returns a new sync.Interface for reconciling pgsql
func NewDatabaseSyncer(r *v1.HarborService, c client.Client, opConfig *config.ConfigFile, eventCli event.Event) syncer.SyncInterface {

	return database.NewDatabaseSyncer(fmt.Sprintf("HarborInstance:%s", r.Spec.InstanceInfo.InstanceName), r, c, opConfig, eventCli)
}

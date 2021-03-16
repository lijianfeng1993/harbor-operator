package sync

import (
	"fmt"
	v1 "harbor-operator/api/v1"
	"harbor-operator/config"
	"harbor-operator/pkg/syncer"
	"harbor-operator/pkg/syncer/kubernetes/event"
	"harbor-operator/pkg/syncer/s3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewMinioSyncer returns a new sync.Interface for reconciling minio bucket
func NewMinioSyncer(r *v1.HarborService, c client.Client, opConfig *config.ConfigFile, eventCli event.Event) syncer.SyncInterface {

	return s3.NewMinioSyncer(fmt.Sprintf("HarborInstance:%s", r.Spec.InstanceInfo.InstanceName), r, c, opConfig, eventCli)
}

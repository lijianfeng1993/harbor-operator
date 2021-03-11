package syncer

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Interface represents a syncer. A syncer persists an object
// (known as subject), into a store (kubernetes apiserver or generic stores)
// and records kubernetes events
type SyncInterface interface {
	// GetObject returns the object for which sync applies
	GetName() interface{}
	// Sync create and update related resource
	Sync(context.Context) (ctrl.Result, error)
	// Delete delete related resource
	Delete(ctx context.Context) error
}

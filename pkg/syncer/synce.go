package syncer

import (
	"context"
)

// Sync mutates the subject of the syncer interface using controller-runtime
// CreateOrUpdate method, when obj is not nil. It takes care of setting owner
// references and recording kubernetes events where appropriate
func Sync(ctx context.Context, syncer SyncInterface) error {
	_, err := syncer.Sync(ctx)
	if err != nil {
		return err
	}
	return nil
}

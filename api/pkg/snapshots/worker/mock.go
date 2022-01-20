package worker

import (
	"context"

	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
)

type inProcessPublisher struct {
	snapshotter snapshotter.Snapshotter
}

func NewSync(snapshotter snapshotter.Snapshotter) Queue {
	return &inProcessPublisher{snapshotter: snapshotter}
}

func (p *inProcessPublisher) Enqueue(_ context.Context, codebaseID, viewID, workspaceID string, paths []string, action snapshots.Action) error {
	_, err := p.snapshotter.Snapshot(codebaseID, workspaceID, action, snapshotter.WithPaths(paths), snapshotter.WithOnView(viewID))
	if err != nil {
		return err
	}
	return nil
}

func (inProcessPublisher) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

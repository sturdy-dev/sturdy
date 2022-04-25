package worker

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
)

type inProcessPublisher struct {
	snapshotter *service_snapshots.Service
}

func NewSync(snapshotter *service_snapshots.Service) Queue {
	return &inProcessPublisher{snapshotter: snapshotter}
}

func (p *inProcessPublisher) Enqueue(_ context.Context, codebaseID codebases.ID, viewID, workspaceID string, action snapshots.Action) error {
	_, err := p.snapshotter.Snapshot(codebaseID, workspaceID, action, service_snapshots.WithOnView(viewID))
	if err != nil {
		return err
	}
	return nil
}

func (inProcessPublisher) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

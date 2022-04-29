package worker

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/users"
	service_users "getsturdy.com/api/pkg/users/service"
)

type inProcessPublisher struct {
	snapshotter  *service_snapshots.Service
	usersService service_users.Service
}

func NewSync(
	snapshotter *service_snapshots.Service,
	usersService service_users.Service,
) Queue {
	return &inProcessPublisher{
		snapshotter:  snapshotter,
		usersService: usersService,
	}
}

func (p *inProcessPublisher) Enqueue(ctx context.Context, codebaseID codebases.ID, viewID, workspaceID string, userID users.ID, action snapshots.Action) error {
	user, err := p.usersService.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	_, err = p.snapshotter.Snapshot(ctx, codebaseID, workspaceID, action, service_snapshots.WithOnView(viewID), service_snapshots.WithUser(user))
	if err != nil {
		return err
	}
	return nil
}

func (inProcessPublisher) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

package graphql

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"

	"github.com/graph-gophers/graphql-go"
)

func (r *WorkspaceRootResolver) UpdateWorkspace(ctx context.Context, args resolvers.UpdateWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.Input.ID))
	if err != nil {
		return nil, errors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, errors.Error(err)
	}

	t := time.Now()
	fields := []db_workspaces.UpdateOption{
		db_workspaces.SetUpdatedAt(&t),
	}
	if args.Input.DraftDescription != nil {
		fields = append(fields, db_workspaces.SetDraftDescription(*args.Input.DraftDescription))
	}
	if args.Input.Name != nil {
		fields = append(fields, db_workspaces.SetName(args.Input.Name))
	}

	if err := r.workspaceWriter.UpdateFields(ctx, ws.ID, fields...); err != nil {
		return nil, errors.Error(err)
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
}

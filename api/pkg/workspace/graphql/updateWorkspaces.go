package graphql

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

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

	if args.Input.DraftDescription != nil {
		ws.DraftDescription = *args.Input.DraftDescription
	}
	if args.Input.Name != nil {
		ws.Name = args.Input.Name
	}

	t := time.Now()
	ws.UpdatedAt = &t

	if err := r.workspaceWriter.Update(ws); err != nil {
		return nil, errors.Error(err)
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
}

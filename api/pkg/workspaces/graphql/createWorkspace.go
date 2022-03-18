package graphql

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/workspaces/service"

	"github.com/graph-gophers/graphql-go"
)

func (r *WorkspaceRootResolver) CreateWorkspace(ctx context.Context, args resolvers.CreateWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	codebaseID := codebases.ID(args.Input.CodebaseID)

	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	cb, err := r.codebaseRepo.Get(codebaseID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Mutually exclusive
	if args.Input.OnTopOfChange != nil && args.Input.OnTopOfChangeWithRevert != nil {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest,
			"onTopOfChange", "can't be set together with onTopOfChangeWithRevert",
			"onTopOfChangeWithRevert", "can't be set together with onTopOfChange",
		)
	}

	// Create request to pass to the old REST API route handler
	req := service.CreateWorkspaceRequest{
		CodebaseID: codebaseID,
		UserID:     userID,
	}
	if args.Input.OnTopOfChange != nil || args.Input.OnTopOfChangeWithRevert != nil {
		var id *graphql.ID
		if args.Input.OnTopOfChange != nil {
			id = args.Input.OnTopOfChange
		} else {
			id = args.Input.OnTopOfChangeWithRevert
		}

		ch, err := r.changeResolver.Change(ctx, resolvers.ChangeArgs{ID: id})
		if err != nil {
			return nil, gqlerrors.Error(err)
		}

		if args.Input.OnTopOfChange != nil {
			id := changes.ID(*args.Input.OnTopOfChange)
			req.BaseChangeID = &id
			req.Name = "On " + ch.Title() // TODO: Use changeService, not resolver
		} else if args.Input.OnTopOfChangeWithRevert != nil {
			id := changes.ID(*args.Input.OnTopOfChangeWithRevert)
			req.BaseChangeID = &id
			req.Revert = true
			req.Name = "Revert " + ch.Title()
		}
	}

	ws, err := r.workspaceService.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
}

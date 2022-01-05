package graphql

import (
	"context"

	"mash/pkg/auth"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/workspace/service"

	"github.com/graph-gophers/graphql-go"
)

func (r *WorkspaceRootResolver) CreateWorkspace(ctx context.Context, args resolvers.CreateWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	codebaseID := string(args.Input.CodebaseID)

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

		ch, err := (*r.changeResolver).Change(ctx, resolvers.ChangeArgs{ID: id})
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		trunkCommitID, err := ch.TrunkCommitID()
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		if trunkCommitID == nil {
			return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "repos without commits are not supported")
		}

		if args.Input.OnTopOfChange != nil {
			req.ChangeID = *trunkCommitID
			req.Name = "On " + ch.Title()
		} else {
			req.RevertChangeID = *trunkCommitID
			req.Name = "Revert " + ch.Title()
		}
	}

	ws, err := r.workspaceService.Create(req)
	if err != nil {
		return nil, err
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
}

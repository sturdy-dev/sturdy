package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type LandRootResovler interface {
	LandWorkspaceChange(context.Context, LandWorkspaceArgs) (WorkspaceResolver, error)
	PushWorkspace(context.Context, PushWorkspaceArgs) (WorkspaceResolver, error)
}

type LandWorkspaceArgs struct {
	Input LandWorkspaceInput
}

type LandWorkspaceInput struct {
	WorkspaceID graphql.ID

	// PatchIDs is deprecated and is not used
	PatchIDs *[]string

	// DiffMaxSize is not on the public API
	// TODO: move this to a more appropriate place
	DiffMaxSize int
}

type PushWorkspaceArgs struct {
	Input PushWorkspaceInput
}

type PushWorkspaceInput struct {
	WorkspaceID                graphql.ID
	LandOnSturdyAndPushTracked *bool
}

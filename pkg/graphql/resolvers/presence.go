package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type PresenceRootResolver interface {
	InternalWorkspacePresence(ctx context.Context, workspaceID string) ([]PresenceResolver, error)

	// Mutations
	ReportWorkspacePresence(ctx context.Context, args ReportWorkspacePresenceArgs) (PresenceResolver, error)

	// Subscriptions
	UpdatedWorkspacePresence(ctx context.Context, args UpdatedWorkspacePresenceArgs) (chan PresenceResolver, error)
}

type PresenceResolver interface {
	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	State() (PresenceState, error)
	LastActiveAt() int32
	Workspace(context.Context) (WorkspaceResolver, error)
}

type PresenceState string

const (
	PresenceStateInvalid PresenceState = "Invalid"
	PresenceStateIdle    PresenceState = "Idle"
	PresenceStateViewing PresenceState = "Viewing"
	PresenceStateCoding  PresenceState = "Coding"
)

type ReportWorkspacePresenceArgs struct {
	Input ReportWorkspacePresenceInput
}

type ReportWorkspacePresenceInput struct {
	WorkspaceID graphql.ID
	State       PresenceState
}

type UpdatedWorkspacePresenceArgs struct {
	WorkspaceID *graphql.ID
}

package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/workspaces"

	"github.com/graph-gophers/graphql-go"
)

type WorkspaceArgs struct {
	ID            graphql.ID
	AllowArchived *bool
}

type WorkspacesArgs struct {
	CodebaseID      graphql.ID
	IncludeArchived *bool
}

type UpdateWorkspaceArgs struct {
	Input UpdateWorkspaceInput
}

type UpdateWorkspaceInput struct {
	ID               graphql.ID
	Name             *string
	DraftDescription *string
}

type ExtractWorkspaceArgs struct {
	Input ExtractWorkspaceInput
}

type ExtractWorkspaceInput struct {
	WorkspaceID graphql.ID
	PatchIDs    []string
}

type UpdatedWorkspaceArgs struct {
	ShortCodebaseID *graphql.ID
	WorkspaceID     *graphql.ID
}

type LandWorkspaceArgs struct {
	Input LandWorkspaceInput
}

type LandWorkspaceInput struct {
	WorkspaceID graphql.ID
	PatchIDs    []string

	// DiffMaxSize is not on the public API
	// TODO: move this to a more appropriate place
	DiffMaxSize int
}

type ArchiveWorkspaceArgs struct {
	ID graphql.ID
}

type UnarchiveWorkspaceArgs struct {
	ID graphql.ID
}

type CreateWorkspaceArgs struct {
	Input CreateWorkspaceInput
}

type CreateWorkspaceInput struct {
	CodebaseID              graphql.ID
	OnTopOfChange           *graphql.ID
	OnTopOfChangeWithRevert *graphql.ID
}

type RemovePatchesArgs struct {
	Input RemovePatchesInput
}

type RemovePatchesInput struct {
	WorkspaceID graphql.ID
	HunkIDs     []string
}

type WorkspaceRootResolver interface {
	// internal
	InternalWorkspace(*workspaces.Workspace) WorkspaceResolver

	Workspace(ctx context.Context, args WorkspaceArgs) (WorkspaceResolver, error)
	Workspaces(ctx context.Context, args WorkspacesArgs) ([]WorkspaceResolver, error)

	// Mutations
	UpdateWorkspace(ctx context.Context, args UpdateWorkspaceArgs) (WorkspaceResolver, error)
	LandWorkspaceChange(ctx context.Context, args LandWorkspaceArgs) (WorkspaceResolver, error)
	ArchiveWorkspace(ctx context.Context, args ArchiveWorkspaceArgs) (WorkspaceResolver, error)
	UnarchiveWorkspace(ctx context.Context, args UnarchiveWorkspaceArgs) (WorkspaceResolver, error)
	CreateWorkspace(ctx context.Context, args CreateWorkspaceArgs) (WorkspaceResolver, error)
	ExtractWorkspace(ctx context.Context, args ExtractWorkspaceArgs) (WorkspaceResolver, error)
	RemovePatches(context.Context, RemovePatchesArgs) (WorkspaceResolver, error)

	// Subscriptions
	UpdatedWorkspace(ctx context.Context, args UpdatedWorkspaceArgs) (<-chan WorkspaceResolver, error)
}

type WorkspaceResolver interface {
	ID() graphql.ID
	Name() string
	Codebase(ctx context.Context) (CodebaseResolver, error)
	Author(context.Context) (AuthorResolver, error)
	CreatedAt() int32
	LastLandedAt() *int32
	ArchivedAt() *int32
	UnarchivedAt() *int32
	UpdatedAt() *int32
	LastActivityAt() int32
	DraftDescription() string
	View(ctx context.Context) (ViewResolver, error)
	Comments() ([]TopCommentResolver, error)
	CommentsCount(context.Context) (int32, error)
	GitHubPullRequest(ctx context.Context) (GitHubPullRequestResolver, error)
	UpToDateWithTrunk(context.Context) (bool, error)
	Conflicts(context.Context) (bool, error)
	HeadChange(ctx context.Context) (ChangeResolver, error)
	Activity(ctx context.Context, args ActivityArgs) ([]ActivityResolver, error)
	Reviews(ctx context.Context) ([]ReviewResolver, error)
	Presence(ctx context.Context) ([]PresenceResolver, error)
	Suggestions(context.Context) ([]SuggestionResolver, error)
	Statuses(context.Context) ([]StatusResolver, error)
	Watchers(context.Context) ([]WorkspaceWatcherResolver, error)
	Suggestion(context.Context) (SuggestionResolver, error)
	SuggestingViews() []ViewResolver
	DiffsCount(context.Context) *int32
}

package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type CodebaseGitHubIntegrationRootResolver interface {
	// Internal APIs
	InternalCodebaseGitHubIntegration(codebaseID graphql.ID) (CodebaseGitHubIntegrationResolver, error)
	InternalGitHubRepositoryByID(id string) (CodebaseGitHubIntegrationResolver, error)

	// Mutations
	UpdateCodebaseGitHubIntegration(ctx context.Context, args UpdateCodebaseGitHubIntegrationArgs) (CodebaseGitHubIntegrationResolver, error)
	CreateWorkspaceFromGitHubBranch(ctx context.Context, args CreateWorkspaceFromGitHubBranchArgs) (WorkspaceResolver, error)
	ImportGitHubPullRequests(ctx context.Context, args ImportGitHubPullRequestsInputArgs) (CodebaseResolver, error)
	RefreshGitHubCodebases(ctx context.Context) ([]CodebaseResolver, error)
}

type CodebaseGitHubIntegrationResolver interface {
	ID() graphql.ID
	Codebase(ctx context.Context) (CodebaseResolver, error)
	Owner() string
	Name() string
	CreatedAt() int32
	UninstalledAt() *int32
	TrackedBranch() *string
	SyncedAt() *int32
	Enabled() bool
	GitHubIsSourceOfTruth() bool
	LastPushErrorMessage() *string
	LastPushAt() *int32
}

type UpdateCodebaseGitHubIntegrationArgs struct {
	Input UpdateCodebaseGitHubIntegrationInput
}

type UpdateCodebaseGitHubIntegrationInput struct {
	ID                    graphql.ID
	Enabled               *bool
	GitHubIsSourceOfTruth *bool
}

type CreateWorkspaceFromGitHubBranchArgs struct {
	Input CreateWorkspaceFromGitHubBranchInput
}

type CreateWorkspaceFromGitHubBranchInput struct {
	CodebaseID graphql.ID
	BranchName string
}

type ImportGitHubPullRequestsInputArgs struct {
	Input ImportGitHubPullRequestsInput
}

type ImportGitHubPullRequestsInput struct {
	CodebaseID graphql.ID
}

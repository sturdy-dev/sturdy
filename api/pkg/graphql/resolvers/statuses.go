package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/statuses"

	"github.com/graph-gophers/graphql-go"
)

type StatusesRootResolver interface {
	// Subscriptions
	UpdatedChangesStatuses(context.Context, UpdatedChangesStatusesArgs) (<-chan ChangeStatusResolver, error)
	UpdatedWorkspacesStatuses(context.Context, UpdatedWorkspacesStatusesArgs) (<-chan WorkspaceStatusResolver, error)
	UpdatedGitHubPullRequestStatuses(context.Context, UpdatedGitHubPullRequestStatusesArgs) (<-chan GitHubPullRequestStatusResolver, error)

	// Internal
	InternalWorkspaceStatuses(context.Context, string) ([]WorkspaceStatusResolver, error)
	InternalChangeStatuses(context.Context, *changes.Change) ([]ChangeStatusResolver, error)
	InternalGitHubPullRequestStatuses(context.Context, *github.PullRequest) ([]GitHubPullRequestStatusResolver, error)

	InternalStatus(*statuses.Status) StatusResolver
}

type UpdatedChangesStatusesArgs struct {
	ChangeIDs []graphql.ID
}

type UpdatedWorkspacesStatusesArgs struct {
	WorkspaceIds []graphql.ID
}

type UpdatedGitHubPullRequestStatusesArgs struct {
	ID graphql.ID
}

type commonStatus interface {
	ID() graphql.ID
	Title() string
	Description() *string
	Type() (StatusType, error)
	Timestamp() int32
	DetailsUrl() *string
}

type StatusResolver interface {
	commonStatus
	ToChangeStatus() (ChangeStatusResolver, bool)
	ToWorkspaceStatus() (WorkspaceStatusResolver, bool)
	ToGitHubPullRequestStatus() (GitHubPullRequestStatusResolver, bool)
}

type WorkspaceStatusResolver interface {
	commonStatus
	Workspace(context.Context) (WorkspaceResolver, error)
	Stale(context.Context) (bool, error)
}

type ChangeStatusResolver interface {
	commonStatus
	Change(context.Context) (ChangeResolver, error)
}

type GitHubPullRequestStatusResolver interface {
	commonStatus
	GitHubPullRequest(context.Context) (GitHubPullRequestResolver, error)
}

type StatusType string

const (
	StatusTypeUndefined StatusType = ""
	StatusTypePending   StatusType = "Pending"
	StatusTypeHealthy   StatusType = "Healthy"
	StatusTypeFailing   StatusType = "Failing"
)

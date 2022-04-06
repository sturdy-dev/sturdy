package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/statuses"

	"github.com/graph-gophers/graphql-go"
)

type StatusesRootResolver interface {
	// Subscriptions
	UpdatedChangesStatuses(context.Context, UpdatedChangesStatusesArgs) (<-chan StatusResolver, error)
	UpdatedGitHubPullRequestStatuses(context.Context, UpdatedGitHubPullRequestStatusesArgs) (<-chan StatusResolver, error)

	// Internal
	InteralStatusesByCodebaseIDAndCommitID(ctx context.Context, codebaseID codebases.ID, commitID string) ([]StatusResolver, error)
	InternalStatus(*statuses.Status) StatusResolver
}

type UpdatedChangesStatusesArgs struct {
	ChangeIDs []graphql.ID
}

type UpdatedGitHubPullRequestStatusesArgs struct {
	ID graphql.ID
}

type StatusResolver interface {
	ID() graphql.ID
	Title() string
	Description() *string
	Type() (StatusType, error)
	Timestamp() int32
	DetailsUrl() *string
	Change(context.Context) (ChangeResolver, error)
	GitHubPullRequest(context.Context) (GitHubPullRequestResolver, error)
}

type StatusType string

const (
	StatusTypeUndefined StatusType = ""
	StatusTypePending   StatusType = "Pending"
	StatusTypeHealthy   StatusType = "Healthy"
	StatusTypeFailing   StatusType = "Failing"
)

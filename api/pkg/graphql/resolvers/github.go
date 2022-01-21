package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type GitHubRootResolver interface {
	GitHubRepositories(ctx context.Context) ([]GitHubRepositoryResolver, error)

	// Mutations
	SetupGitHubRepository(ctx context.Context, args SetupGitHubRepositoryArgs) (CodebaseResolver, error)
}

type GitHubRepositoryResolver interface {
	ID() graphql.ID
	GitHubInstallationID() graphql.ID
	GitHubRepositoryID() graphql.ID
	GitHubOwner() string
	GitHubName() string
	Codebase(ctx context.Context) (CodebaseResolver, error)
}

type SetupGitHubRepositoryArgs struct {
	Input SetupGitHubRepositoryInput
}

type SetupGitHubRepositoryInput struct {
	GitHubInstallationID graphql.ID
	GitHubRepositoryID   graphql.ID
	OrganizationID       graphql.ID
}

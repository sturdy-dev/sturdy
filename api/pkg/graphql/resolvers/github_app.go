package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
)

type GitHubAppRootResolver interface {
	GitHubApp() GitHubApp
}

type GitHubApp interface {
	ID() graphql.ID
	Name() string
	ClientID() string
	Validation() GithubValidationApp
}

type GithubValidationApp interface {
	ID() graphql.ID
	Ok(ctx context.Context) (bool, error)
	MissingPermissions(ctx context.Context) ([]string, error)
	MissingEvents(ctx context.Context) ([]string, error)
}

package resolvers

import "github.com/graph-gophers/graphql-go"

type GitHubAppRootResolver interface {
	GitHubApp() GitHubApp
}

type GitHubApp interface {
	ID() graphql.ID
	Name() string
	ClientID() string
}

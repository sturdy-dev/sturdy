package enterprise

import (
	"mash/pkg/github/config"
	"mash/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type gitHubAppRootResolver struct {
	conf config.GitHubAppConfig
}

func NewGitHubAppRootResolver(conf config.GitHubAppConfig) resolvers.GitHubAppRootResolver {
	return &gitHubAppRootResolver{
		conf: conf,
	}
}

func (r *gitHubAppRootResolver) GitHubApp() resolvers.GitHubApp {
	return &gitHubAppResolver{root: r}
}

type gitHubAppResolver struct {
	root *gitHubAppRootResolver
}

func (r *gitHubAppResolver) ID() graphql.ID {
	return "sturdy"
}

func (r *gitHubAppResolver) Name() string {
	return r.root.conf.GitHubAppName
}

func (r *gitHubAppResolver) ClientID() string {
	return r.root.conf.GitHubAppClientID
}

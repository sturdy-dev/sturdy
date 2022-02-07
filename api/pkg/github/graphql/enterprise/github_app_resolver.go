package enterprise

import (
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type gitHubAppRootResolver struct {
	conf *config.GitHubAppConfig
	meta *config.GitHubAppMetadata
}

func NewGitHubAppRootResolver(conf *config.GitHubAppConfig, meta *config.GitHubAppMetadata) resolvers.GitHubAppRootResolver {
	return &gitHubAppRootResolver{
		conf: conf,
		meta: meta,
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
	return r.root.meta.Slug
}

func (r *gitHubAppResolver) ClientID() string {
	return r.root.conf.ClientID
}

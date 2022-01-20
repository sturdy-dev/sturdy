package oss

import (
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type gitHubAppRootResolver struct{}

func NewGitHubAppRootResolver() resolvers.GitHubAppRootResolver {
	return &gitHubAppRootResolver{}
}

func (r *gitHubAppRootResolver) GitHubApp() resolvers.GitHubApp {
	return nil
}

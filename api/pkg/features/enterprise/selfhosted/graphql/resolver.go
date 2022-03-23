package graphql

import (
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type FeaturesRootResolver struct {
	githubConfig *config.GitHubAppConfig
}

func NewFeaturesRootResolver(githubConfig *config.GitHubAppConfig) resolvers.FeaturesRootResolver {
	return &FeaturesRootResolver{
		githubConfig: githubConfig,
	}
}

func (r *FeaturesRootResolver) isGitHubEnabled() bool {
	return r.githubConfig.ID != 0 &&
		r.githubConfig.ClientID != "" &&
		r.githubConfig.Secret != "" &&
		r.githubConfig.PrivateKeyPath != ""
}

func (r *FeaturesRootResolver) Features() []resolvers.Feature {
	ff := []resolvers.Feature{
		resolvers.FeatureBuildkite,
		resolvers.FeatureRemote,
		resolvers.SelfHostedLicense,
	}
	if r.isGitHubEnabled() {
		ff = append(ff, resolvers.FeatureGitHub)
	} else {
		ff = append(ff, resolvers.FeatureGitHubNotConfigured)
	}
	return ff
}

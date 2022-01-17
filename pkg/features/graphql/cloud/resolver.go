package cloud

import "mash/pkg/graphql/resolvers"

type FeaturesRootResolver struct{}

func NewFeaturesRootResolver() resolvers.FeaturesRootResolver {
	return &FeaturesRootResolver{}
}

func (r *FeaturesRootResolver) Features() []resolvers.Feature {
	return []resolvers.Feature{
		resolvers.FeatureBuildkite,
		resolvers.FeatureGitHub,
		resolvers.FeatureMultiTenancy,
	}
}

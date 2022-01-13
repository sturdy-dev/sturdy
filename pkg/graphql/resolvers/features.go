package resolvers

type FeaturesRootResolver interface {
	Features() []Feature
}

type Feature string

const (
	FeatureBuildkite    Feature = "Buildkite"
	FeatureGitHub       Feature = "GitHub"
	FeatureMultiTenancy Feature = "MultiTenancy"
)

package resolvers

type FeaturesRootResolver interface {
	Features() []Feature
}

type Feature string

const (
	FeatureBuildkite    Feature = "Buildkite"
	FeatureGitHub       Feature = "GitHub"
	FeatureMultiTenancy Feature = "MultiTenancy"

	OrganizationSubscriptions Feature = "OrganizationSubscriptions"
	SelfHostedLicense         Feature = "SelfHostedLicense"

	FeaturePasswordAuth Feature = "PasswordAuth"
	FeatureEmailAuth    Feature = "EmailAuth"
)

package features

type Feature string

const (
	FeatureUndefined Feature = ""
	FeatureGitHub    Feature = "github"
	FeatureBuildkite Feature = "buildkite"
)

package resolvers

type FeaturesRootResolver interface {
	Features() []Feature
}

type Feature string

const (
	FeatureBuildkite Feature = "Buildkite"
	FeatureRemote    Feature = "Remote"

	FeatureGitHub              Feature = "GitHub"              // If the GitHub feature is available and ready to use
	FeatureGitHubNotConfigured Feature = "GitHubNotConfigured" // If the GitHub feature is available, but not ready to use

	FeatureMultiTenancy Feature = "MultiTenancy"

	FeatureOrganizationSubscriptions Feature = "OrganizationSubscriptions"
	SelfHostedLicense                Feature = "SelfHostedLicense"

	FeatureEmails Feature = "Emails"

	FeatureDownloadChanges Feature = "DownloadChanges"
)

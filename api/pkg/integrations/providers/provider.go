package providers

type ProviderType string

const (
	ProviderTypeUndefined ProviderType = ""
	ProviderTypeBuild     ProviderType = "build"
)

type ProviderName string

const (
	ProviderNameUndefined ProviderName = ""
	ProviderNameBuildkite ProviderName = "buildkite"
	ProviderNameGithub    ProviderName = "github"
)

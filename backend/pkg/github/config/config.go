package config

type GitHubAppConfig struct {
	GitHubAppID             int64
	GitHubAppName           string // Name used on installation URLs
	GitHubAppClientID       string
	GitHubAppSecret         string
	GitHubAppPrivateKeyPath string
}

package config

type GitHubAppConfig struct {
	ID int64 `long:"id" description:"GitHub App ID"`
	// Name used on installation URLs
	Name           string `long:"name" description:"GitHub App Name"`
	ClientID       string `long:"client-id" description:"GitHub App Client ID"`
	Secret         string `long:"secret" description:"GitHub App Secret"`
	PrivateKeyPath string `long:"private-key-path" description:"Path to GitHub App Private Key"`
}

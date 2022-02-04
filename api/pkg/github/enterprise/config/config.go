package config

type GitHubAppConfig struct {
	ID int64 `long:"id" description:"GitHub App ID" env:"ID"`
	// Name used on installation URLs
	Name           string `long:"name" description:"GitHub App Name" env:"NAME"`
	ClientID       string `long:"client-id" description:"GitHub App Client ID" env:"CLIENT_ID"`
	Secret         string `long:"secret" description:"GitHub App Secret" env:"SECRET"`
	PrivateKeyPath string `long:"private-key-path" description:"Path to GitHub App Private Key" env:"PRIVATE_KEY_PATH"`
}

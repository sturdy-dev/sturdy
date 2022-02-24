package selfhosted

import (
	"getsturdy.com/api/pkg/analytics/proxy"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/users/avatars/uploader"

	"github.com/jessevdk/go-flags"
)

type Configuration struct {
	configuration.Base

	GitHub    *config.GitHubAppConfig `flags-group:"github-app" namespace:"github-app" env-namespace:"STURDY_GITHUB_APP"`
	Analytics *proxy.Configuration    `flags-group:"analytics" namespace:"analytics"`
	Avatars   *uploader.Configuration `flags-group:"avatars" namespace:"users.avatars"`
}

func New() (Configuration, error) {
	cfg := Configuration{}

	parser := flags.NewParser(&cfg, flags.HelpFlag)
	_, err := parser.Parse()
	return cfg, err
}

package selfhosted

import (
	"getsturdy.com/api/pkg/analytics/proxy"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/github/enterprise/config"

	"github.com/jessevdk/go-flags"
)

type Configuration struct {
	configuration.Base

	GitHub    *config.GitHubAppConfig `flags-group:"github-app" namespace:"github-app"`
	Analytics *proxy.Configuration    `flags-group:"analytics" namespace:"analytics"`
}

func New() (Configuration, error) {
	cfg := Configuration{}

	parser := flags.NewParser(&cfg, flags.HelpFlag)
	_, err := parser.Parse()
	return cfg, err
}

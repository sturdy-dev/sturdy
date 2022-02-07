package cloud

import (
	"getsturdy.com/api/pkg/analytics/enterprise/cloud/posthog"
	"getsturdy.com/api/pkg/configuration"
	emails "getsturdy.com/api/pkg/emails/enterprise/cloud"
	"getsturdy.com/api/pkg/github/enterprise/config"
	queue "getsturdy.com/api/pkg/queue/enterprise/cloud"

	"github.com/jessevdk/go-flags"
)

type Configuration struct {
	configuration.Base

	GitHub    *config.GitHubAppConfig `flags-group:"github-app" namespace:"github-app"`
	Analytics *posthog.Configuration  `flags-group:"analytics" namespace:"analytics"`
	Emails    *emails.Configuration   `flags-group:"emails" namespace:"emails"`
	Queue     *queue.Configuration    `flags-group:"queue" namespace:"queue"`
}

func New() (Configuration, error) {
	cfg := Configuration{}

	parser := flags.NewParser(&cfg, flags.HelpFlag)
	_, err := parser.Parse()
	return cfg, err
}

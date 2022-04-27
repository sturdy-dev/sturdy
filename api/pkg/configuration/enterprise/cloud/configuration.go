package cloud

import (
	"errors"
	"fmt"
	"os"

	posthog "getsturdy.com/api/pkg/analytics/enterprise/cloud/posthog/configuration"
	aws "getsturdy.com/api/pkg/aws/enterprise/cloud/configuration"
	"getsturdy.com/api/pkg/configuration"
	service_change_downloads "getsturdy.com/api/pkg/downloads/enterprise/cloud/service/configuration"
	emails "getsturdy.com/api/pkg/emails/enterprise/cloud/configuration"
	"getsturdy.com/api/pkg/github/enterprise/config"
	queue "getsturdy.com/api/pkg/queue/enterprise/cloud/configuration"

	"github.com/jessevdk/go-flags"
)

type Configuration struct {
	configuration.Base

	AWS              *aws.Configuration                      `flags-group:"aws" namespace:"aws"`
	GitHub           *config.GitHubAppConfig                 `flags-group:"github-app" namespace:"github-app"`
	Analytics        *posthog.Configuration                  `flags-group:"analytics" namespace:"analytics"`
	Emails           *emails.Configuration                   `flags-group:"emails" namespace:"emails"`
	Queue            *queue.Configuration                    `flags-group:"queue" namespace:"queue"`
	ChangesDownloads *service_change_downloads.Configuration `flags-group:"downloads" namespace:"downloads"`
}

func New() (Configuration, error) {
	cfg := Configuration{}

	parser := flags.NewParser(&cfg, flags.HelpFlag)
	var flagsErr *flags.Error
	if _, err := parser.Parse(); errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(0)
		panic("unreachable")
	} else {
		return cfg, err
	}
}

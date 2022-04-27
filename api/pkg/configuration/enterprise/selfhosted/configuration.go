package selfhosted

import (
	"errors"
	"fmt"
	"os"

	proxy "getsturdy.com/api/pkg/analytics/proxy/configuration"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/github/enterprise/config"
	uploader "getsturdy.com/api/pkg/users/avatars/uploader/configuration"

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
	var flagsErr *flags.Error
	if _, err := parser.Parse(); errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(0)
		panic("unreachable")
	} else {
		return cfg, err
	}
}

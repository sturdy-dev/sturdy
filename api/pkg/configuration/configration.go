package configuration

import (
	"errors"
	"fmt"
	"os"

	proxy "getsturdy.com/api/pkg/analytics/proxy/configuration"
	service_ci "getsturdy.com/api/pkg/ci/service/configuration"
	db "getsturdy.com/api/pkg/db/configuration"
	"getsturdy.com/api/pkg/di"
	gitserver "getsturdy.com/api/pkg/gitserver/configuration"
	http "getsturdy.com/api/pkg/http/configuration"
	logger "getsturdy.com/api/pkg/logger/configuration"
	metrics "getsturdy.com/api/pkg/metrics/configuration"
	pprof "getsturdy.com/api/pkg/pprof/configuration"
	uploader "getsturdy.com/api/pkg/users/avatars/uploader/configuration"
	provider "getsturdy.com/api/vcs/provider/configuration"

	"github.com/jessevdk/go-flags"
)

type Base struct {
	di.Out

	Provider *provider.Configuration   `flags-group:"vcs" namespace:"vcs"`
	DB       *db.Configuration         `flags-group:"db" namespace:"db"`
	CI       *service_ci.Configuration `flags-group:"ci" namespace:"ci"`
	HTTP     *http.Configuration       `flags-group:"http" namespace:"http"`
	Git      *gitserver.Configuration  `flags-group:"git" namespace:"git"`
	Pprof    *pprof.Configuration      `flags-group:"pprof" namespace:"pprof"`
	Metrics  *metrics.Configuration    `flags-group:"metrics" namespace:"metrics"`
	Logger   *logger.Configuration     `flags-group:"logger" namespace:"logger"`
}

type Configuration struct {
	Base

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

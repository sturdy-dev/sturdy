package configuration

import (
	"getsturdy.com/api/pkg/analytics/proxy"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/gitserver"
	"getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/metrics"
	"getsturdy.com/api/pkg/pprof"
	"getsturdy.com/api/pkg/users/avatars/uploader"
	"getsturdy.com/api/vcs/provider"

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
	Avatars  *uploader.Configuration   `flags-group:"avatars" namespace:"users.avatars"`
}

type Configuration struct {
	Base

	Analytics *proxy.Configuration `flags-group:"analytics" namespace:"analytics"`
}

func New() (Configuration, error) {
	cfg := Configuration{}

	parser := flags.NewParser(&cfg, flags.HelpFlag)
	_, err := parser.Parse()
	return cfg, err
}

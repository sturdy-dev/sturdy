package configuration

import (
	"getsturdy.com/api/pkg/analytics/proxy"
	"getsturdy.com/api/pkg/aws"
	service_change "getsturdy.com/api/pkg/change/service"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/gitserver"
	"getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/metrics"
	"getsturdy.com/api/pkg/pprof"
	"getsturdy.com/api/vcs/provider"

	"github.com/jessevdk/go-flags"
)

type Base struct {
	di.Out

	AWS      *aws.Configuration            `flags-group:"aws" namespace:"aws"`
	Provider *provider.Configuration       `flags-group:"vcs" namespace:"vcs"`
	DB       *db.Configuration             `flags-group:"db" namespace:"db"`
	CI       *service_ci.Configuration     `flags-group:"ci" namespace:"ci"`
	Change   *service_change.Configuration `flags-group:"changes" namespace:"changes"`
	HTTP     *http.Configuration           `flags-group:"http" namespace:"http"`
	Git      *gitserver.Configuration      `flags-group:"git" namespace:"git"`
	Pprof    *pprof.Configuration          `flags-group:"pprof" namespace:"pprof"`
	Metrics  *metrics.Configuration        `flags-group:"metrics" namespace:"metrics"`
	Logger   *logger.Configuration         `flags-group:"logger" namespace:"logger"`
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

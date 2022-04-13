package configuration

import (
	"os"
	"time"

	proxy "getsturdy.com/api/pkg/analytics/proxy/configuration"
	service_ci "getsturdy.com/api/pkg/ci/service/configuration"
	"getsturdy.com/api/pkg/configuration/flags"
	db "getsturdy.com/api/pkg/db/configuration"
	"getsturdy.com/api/pkg/di"
	gitserver "getsturdy.com/api/pkg/gitserver/configuration"
	http "getsturdy.com/api/pkg/http/configuration"
	"getsturdy.com/api/pkg/internal/sturdytest"
	logger "getsturdy.com/api/pkg/logger/configuration"
	metrics "getsturdy.com/api/pkg/metrics/configuration"
	pprof "getsturdy.com/api/pkg/pprof/configuration"
	uploader "getsturdy.com/api/pkg/users/avatars/uploader/configuration"
	provider "getsturdy.com/api/vcs/provider/configuration"
)

func TestModule(c *di.Container) {
	c.Register(func() (Configuration, error) {
		var lfsAddr flags.Addr
		lfsAddrStr := "localhost:8888"
		if n := os.Getenv("E2E_LFS_HOSTNAME"); n != "" {
			lfsAddrStr = n
		}
		if err := lfsAddr.UnmarshalFlag(lfsAddrStr); err != nil {
			return Configuration{}, err
		}

		var dbURL flags.URL
		if err := dbURL.UnmarshalFlag(sturdytest.PsqlDbSourceForTesting()); err != nil {
			return Configuration{}, err
		}

		var metricsAddr flags.Addr
		if err := metricsAddr.UnmarshalFlag("127.0.0.1:2112"); err != nil {
			return Configuration{}, err
		}

		var pprofAddr flags.Addr
		if err := pprofAddr.UnmarshalFlag("127.0.0.1:6060"); err != nil {
			return Configuration{}, err
		}

		var httpAddr flags.Addr
		if err := httpAddr.UnmarshalFlag("127.0.0.1:3000"); err != nil {
			return Configuration{}, err
		}

		tmpPath, err := os.MkdirTemp("", "sturdy_test")
		if err != nil {
			return Configuration{}, err
		}

		return Configuration{
			Base: Base{
				Provider: &provider.Configuration{
					ReposPath: tmpPath,
					LFS:       &provider.GitLFSConfiguration{Addr: lfsAddr},
				},
				DB: &db.Configuration{
					URL:            dbURL,
					ConnectTimeout: time.Second,
				},
				CI:      &service_ci.Configuration{PublicAPIHostname: "localhost"},
				HTTP:    &http.Configuration{Addr: httpAddr},
				Git:     &gitserver.Configuration{},
				Pprof:   &pprof.Configuration{Addr: pprofAddr},
				Metrics: &metrics.Configuration{Addr: metricsAddr},
				Logger: &logger.Configuration{
					Level: "INFO",
				},
			},

			Analytics: &proxy.Configuration{Disable: true},
			Avatars:   &uploader.Configuration{},
		}, nil
	})
}

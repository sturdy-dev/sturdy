package module

import (
	"os"
	"time"

	"getsturdy.com/api/pkg/analytics/proxy"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/configuration/flags"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/gitserver"
	"getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/internal/sturdytest"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/metrics"
	"getsturdy.com/api/pkg/pprof"
	"getsturdy.com/api/pkg/users/avatars/uploader"
	"getsturdy.com/api/vcs/provider"
)

func TestingModule(c *di.Container) {
	c.Register(func() configuration.Configuration {
		var lfsAddr flags.Addr
		lfsAddrStr := "localhost:8888"
		if n := os.Getenv("E2E_LFS_HOSTNAME"); n != "" {
			lfsAddrStr = n
		}
		if err := lfsAddr.UnmarshalFlag(lfsAddrStr); err != nil {
			panic(err)
		}

		var dbURL flags.URL
		if err := dbURL.UnmarshalFlag(sturdytest.PsqlDbSourceForTesting()); err != nil {
			panic(err)
		}

		var metricsAddr flags.Addr
		if err := metricsAddr.UnmarshalFlag("127.0.0.1:2112"); err != nil {
			panic(err)
		}

		var pprofAddr flags.Addr
		if err := pprofAddr.UnmarshalFlag("127.0.0.1:6060"); err != nil {
			panic(err)
		}

		var httpAddr flags.Addr
		if err := httpAddr.UnmarshalFlag("127.0.0.1:3000"); err != nil {
			panic(err)
		}

		tmpPath, err := os.MkdirTemp("", "sturdy_test")
		if err != nil {
			panic(err)
		}

		return configuration.Configuration{
			Base: configuration.Base{
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
		}
	})
}

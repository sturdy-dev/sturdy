package enterprise

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/github/enterprise/db"
	"getsturdy.com/api/pkg/github/enterprise/service"
	"getsturdy.com/api/pkg/github/enterprise/workers"
)

func Module(c *di.Container) {
	c.Import(workers.Module)
	c.Import(db.Module)
	c.Register(service.New)
	c.Register(func() (client.InstallationClientProvider, client.PersonalClientProvider, client.AppClientProvider) {
		return client.NewInstallationClient, client.NewPersonalClient, client.NewAppClient
	})

	// todo: hack to solve circular import dependency
	iq := new(service.ImporterQueue)
	c.Register(func() *service.ImporterQueue {
		return iq
	})

	type importerHack struct{}
	c.Register(func(wq workers.ImporterQueue) importerHack {
		*iq = wq
		return struct{}{}
	})

	// todo: hack to solve circular import dependency
	cq := new(service.ClonerQueue)
	c.Register(func() *service.ClonerQueue {
		return cq
	})
	type clonerHack struct{}
	c.Register(func(wq *workers.ClonerQueue) clonerHack {
		*cq = wq
		return struct{}{}
	})

	// Get and provide GitHubAppMetadata
	c.Register(func(svc *service.Service, logger *zap.Logger) *config.GitHubAppMetadata {
		log := logger.Named("github-app")
		meta, err := svc.GetAppMetadata()
		switch {
		case errors.Is(err, service.ErrNotSetup):
			log.Info("no app has been configured")
			return nil
		case err != nil:
			log.Fatal("unable to get metadata from github", zap.Error(err))
			return nil
		default:
			log.Info("github app metadata", zap.Any("meta", meta))
			return meta
		}
	})
}

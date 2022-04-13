package graphql

import (
	"errors"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	"getsturdy.com/api/pkg/github/enterprise/service"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	meta_workspaces "getsturdy.com/api/pkg/workspaces/meta"
	"getsturdy.com/api/vcs/executor"
	"go.uber.org/zap"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(configuration.Module)
	c.Import(db_github.Module)
	c.Import(executor.Module)
	c.Import(logger.Module)
	c.Import(meta_workspaces.Module)
	c.Import(db_workspaces.Module)
	c.Import(snapshotter.Module)
	c.Import(github_client.Module)
	c.Import(db_snapshots.Module)
	c.Import(service_auth.Module)
	c.Import(service_codebase.Module)
	c.Import(resolvers.Module)
	c.Register(NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver))
	c.Register(NewGitHubAppRootResolver)
	c.Register(NewCodebaseGitHubIntegrationRootResolver)
	c.Register(NewGitHubRootResolver)

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

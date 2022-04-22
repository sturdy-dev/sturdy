package executor

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/vcs/provider"
	"go.uber.org/zap"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(provider.Module)
	c.Register(func(logger *zap.Logger, repoProvider provider.RepoProvider) Provider {
		return NewProvider(logger, repoProvider, 3)
	})
}

package executor

import (
	"mash/vcs/provider"

	"go.uber.org/zap"
)

type Provider interface {
	New() Executor
}

type executorProvider struct {
	logger       *zap.Logger
	repoProvider provider.RepoProvider

	locks *locker
}

func NewProvider(logger *zap.Logger, repoProvider provider.RepoProvider) Provider {
	return &executorProvider{
		logger:       logger.Named("gitExecutor"),
		repoProvider: repoProvider,
		locks:        newLocker(repoProvider),
	}
}

func (p *executorProvider) New() Executor {
	return newExecutor(p.logger, p.repoProvider, p.locks)
}

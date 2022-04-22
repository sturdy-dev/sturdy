package executor

import (
	"getsturdy.com/api/vcs/provider"

	"go.uber.org/zap"
)

type Provider interface {
	New() Executor
}

type executorProvider struct {
	logger           *zap.Logger
	repoProvider     provider.RepoProvider
	minTmpBufferSize int

	locks *locker
}

func NewProvider(logger *zap.Logger, repoProvider provider.RepoProvider, minTmpBufferSize int) Provider {
	return &executorProvider{
		logger:           logger.Named("gitExecutor"),
		repoProvider:     repoProvider,
		minTmpBufferSize: minTmpBufferSize,
		locks:            newLocker(repoProvider),
	}
}

func (p *executorProvider) New() Executor {
	return newExecutor(p.logger, p.repoProvider, p.locks, p.minTmpBufferSize)
}

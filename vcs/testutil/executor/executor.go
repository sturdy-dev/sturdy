package executor

import (
	"testing"

	"go.uber.org/zap"

	"mash/vcs/executor"
	"mash/vcs/testutil"
)

func TestingExecutorProvider(t *testing.T) executor.Provider {
	return executor.NewProvider(zap.NewNop(), testutil.TestingRepoProvider(t))
}

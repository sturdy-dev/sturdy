package executor

import (
	"testing"

	"go.uber.org/zap"

	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/testutil"
)

func TestingExecutorProvider(t *testing.T) executor.Provider {
	return executor.NewProvider(zap.NewNop(), testutil.TestingRepoProvider(t))
}

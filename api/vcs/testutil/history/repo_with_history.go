package history

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
)

func CreateRepoWithRootCommit(t *testing.T, executorProvider executor.Provider) string {
	codebaseID := uuid.NewString()
	var path string
	err := executorProvider.New().AllowRebasingState().Schedule(func(repoProvider provider.RepoProvider) error {
		path = repoProvider.TrunkPath(codebaseID)
		if _, err := vcs.CreateBareRepoWithRootCommit(path); err != nil {
			return err
		}
		return nil
	}).ExecTrunk(codebaseID, "testutilCreateRepoWithRootCommit")
	assert.NoError(t, err)
	return path
}

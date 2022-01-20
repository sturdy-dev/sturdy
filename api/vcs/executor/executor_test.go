package executor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"
	"getsturdy.com/api/vcs/testutil"
)

func TestExecutor_AllowRebasingState(t *testing.T) {
	exec := NewProvider(zap.NewNop(), testutil.TestingRepoProvider(t))

	// create repo
	err := exec.New().AllowRebasingState().Schedule(func(repoProvider provider.RepoProvider) error {
		path := repoProvider.ViewPath("cb", "vw")
		repo, err := vcs.CreateNonBareRepoWithRootCommit(path, "testtrunk")
		if err != nil {
			return fmt.Errorf("failed to create trunk: %w", err)
		}

		assert.NoError(t, repo.CreateNewBranchOnHEAD("b1"))
		assert.NoError(t, repo.CreateNewBranchOnHEAD("b2"))

		assert.NoError(t, repo.CheckoutBranchWithForce("b1"))
		assert.NoError(t, ioutil.WriteFile(filepath.Join(path, "a.txt"), []byte("foo"), 0o644))
		cb1, err := repo.AddAndCommit("hey-b1")
		assert.NoError(t, err)

		assert.NoError(t, repo.CheckoutBranchWithForce("b2"))
		assert.NoError(t, ioutil.WriteFile(filepath.Join(path, "a.txt"), []byte("bar"), 0o644))
		cb2, err := repo.AddAndCommit("hey-b2")
		assert.NoError(t, err)

		rb, _, err := repo.InitRebaseRaw(cb1, cb2)
		assert.NoError(t, err)
		assert.NotNil(t, rb)
		assert.True(t, repo.IsRebasing())
		return nil
	}).ExecView("cb", "vw", "createView")
	assert.NoError(t, err)

	// the view is not conflicting, try to open again without conflicts allowed
	err = exec.New().Read(func(reader vcs.RepoReader) error {
		return nil
	}).ExecView("cb", "vw", "tryToOpen")
	assert.ErrorIs(t, err, ErrIsRebasing)

	// can open if rebasing is allowed
	err = exec.New().AllowRebasingState().Read(func(reader vcs.RepoReader) error {
		assert.True(t, reader.IsRebasing())
		return nil
	}).ExecView("cb", "vw", "tryToOpenWithAllowed")
	assert.NoError(t, err)
}

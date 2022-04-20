package executor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/codebases"
	vcs_codebases "getsturdy.com/api/pkg/codebases/vcs"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"
	"getsturdy.com/api/vcs/testutil"
)

func syncMapSize(m *sync.Map) int {
	size := 0
	m.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}

func TestExecutor_TemporaryView_must_createtwo(t *testing.T) {
	exec := NewProvider(zap.NewNop(), testutil.TestingRepoProvider(t))

	codebaseID := codebases.ID("cb")

	assert.NoError(t, exec.New().
		AllowRebasingState().
		Schedule(vcs_codebases.Create(codebaseID)).
		ExecTrunk(codebaseID, "createTrunk"), "failed to create trunk")

	recordViewIDs := make(chan struct{})
	viewIDs := &sync.Map{}
	recordViewID := func(reader vcs.RepoReader) error {
		<-recordViewIDs
		viewIDs.Store(*reader.ViewID(), true)
		return nil
	}

	viewOneUsed := make(chan struct{})
	viewTwoUsed := make(chan struct{})
	go func() {
		// this will create a new view, becasuse there are no tmp views yet
		assert.NoError(t, exec.New().Read(recordViewID).ExecTemporaryView(codebaseID, "test1"))
		close(viewOneUsed)
	}()
	go func() {
		// this should create a new view, becasuse there are no tmp views available
		assert.NoError(t, exec.New().Read(recordViewID).ExecTemporaryView(codebaseID, "test2"))
		close(viewTwoUsed)
	}()

	close(recordViewIDs)
	<-viewOneUsed
	<-viewTwoUsed

	// this will reuse one of the existing tmp views
	assert.NoError(t, exec.New().Read(recordViewID).ExecTemporaryView(codebaseID, "test3"))
	assert.Equal(t, 2, syncMapSize(viewIDs), "expexted two temporary views")
}

func TestExecutor_TemporaryView_must_reuse(t *testing.T) {
	exec := NewProvider(zap.NewNop(), testutil.TestingRepoProvider(t))

	codebaseID := codebases.ID("cb")

	assert.NoError(t, exec.New().
		AllowRebasingState().
		Schedule(vcs_codebases.Create(codebaseID)).
		ExecTrunk(codebaseID, "createTrunk"), "failed to create trunk")

	viewIDs := map[string]bool{}
	recordViewID := func(reader vcs.RepoReader) error {
		viewIDs[*reader.ViewID()] = true
		return nil
	}

	assert.NoError(t, exec.New().Read(recordViewID).ExecTemporaryView(codebaseID, "test1"))
	assert.NoError(t, exec.New().Read(recordViewID).ExecTemporaryView(codebaseID, "test2"))
	assert.Len(t, viewIDs, 1, "expected temporary view to be reused")
}

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

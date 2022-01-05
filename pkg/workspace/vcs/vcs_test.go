package vcs

import (
	"io/ioutil"
	codebasevcs "mash/pkg/codebase/vcs"
	"mash/vcs"
	"mash/vcs/testutil"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestCreateWorkspace(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	codebaseID := "codebaseID"
	err := codebasevcs.Create(repoProvider, codebaseID)
	assert.NoError(t, err)

	workspaceID := "workspaceID"
	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = Create(trunkRepo, workspaceID)
	assert.NoError(t, err)
}

func TestListChangesInWorkspace(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	codebaseID := "codebaseID"
	err := codebasevcs.Create(repoProvider, codebaseID)
	assert.NoError(t, err)

	workspaceID := "workspaceID"
	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = Create(trunkRepo, workspaceID)
	assert.NoError(t, err)

	changes, err := ListChanges(repoProvider, codebaseID, workspaceID, 100)
	assert.NoError(t, err)
	assert.Len(t, changes, 0)

	trunkPath := repoProvider.TrunkPath(codebaseID)

	viewPath, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	clonedRepo, err := vcs.CloneRepo(trunkPath, viewPath)
	assert.NoError(t, err)

	// The combination FetchOriginCLI, CreateBranchTrackingUpstream, CheckoutBranchWithForce
	// is common enough and could be a function in vcs
	err = clonedRepo.FetchOriginCLI()
	assert.NoError(t, err)
	err = clonedRepo.CreateBranchTrackingUpstream(workspaceID)
	assert.NoError(t, err)
	err = clonedRepo.CheckoutBranchWithForce(workspaceID)
	assert.NoError(t, err)

	// clone the repo so that we can add some more commits
	_, err = clonedRepo.AddAndCommit("foo0")
	assert.NoError(t, err)
	_, err = clonedRepo.AddAndCommit("foo1")
	assert.NoError(t, err)
	_, err = clonedRepo.AddAndCommit("foo2")
	assert.NoError(t, err)

	err = clonedRepo.Push(zap.NewNop(), workspaceID)
	assert.NoError(t, err)

	changes, err = ListChanges(repoProvider, codebaseID, workspaceID, 100)
	assert.NoError(t, err)

	if assert.Len(t, changes, 3) {
		assert.Equal(t, "foo2", changes[0].RawCommitMessage)
		assert.False(t, changes[0].IsLanded)

		assert.Equal(t, "foo1", changes[1].RawCommitMessage)
		assert.False(t, changes[1].IsLanded)

		assert.Equal(t, "foo0", changes[2].RawCommitMessage)
		assert.False(t, changes[2].IsLanded)
	}
}

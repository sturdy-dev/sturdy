package vcs_test

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	vcs_change "getsturdy.com/api/pkg/changes/vcs"
	"getsturdy.com/api/pkg/codebases"
	codebasevcs "getsturdy.com/api/pkg/codebases/vcs"
	"getsturdy.com/api/pkg/unidiff"
	viewsvcs "getsturdy.com/api/pkg/view/vcs"
	vcs_workspace "getsturdy.com/api/pkg/workspaces/vcs"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
	"getsturdy.com/api/vcs/testutil"

	git "github.com/libgit2/git2go/v33"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	sig = git.Signature{Name: "test", Email: "test@driva.dev", When: time.Now()}
)

func getDiffs(t *testing.T, repo vcs.RepoReader) []unidiff.FileDiff {
	gitDiffs, err := repo.Diffs()
	assert.NoError(t, err)
	defer gitDiffs.Free()

	diffs, err := unidiff.NewUnidiff(unidiff.NewGitPatchReader(gitDiffs), zap.NewNop()).
		WithExpandedHunks().Decorate()
	assert.NoError(t, err)

	return diffs
}

func TestAddModifyDeleteBinaryFile(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	codebaseID := codebases.ID("codebaseID")
	workspaceID := "workspaceID"
	viewID := "viewID"
	setupCodebase(t, repoProvider, codebaseID, workspaceID, viewID)

	repo, err := repoProvider.ViewRepo(codebaseID, viewID)
	assert.NoError(t, err)
	viewRoot := repoProvider.ViewPath(codebaseID, viewID)

	// Create binary files
	assert.NoError(t, ioutil.WriteFile(viewRoot+"/normal.bin", []byte{0, 0, 0, 0}, 0666))
	assert.NoError(t, ioutil.WriteFile(viewRoot+"/with space.bin", []byte{0, 0, 0, 0}, 0666))

	diffs := getDiffs(t, repo)
	assert.Len(t, diffs, 2)

	_, err = vcs_change.CreateChangeFromPatchesOnRepo(context.Background(), zap.NewNop(), repo, codebaseID, allHunkIDs(diffs), "adding binary files", sig)
	assert.NoError(t, err)

	// Expect empty diff
	diffs = getDiffs(t, repo)
	assert.Len(t, diffs, 0)

	// Modify binary files
	assert.NoError(t, ioutil.WriteFile(viewRoot+"/normal.bin", []byte{0, 0, 0, 0, 1}, 0666))
	assert.NoError(t, ioutil.WriteFile(viewRoot+"/with space.bin", []byte{0, 0, 0, 0, 1}, 0666))

	diffs = getDiffs(t, repo)
	assert.Len(t, diffs, 2)

	_, err = vcs_change.CreateChangeFromPatchesOnRepo(context.Background(), zap.NewNop(), repo, codebaseID, allHunkIDs(diffs), "modify binary files", sig)
	assert.NoError(t, err)

	// Expect empty diff
	diffs = getDiffs(t, repo)
	assert.Len(t, diffs, 0)

	// Remove binary files
	assert.NoError(t, os.Remove(viewRoot+"/normal.bin"))
	assert.NoError(t, os.Remove(viewRoot+"/with space.bin"))

	diffs = getDiffs(t, repo)
	assert.Len(t, diffs, 2)

	_, err = vcs_change.CreateChangeFromPatchesOnRepo(context.Background(), zap.NewNop(), repo, codebaseID, allHunkIDs(diffs), "remove binary files", sig)
	assert.NoError(t, err)

	// Expect empty diff
	diffs = getDiffs(t, repo)
	assert.Len(t, diffs, 0)
}

func setupCodebase(t *testing.T, repoProvider provider.RepoProvider, codebaseID codebases.ID, workspaceID, viewID string) {
	err := codebasevcs.Create(codebaseID)(repoProvider)
	assert.NoError(t, err)

	repo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = vcs_workspace.Create(repo, workspaceID)
	assert.NoError(t, err)

	err = viewsvcs.Create(codebaseID, workspaceID, viewID)(repoProvider)
	assert.NoError(t, err)
}

func TestRevert(t *testing.T) {
	cases := []struct {
		binaryFileName string
	}{
		{
			"new-binary.txt",
		},
		{
			"new binary with space.txt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.binaryFileName, func(t *testing.T) {
			repoProvider := testutil.TestingRepoProvider(t)
			executorProvider := executor.NewProvider(zap.NewNop(), repoProvider)
			codebaseID := codebases.ID("codebaseID")
			workspaceID := "workspaceID"
			viewID := "viewID"
			setupCodebase(t, repoProvider, codebaseID, workspaceID, viewID)

			viewRoot := repoProvider.ViewPath(codebaseID, viewID)

			// Add a file
			err := ioutil.WriteFile(viewRoot+"/new.txt", []byte("hello world"), 0666)
			assert.NoError(t, err)

			// Add a file that git will consider to be binary
			err = ioutil.WriteFile(viewRoot+"/"+tc.binaryFileName, []byte{0, 0, 0, 0}, 0666)
			assert.NoError(t, err)

			repo, err := vcs.OpenRepo(viewRoot)
			assert.NoError(t, err)

			diffs := getDiffs(t, repo)
			assert.Len(t, diffs, 2)

			t.Logf("diffs: %+v", diffs)

			// Remove all diffs
			assert.NoError(t, executorProvider.New().Write(vcs_workspace.Remove(zap.NewNop(), allHunkIDs(diffs)...)).ExecView(codebaseID, viewID, "remove all diffs"))

			// No more diffs!
			diffs = getDiffs(t, repo)
			assert.Empty(t, diffs)

			// Add a binary file, and commit it
			// Then modify the binary file and revert its
			err = ioutil.WriteFile(viewRoot+"/"+tc.binaryFileName, []byte{0, 0, 0, 0}, 0666)
			assert.NoError(t, err)

			_, err = repo.AddAndCommit("Add binary file")
			assert.NoError(t, err)

			// There should be no diffs at this point
			diffs = getDiffs(t, repo)
			assert.Empty(t, diffs)

			// Modify the binary file
			err = ioutil.WriteFile(viewRoot+"/"+tc.binaryFileName, []byte{0, 0, 0, 0, 1, 1, 1}, 0666)
			assert.NoError(t, err)

			// There should be one diff
			diffs = getDiffs(t, repo)
			assert.Len(t, diffs, 1)

			t.Logf("diffs: %+v", diffs)

			// Undo the diff
			assert.NoError(t, executorProvider.New().Write(vcs_workspace.Remove(zap.NewNop(), allHunkIDs(diffs)...)).ExecView(codebaseID, viewID, "remove all diffs"))

			// No more diffs!
			diffs = getDiffs(t, repo)
			assert.Empty(t, diffs)

			// Remove the binary file
			err = os.Remove(viewRoot + "/" + tc.binaryFileName)
			assert.NoError(t, err)

			// There should be one diff
			diffs = getDiffs(t, repo)
			assert.Len(t, diffs, 1)

			t.Logf("diffs: %+v", diffs)

			// Undo the diff
			assert.NoError(t, executorProvider.New().Write(vcs_workspace.Remove(zap.NewNop(), allHunkIDs(diffs)...)).ExecView(codebaseID, viewID, "remove all diffs"))

			// No more diffs!
			diffs = getDiffs(t, repo)
			assert.Empty(t, diffs)

			// Rename the binary file
			err = os.Rename(viewRoot+"/"+tc.binaryFileName, viewRoot+"/renamed-"+tc.binaryFileName)
			assert.NoError(t, err)

			// 1 diff (renamed)
			diffs = getDiffs(t, repo)
			assert.Len(t, diffs, 1)
			assert.True(t, diffs[0].IsMoved)

			t.Logf("diffs: %+v", diffs)

			// Undo the diffs
			assert.NoError(t, executorProvider.New().Write(vcs_workspace.Remove(zap.NewNop(), allHunkIDs(diffs)...)).ExecView(codebaseID, viewID, "remove all diffs"))

			// No more diffs!
			diffs = getDiffs(t, repo)
			assert.Empty(t, diffs)
		})
	}
}

func TestStagedCleanup(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	codebaseID := codebases.ID("codebaseID")
	workspaceID := "workspaceID"
	viewID := "viewID"
	setupCodebase(t, repoProvider, codebaseID, workspaceID, viewID)

	viewPath := repoProvider.ViewPath(codebaseID, viewID)
	t.Log(viewPath)

	err := ioutil.WriteFile(path.Join(viewPath, "a.txt"), []byte("aaa\n"), 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(viewPath, "b.txt"), []byte("bbb\n"), 0777)
	assert.NoError(t, err)

	// Add a.txt to the "staging area" (adding it to the index, without creating a commit)
	cmd := exec.Command("git", "add", "a.txt")
	cmd.Dir = viewPath
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Logf("add: %+v", out)

	r, err := repoProvider.ViewRepo(codebaseID, viewID)
	assert.NoError(t, err)

	// Get diffs
	diffs, err := r.CurrentDiff()
	assert.NoError(t, err)
	defer diffs.Free()

	fileDiffs, err := unidiff.NewUnidiff(unidiff.NewGitPatchReader(diffs), zap.NewNop()).WithExpandedHunks().Decorate()
	assert.NoError(t, err)

	// a.txt and b.txt
	assert.Len(t, fileDiffs, 2)

	var hunkIds []string
	for _, fd := range fileDiffs {
		for _, h := range fd.Hunks {
			hunkIds = append(hunkIds, h.ID)
		}
	}
	assert.Len(t, hunkIds, 2)

	// This should work! Even if we're trying to add an already added file.
	_, err = vcs_change.CreateChangeFromPatchesOnRepo(context.Background(), zap.NewNop(), r, codebaseID, hunkIds, "added all", sig)
	assert.NoError(t, err)

	// Test that it's safe to run CleanStaged even if the "staging area" is empty
	assert.NoError(t, r.CleanStaged())
	assert.NoError(t, r.CleanStaged())
	assert.NoError(t, r.CleanStaged())
}

func allHunkIDs(diffs []unidiff.FileDiff) []string {
	var res []string
	for _, diff := range diffs {
		for _, hunk := range diff.Hunks {
			res = append(res, hunk.ID)
		}
	}
	return res
}

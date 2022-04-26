package vcs

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"getsturdy.com/api/pkg/codebases"
	codebasevcs "getsturdy.com/api/pkg/codebases/vcs"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/view/vcs"
	workspacevcs "getsturdy.com/api/pkg/workspaces/vcs"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/sourcegraph/go-diff/diff"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func diffSnapshots(logger *zap.Logger, viewProvider provider.ViewProvider, codebaseID codebases.ID, viewID, snapshotCommitID, parentSnapshotCommitID string) ([]unidiff.FileDiff, error) {
	repo, err := viewProvider.ViewRepo(codebaseID, viewID)
	if err != nil {
		return nil, err
	}

	diffs, err := repo.DiffCommits(snapshotCommitID, parentSnapshotCommitID)
	if err != nil {
		return nil, err
	}
	defer diffs.Free()

	res, err := unidiff.NewUnidiff(unidiff.NewGitPatchReader(diffs), logger).Decorate()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func TestSnapshot(t *testing.T) {
	reposBasePath, repoProvider := reposBasePath(t)
	codebaseID := codebases.ID("codebaseID")
	err := codebasevcs.Create(codebaseID)(repoProvider)
	assert.NoError(t, err)

	workspaceID := "workspaceID"
	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = workspacevcs.Create(trunkRepo, workspaceID)
	assert.NoError(t, err)

	viewID := "viewID"
	err = vcs.Create(codebaseID, workspaceID, viewID)(repoProvider)
	assert.NoError(t, err)

	viewRepo, err := repoProvider.ViewRepo(codebaseID, viewID)
	assert.NoError(t, err)

	// Write some files
	viewPath := path.Join(reposBasePath, codebaseID.String(), viewID)
	assert.NoError(t, ioutil.WriteFile(viewPath+"/a.txt", []byte("hello a"), 0o666))
	assert.NoError(t, ioutil.WriteFile(viewPath+"/b.txt", []byte("hello b"), 0o666))

	// Snapshot
	firstSnapshotCommitID, err := SnapshotOnViewRepo(zap.NewNop(), viewRepo, codebaseID, uuid.New().String())
	assert.NoError(t, err)
	diffs, _, err := viewRepo.ShowCommit(firstSnapshotCommitID)
	assert.NoError(t, err)
	for _, d := range diffs {
		fd, err := diff.ParseFileDiff([]byte(d))
		assert.NoError(t, err)
		fileDiffMode, _, _, err := unidiff.DiffFileStat(fd)
		assert.NoError(t, err)
		assert.Equal(t, unidiff.FileDiffModeNew, fileDiffMode)
	}

	// Update files some more and snapshot
	assert.NoError(t, ioutil.WriteFile(viewPath+"/a.txt", []byte("hello a2"), 0o666))
	assert.NoError(t, ioutil.WriteFile(viewPath+"/b.txt", []byte("hello b2"), 0o666))
	secondSnapshotCommitID, err := SnapshotOnViewRepo(zap.NewNop(), viewRepo, codebaseID, uuid.New().String())
	assert.NoError(t, err)

	diffs, _, err = viewRepo.ShowCommit(secondSnapshotCommitID)
	assert.NoError(t, err)
	for _, d := range diffs {
		fd, err := diff.ParseFileDiff([]byte(d))
		assert.NoError(t, err)
		fileDiffMode, _, _, err := unidiff.DiffFileStat(fd)
		assert.NoError(t, err)
		// TODO: Compare against the previous snapshot! should be FileDiffModeChanged
		assert.Equal(t, unidiff.FileDiffModeNew, fileDiffMode)
	}

	// Diff the snapshots
	snapshotDiffs, err := diffSnapshots(zap.NewNop(), repoProvider, codebaseID, viewID, firstSnapshotCommitID, secondSnapshotCommitID)
	assert.NoError(t, err)
	assert.Equal(t, []unidiff.FileDiff{
		{
			OrigName:      "a.txt",
			NewName:       "a.txt",
			PreferredName: "a.txt",
			IsDeleted:     false,
			IsNew:         false,
			IsMoved:       false,
			Hunks: []unidiff.Hunk{{
				ID:         "354be1e8201711b68e87c164d107b07aab3b3f88610f85731520865e68387caf",
				Patch:      "diff --git \"a/a.txt\" \"b/a.txt\"\nindex 4d657e1..38cc63e 100644\n--- \"a/a.txt\"\n+++ \"b/a.txt\"\n@@ -1,1 +1,1 @@\n-hello a\n\\ No newline at end of file\n+hello a2\n\\ No newline at end of file\n",
				IsOutdated: false,
				IsApplied:  false,
			}},
		},
		{
			OrigName:      "b.txt",
			NewName:       "b.txt",
			PreferredName: "b.txt",
			IsDeleted:     false,
			IsNew:         false,
			IsMoved:       false,
			Hunks: []unidiff.Hunk{{
				ID:         "c2f9bf1661c77b3c8f8b41768636a8fdd762be72c551d2ba6e970d6b15de3d26",
				Patch:      "diff --git \"a/b.txt\" \"b/b.txt\"\nindex c53170f..a2af30e 100644\n--- \"a/b.txt\"\n+++ \"b/b.txt\"\n@@ -1,1 +1,1 @@\n-hello b\n\\ No newline at end of file\n+hello b2\n\\ No newline at end of file\n",
				IsOutdated: false,
				IsApplied:  false,
			}},
		}}, snapshotDiffs)

	// Remove a file and snapshot
	assert.NoError(t, os.Remove(viewPath+"/a.txt"))
	assert.NoError(t, ioutil.WriteFile(viewPath+"/b.txt", []byte("hello b3"), 0o666))
	_, err = SnapshotOnViewRepo(zap.NewNop(), viewRepo, codebaseID, uuid.New().String())
	assert.NoError(t, err)

	// TODO: Verify snapshot commits, current working directory, etc.
}

func reposBasePath(t *testing.T) (string, provider.RepoProvider) {
	reposBasePath := t.TempDir()
	return reposBasePath, provider.New(reposBasePath, "")
}

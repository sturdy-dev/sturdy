package vcs

import (
	"mash/vcs/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCodebase(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	err := Create(repoProvider, "codebaseID")
	assert.NoError(t, err)
}

func TestImportCodebaseFromGit(t *testing.T) {
	gitURL := "https://github.com/tantivy-search/tantivy-cli.git"
	repoProvider := testutil.TestingRepoProvider(t)
	err := Import(repoProvider, "codebaseID", gitURL)
	assert.NoError(t, err)
}

func TestListChangesInCodebaseTrunk(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	codebaseID := "codebaseID"
	err := Create(repoProvider, codebaseID)
	assert.NoError(t, err)
	repo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	log, err := ListChanges(repo, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(log))
}

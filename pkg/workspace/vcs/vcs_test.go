package vcs

import (
	"testing"

	codebasevcs "mash/pkg/codebase/vcs"
	"mash/vcs/testutil"

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

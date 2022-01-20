package vcs

import (
	codebasevcs "getsturdy.com/api/pkg/codebase/vcs"
	workspacevcs "getsturdy.com/api/pkg/workspace/vcs"
	"getsturdy.com/api/vcs/provider"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateView(t *testing.T) {
	repoProvider := newRepoProvider(t)
	codebaseID := "codebaseID"
	err := codebasevcs.Create(repoProvider, codebaseID)
	assert.NoError(t, err)

	workspaceID := "workspaceID"
	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = workspacevcs.Create(trunkRepo, workspaceID)
	assert.NoError(t, err)

	viewID := "viewID"
	err = Create(repoProvider, codebaseID, workspaceID, viewID)
	assert.NoError(t, err)
}

func TestSetWorkspace(t *testing.T) {
	repoProvider := newRepoProvider(t)
	codebaseID := "codebaseID"
	err := codebasevcs.Create(repoProvider, codebaseID)
	assert.NoError(t, err)

	workspaceID := "workspaceID"
	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = workspacevcs.Create(trunkRepo, workspaceID)
	assert.NoError(t, err)

	viewID := "viewID"
	err = Create(repoProvider, codebaseID, workspaceID, viewID)
	assert.NoError(t, err)

	newWorkspaceID := "ws2"
	err = workspacevcs.Create(trunkRepo, newWorkspaceID)
	assert.NoError(t, err)

	err = SetWorkspace(repoProvider, codebaseID, viewID, newWorkspaceID)
	assert.NoError(t, err)
}

func newRepoProvider(t *testing.T) provider.RepoProvider {
	reposBasePath, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)
	return provider.New(reposBasePath, "")
}

package vcs_test

import (
	"io/ioutil"
	"os"
	"testing"

	vcs_codebase "getsturdy.com/api/pkg/codebase/vcs"
	"getsturdy.com/api/pkg/view/vcs"
	vcs_workspace "getsturdy.com/api/pkg/workspace/vcs"
	"getsturdy.com/api/vcs/provider"

	"github.com/stretchr/testify/assert"
)

func TestCreateView(t *testing.T) {
	repoProvider := newRepoProvider(t)
	codebaseID := "codebaseID"
	err := vcs_codebase.Create(repoProvider, codebaseID)
	assert.NoError(t, err)

	workspaceID := "workspaceID"
	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = vcs_workspace.Create(trunkRepo, workspaceID)
	assert.NoError(t, err)

	viewID := "viewID"
	err = vcs.Create(repoProvider, codebaseID, workspaceID, viewID)
	assert.NoError(t, err)
}

func TestSetWorkspace(t *testing.T) {
	repoProvider := newRepoProvider(t)
	codebaseID := "codebaseID"
	err := vcs_codebase.Create(repoProvider, codebaseID)
	assert.NoError(t, err)

	workspaceID := "workspaceID"
	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	assert.NoError(t, err)
	err = vcs_workspace.Create(trunkRepo, workspaceID)
	assert.NoError(t, err)

	viewID := "viewID"
	err = vcs.Create(repoProvider, codebaseID, workspaceID, viewID)
	assert.NoError(t, err)

	newWorkspaceID := "ws2"
	err = vcs_workspace.Create(trunkRepo, newWorkspaceID)
	assert.NoError(t, err)

	err = vcs.SetWorkspace(repoProvider, codebaseID, viewID, newWorkspaceID)
	assert.NoError(t, err)
}

func newRepoProvider(t *testing.T) provider.RepoProvider {
	reposBasePath, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)
	return provider.New(reposBasePath, "")
}

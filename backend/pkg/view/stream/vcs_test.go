package stream

import (
	"io/ioutil"
	vcs_codebase "mash/pkg/codebase/vcs"
	vcs_view "mash/pkg/view/vcs"
	vcs_workspace "mash/pkg/workspace/vcs"
	"mash/vcs/provider"
	"os"
	"testing"

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
	err = vcs_view.Create(repoProvider, codebaseID, workspaceID, viewID)
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
	err = vcs_view.Create(repoProvider, codebaseID, workspaceID, viewID)
	assert.NoError(t, err)

	newWorkspaceID := "ws2"
	err = vcs_workspace.Create(trunkRepo, newWorkspaceID)
	assert.NoError(t, err)

	err = vcs_view.SetWorkspace(repoProvider, codebaseID, viewID, newWorkspaceID)
	assert.NoError(t, err)
}

func newRepoProvider(t *testing.T) provider.RepoProvider {
	reposBasePath, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)
	return provider.New(reposBasePath, "")
}

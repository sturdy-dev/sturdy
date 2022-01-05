package vcs

import (
	"fmt"
	"io/ioutil"
	"mash/vcs"
	"mash/vcs/provider"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBehindAheadCount(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	repoProvider := provider.New(tmpBase, "")

	codebaseID := "codebaseID"
	workspaceID := "workspaceID"

	pathBase := path.Join(tmpBase, codebaseID, "trunk")
	clientA := path.Join(tmpBase, codebaseID, workspaceID)
	_, err = vcs.CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}
	repoA, err := vcs.CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}
	t.Log(pathBase)

	// Setup branch
	assert.NoError(t, repoA.CreateNewBranchOnHEAD(workspaceID))
	assert.NoError(t, repoA.CheckoutBranchWithForce(workspaceID))
	assert.NoError(t, repoA.Push(zap.NewNop(), workspaceID))

	// Not behind
	behindCount, aheadCount, err := BehindAheadCount(repoProvider, codebaseID, workspaceID)
	assert.NoError(t, err)
	assert.Equal(t, 0, behindCount)
	assert.Equal(t, 0, aheadCount)

	// Switch to sturdytrunk and make some changes there
	assert.NoError(t, repoA.CheckoutBranchWithForce("sturdytrunk"))
	for i := 1; i < 4; i++ {
		assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte(fmt.Sprintf("trunk %d\n", i)), 0666))
		_, err = repoA.AddAndCommit(fmt.Sprintf("%d", i))
		assert.NoError(t, err)
		assert.NoError(t, repoA.Push(zap.NewNop(), "sturdytrunk"))

		// Should be behind now
		behindCount, aheadCount, err = BehindAheadCount(repoProvider, codebaseID, workspaceID)
		assert.NoError(t, err)
		assert.Equal(t, i, behindCount)
		assert.Equal(t, 0, aheadCount)
	}

	// Switch back to the branch, and make some changes that are _ahead_
	assert.NoError(t, repoA.CheckoutBranchWithForce(workspaceID))
	for i := 1; i < 4; i++ {
		assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte(fmt.Sprintf("workspace %d\n", i)), 0666))
		_, err = repoA.AddAndCommit(fmt.Sprintf("%d", i))
		assert.NoError(t, err)
		assert.NoError(t, repoA.Push(zap.NewNop(), workspaceID))

		// Should be ahead
		behindCount, aheadCount, err = BehindAheadCount(repoProvider, codebaseID, workspaceID)
		assert.NoError(t, err)
		assert.Equal(t, 3, behindCount)
		assert.Equal(t, i, aheadCount)
	}
}

package vcs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRevert(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	_, err = CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}
	repoA, err := CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}

	t.Log(clientA)

	assert.NoError(t, repoA.CreateNewBranchOnHEAD("branch"))
	assert.NoError(t, repoA.CheckoutBranchWithForce("branch"))

	// Create first commit
	assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("first!"), 0o666))
	_, err = repoA.AddAndCommit("first")
	assert.NoError(t, err)

	// Create second commit
	assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("second!"), 0o666))
	secondCommitID, err := repoA.AddAndCommit("second")
	assert.NoError(t, err)

	// Third non conflicting commit
	assert.NoError(t, ioutil.WriteFile(clientA+"/b.txt", []byte("third!"), 0o666))
	_, err = repoA.AddAndCommit("third")
	assert.NoError(t, err)

	err = repoA.RevertHEAD(secondCommitID)
	assert.NoError(t, err)

	// Check a.txt
	content, err := ioutil.ReadFile(clientA + "/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "first!", string(content))

	// Check b.txt
	content, err = ioutil.ReadFile(clientA + "/b.txt")
	assert.NoError(t, err)
	assert.Equal(t, "third!", string(content))

	// Check log (no commit created)
	logs, err := repoA.LogBranchUntilTrunk("branch", 4)
	assert.NoError(t, err)
	assert.Equal(t, "third", logs[0].RawCommitMessage)
}

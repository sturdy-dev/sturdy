package vcs

import (
	"io/ioutil"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestListBranches(t *testing.T) {
	repoPath, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	repo, err := CreateBareRepoWithRootCommit(repoPath)
	if err != nil {
		panic(err)
	}

	assert.NoError(t, repo.CreateNewBranchOnHEAD("foo-1"))
	assert.NoError(t, repo.CreateNewBranchOnHEAD("foo-2"))
	assert.NoError(t, repo.CreateNewBranchOnHEAD("foo-3"))

	branches, err := repo.listBranches()
	assert.NoError(t, err)

	assert.Equal(t, []string{"foo-1", "foo-2", "foo-3", "master", "sturdytrunk"}, branches)
}

func TestDiffFromBare(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"

	t.Log("Creating bare base repo at", pathBase)
	bareRepo, err := CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}

	t.Log("Cloning repo to", clientA)
	repoA, err := CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}

	err = repoA.CreateNewBranchOnHEAD("a-branch-name")
	assert.NoError(t, err)

	err = repoA.CheckoutBranchWithForce("a-branch-name")
	assert.NoError(t, err)

	// Add files and commit in client-a
	err = ioutil.WriteFile(clientA+"/foo.txt", []byte("foo foo foo"), 0777)
	assert.NoError(t, err)

	_, err = repoA.AddAndCommit("Commit in a!")
	assert.NoError(t, err)

	err = repoA.Push(zap.NewNop(), "a-branch-name")
	assert.NoError(t, err)

	// Verify logs on sturdytrunk
	logsMaster, err := bareRepo.LogHead(10)
	assert.NoError(t, err)
	assert.Len(t, logsMaster, 1)

	// Verify logs on the new branch, on the bareRepo
	logsBareBranch, err := bareRepo.LogBranchUntilTrunk("a-branch-name", 10)
	assert.NoError(t, err)
	assert.Len(t, logsBareBranch, 1)
}

func TestHeadBranch(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"

	t.Log("Creating bare base repo at", pathBase)
	_, err = CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}

	t.Log("Cloning repo to", clientA)
	repoA, err := CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}

	branchName, err := repoA.HeadBranch()
	assert.NoError(t, err)
	assert.Equal(t, "sturdytrunk", branchName)

	err = repoA.CreateNewBranchOnHEAD("a-branch-name")
	assert.NoError(t, err)

	err = repoA.CheckoutBranchWithForce("a-branch-name")
	assert.NoError(t, err)

	branchName, err = repoA.HeadBranch()
	assert.NoError(t, err)
	assert.Equal(t, "a-branch-name", branchName)
}

func TestPushPull(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	clientB := tmpBase + "client-b"

	t.Logf("b=%s", clientB)

	logger := zap.NewNop()

	_, err = CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}

	repoA, err := CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}
	repoB, err := CloneRepo(pathBase, clientB)
	if err != nil {
		panic(err)
	}

	err = repoA.CreateNewBranchOnHEAD("a-branch-name")
	assert.NoError(t, err)

	err = repoA.CheckoutBranchWithForce("a-branch-name")
	assert.NoError(t, err)

	err = ioutil.WriteFile(clientA+"/hello.txt", []byte("hello world"), 0o644)
	assert.NoError(t, err)

	_, err = repoA.AddAndCommit("commit")
	assert.NoError(t, err)

	err = repoA.Push(logger, "a-branch-name")
	assert.NoError(t, err)

	// fetch on B, and verify contents
	err = repoB.FetchBranch("a-branch-name")
	assert.NoError(t, err)

	// Should fail, no local branch exists
	err = repoB.CheckoutBranchSafely("a-branch-name")
	assert.Error(t, err)

	// Create as new branch
	err = repoB.CreateBranchTrackingUpstream("a-branch-name")
	assert.NoError(t, err)
	err = repoB.CheckoutBranchSafely("a-branch-name")
	assert.NoError(t, err)

	contents, err := ioutil.ReadFile(clientB + "/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(contents))

	// repo a create new commit, and push it
	err = ioutil.WriteFile(clientA+"/hello.txt", []byte("hello world 2222"), 0o644)
	assert.NoError(t, err)
	_, err = repoA.AddAndCommit("commit")
	assert.NoError(t, err)
	err = repoA.Push(logger, "a-branch-name")
	assert.NoError(t, err)

	// fetch again on repo b, and check contents
	err = repoB.CheckoutBranchSafely("sturdytrunk")
	assert.NoError(t, err)
	err = repoB.FetchBranch("a-branch-name")
	assert.NoError(t, err)
	err = repoB.CreateBranchTrackingUpstream("a-branch-name")
	assert.NoError(t, err)
	err = repoB.CheckoutBranchSafely("a-branch-name")
	assert.NoError(t, err)
	contents, err = ioutil.ReadFile(clientB + "/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello world 2222", string(contents))

	// create a new branch on A
	err = repoA.CreateNewBranchOnHEAD("new-branch")
	assert.NoError(t, err)
	err = repoA.Push(logger, "new-branch")
	assert.NoError(t, err)

	// DONT fetch new-branch on B, should not exist
	err = repoB.CheckoutBranchSafely("new-branch")
	assert.Error(t, err)

	// fetch something else, should still not exist
	err = repoB.FetchBranch("a-branch-name")
	assert.NoError(t, err)
	err = repoB.CheckoutBranchSafely("new-branch")
	assert.Error(t, err)
}

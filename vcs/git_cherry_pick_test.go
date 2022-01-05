package vcs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCherryPickNoConflicts(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	clientB := tmpBase + "client-b"
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

	t.Log(pathBase)

	branchName := "branchA"

	logger := zap.NewNop()

	// Create one commit in A
	err = repoA.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoA.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\n"), 0666)
	assert.NoError(t, err)
	_, err = repoA.AddAndCommit("Commit 1 (in A)")
	assert.NoError(t, err)
	err = repoA.Push(logger, branchName)
	assert.NoError(t, err)

	// Create two commits in B
	err = repoB.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoB.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/b.txt", []byte("a new file\n"), 0666)
	assert.NoError(t, err)
	_, err = repoB.AddAndCommit("Commit 1 (in B)")
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/b.txt", []byte("a new file\nupdated\n"), 0666)
	assert.NoError(t, err)
	_, err = repoB.AddAndCommit("Commit 2 (in B)")
	assert.NoError(t, err)

	// Just for the sake of it, make sure that the push fails (non-fast-forward)
	err = repoB.Push(logger, branchName)
	assert.Error(t, err)

	// Ok, here we go. Start the rebasing
	err = repoB.FetchOriginCLI()
	assert.NoError(t, err)
	rebasing, _, err := repoB.InitRebase("origin", branchName)
	assert.NoError(t, err)

	status, err := rebasing.Status()
	assert.NoError(t, err)
	assert.Equal(t, RebaseCompleted, status)

	logs, err := repoB.LogHead(10)
	assert.NoError(t, err)

	if assert.Len(t, logs, 4) { // Root + 1 + 2
		assert.Equal(t, "Commit 2 (in B)", logs[0].RawCommitMessage)
		assert.Equal(t, "Commit 1 (in B)", logs[1].RawCommitMessage)
		assert.Equal(t, "Commit 1 (in A)", logs[2].RawCommitMessage)
		assert.Equal(t, "Root Commit", logs[3].RawCommitMessage)
	}
}

/*
func TestRebaseWithConflict(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	clientB := tmpBase + "client-b"
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

	t.Log(pathBase)

	branchName := "branchA"

	// Create one commit in A
	err = repoA.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoA.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\nhello world\nhello world\n"), 0666)
	assert.NoError(t, err)
	err = repoA.AddAndCommit("Commit 1 (in A)")
	assert.NoError(t, err)
	err = repoA.Push(branchName)
	assert.NoError(t, err)

	// Create two commits in B
	// The first one is conflicting with the commit in repo A
	err = repoB.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoB.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy\nhello sturdy\nhello sturdy\n"), 0666)
	assert.NoError(t, err)
	err = repoB.AddAndCommit("Commit 1 (in B)")
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/b.txt", []byte("a new file\nupdated\n"), 0666)
	assert.NoError(t, err)
	err = repoB.AddAndCommit("Commit 2 (in B)")
	assert.NoError(t, err)

	// Just for the sake of it, make sure that the push fails (non-fast-forward)
	err = repoB.Push(branchName)
	assert.Error(t, err)

	// Ok, here we go. Start the rebasing
	err = repoB.FetchOriginCLI()
	assert.NoError(t, err)
	rebasing, err := repoB.InitRebase("origin", branchName)
	if assert.NoError(t, err) {
		status, err := rebasing.Status()
		assert.NoError(t, err)
		assert.Equal(t, RebaseHaveConflicts, status)
	}

	files, err := rebasing.ConflictingFiles()
	assert.NoError(t, err)
	if assert.Len(t, files, 1) {
		assert.Equal(t, "a.txt", files[0])

		// TODO: Support diffing in this case!
		// err = rebasing.ConflictDiff(files[0])
		// assert.NoError(t, err)
	}

	err = rebasing.ResolveFile("a.txt", "workspace")
	assert.NoError(t, err)

	_, err = rebasing.Continue()
	assert.NoError(t, err)

	status, err := rebasing.Status()
	assert.NoError(t, err)
	assert.Equal(t, RebaseCompleted, status)

	logs, err := repoB.LogHead(10)
	assert.NoError(t, err)

	t.Logf("%+v", logs)
	if assert.Len(t, logs, 4) { // Root + 1 + 2
		assert.Equal(t, "Commit 2 (in B)", logs[0].Message)
		assert.Equal(t, "Commit 1 (in B)", logs[1].Message)
		assert.Equal(t, "Commit 1 (in A)", logs[2].Message)
		assert.Equal(t, "Root Commit", logs[3].Message)
	}
}

func TestRebaseWithConflictReOpen(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	clientB := tmpBase + "client-b"
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

	t.Log(pathBase)

	branchName := "branchA"

	// Create one commit in A
	err = repoA.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoA.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\nhello world\nhello world\n"), 0666)
	assert.NoError(t, err)
	err = repoA.AddAndCommit("Commit 1 (in A)")
	assert.NoError(t, err)
	err = repoA.Push(branchName)
	assert.NoError(t, err)

	// Create two commits in B
	// The first one is conflicting with the commit in repo A
	err = repoB.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoB.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy\nhello sturdy\nhello sturdy\n"), 0666)
	assert.NoError(t, err)
	err = repoB.AddAndCommit("Commit 1 (in B)")
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/b.txt", []byte("a new file\nupdated\n"), 0666)
	assert.NoError(t, err)
	err = repoB.AddAndCommit("Commit 2 (in B)")
	assert.NoError(t, err)

	// Just for the sake of it, make sure that the push fails (non-fast-forward)
	err = repoB.Push(branchName)
	assert.Error(t, err)

	// Ok, here we go. Start the rebasing
	// Creating a new scope (for the rebasing-variable), to make sure that it's not accidentally used afterwards
	{
		err = repoB.FetchOriginCLI()
		assert.NoError(t, err)
		rebasing, err := repoB.InitRebase("origin", branchName)
		if assert.NoError(t, err) {
			status, err := rebasing.Status()
			assert.NoError(t, err)
			assert.Equal(t, RebaseHaveConflicts, status)
		}

		files, err := rebasing.ConflictingFiles()
		assert.NoError(t, err)
		if assert.Len(t, files, 1) {
			assert.Equal(t, "a.txt", files[0])
		}
	}

	// Re-Open the rebasing operation
	rebasing2, err := repoB.OpenRebase()
	if assert.NoError(t, err) {
		status, err := rebasing2.Status()
		assert.NoError(t, err)
		assert.Equal(t, RebaseHaveConflicts, status)
	}

	files2, err := rebasing2.ConflictingFiles()
	assert.NoError(t, err)
	if assert.Len(t, files2, 1) {
		assert.Equal(t, "a.txt", files2[0])
	}

	err = rebasing2.ResolveFile("a.txt", "workspace")
	assert.NoError(t, err)

	stoppedFromConflicts, err := rebasing2.Continue()
	assert.NoError(t, err)
	assert.False(t, stoppedFromConflicts)

	status, err := rebasing2.Status()
	assert.NoError(t, err)
	assert.Equal(t, RebaseCompleted, status)

	logs, err := repoB.LogHead(10)
	assert.NoError(t, err)

	t.Logf("%+v", logs)
	if assert.Len(t, logs, 4) { // Root + 1 + 2
		assert.Equal(t, "Commit 2 (in B)", logs[0].Message)
		assert.Equal(t, "Commit 1 (in B)", logs[1].Message)
		assert.Equal(t, "Commit 1 (in A)", logs[2].Message)
		assert.Equal(t, "Root Commit", logs[3].Message)
	}

	// It should not be possible to open a rebasing action now, as none is ongoing
	_, err = repoB.OpenRebase()
	assert.True(t, errors.Is(err, NoRebaseInProgress))
}

func TestRebaseReOpenNoRebasingOngoing(t *testing.T) {
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
	branchName := "branchA"

	// Create one commit in A
	err = repoA.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoA.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)

	_, err = repoA.OpenRebase()
	assert.True(t, errors.Is(err, NoRebaseInProgress))
}

func TestRebaseWithConflictCommonAncestor(t *testing.T) {
	testCases := []struct {
		resolveVersion     string
		expectedCommits    []string
		withUnsavedChanges bool
	}{
		{
			resolveVersion: "workspace",
			expectedCommits: []string{
				"Commit 1 (in B)",
				"Commit 1 (in A)",
				"Commit 2 (shared history)",
				"Commit 1 (shared history)",
				"Root Commit",
			},
		},
		{
			resolveVersion: "workspace",
			expectedCommits: []string{
				"Commit 1 (in B)",
				"Commit 1 (in A)",
				"Commit 2 (shared history)",
				"Commit 1 (shared history)",
				"Root Commit",
			},
			withUnsavedChanges: true,
		},
		{
			resolveVersion: "trunk",
			expectedCommits: []string{
				// "Commit 1 (in B)", this commit is not there when the trunk version is picked
				"Commit 1 (in A)",
				"Commit 2 (shared history)",
				"Commit 1 (shared history)",
				"Root Commit",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.resolveVersion, func(t *testing.T) {

			tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
			assert.NoError(t, err)

			pathBase := tmpBase + "base"
			clientA := tmpBase + "client-a"
			clientB := tmpBase + "client-b"
			_, err = CreateBareRepoWithRootCommit(pathBase)
			if err != nil {
				panic(err)
			}
			repoA, err := CloneRepo(pathBase, clientA)
			if err != nil {
				panic(err)
			}

			t.Log(pathBase)

			branchName := "branchA"

			var multilinesToBreakDiffContext = strings.Repeat("content\n", 10)

			// Create some common history
			assert.NoError(t, repoA.CreateNewBranchOnHEAD(branchName))
			assert.NoError(t, repoA.CheckoutBranchWithForce(branchName))
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\nhello world\nhello world\n"+multilinesToBreakDiffContext+multilinesToBreakDiffContext), 0666))
			assert.NoError(t, repoA.AddAndCommit("Commit 1 (shared history)"))
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\nhello second\nhello world\n"+multilinesToBreakDiffContext+multilinesToBreakDiffContext), 0666))
			assert.NoError(t, repoA.AddAndCommit("Commit 2 (shared history)"))
			assert.NoError(t, repoA.Push(branchName))

			// Setup Repo B
			repoB, err := CloneRepo(pathBase, clientB)
			if err != nil {
				panic(err)
			}

			// Create one commit in A
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("hello sturdy in trunk\nhello second\nhello world\n"+multilinesToBreakDiffContext+"added here\n"+multilinesToBreakDiffContext), 0666))
			assert.NoError(t, repoA.AddAndCommit("Commit 1 (in A)"))
			assert.NoError(t, repoA.Push(branchName))

			// Create a conflicting commit in B
			assert.NoError(t, repoB.CreateBranchTrackingUpstream(branchName))
			assert.NoError(t, repoB.CheckoutBranchWithForce(branchName))
			assert.NoError(t, ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy in workspace\nhello second\nhello world\n"+multilinesToBreakDiffContext+"added here, but spicy\n"+multilinesToBreakDiffContext), 0666))
			assert.NoError(t, repoB.AddAndCommit("Commit 1 (in B)"))

			// Just for the sake of it, make sure that the push fails (non-fast-forward)
			err = repoB.Push(branchName)
			assert.Error(t, err)

			// Make a local change that is not commited
			if tc.withUnsavedChanges {
				assert.NoError(t, ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy ~~local~~\nhello second\nhello world\n"+multilinesToBreakDiffContext+"added here, but spicy\n"+multilinesToBreakDiffContext), 0666))
			}

			// Ok, here we go. Start the rebasing
			err = repoB.FetchOriginCLI()
			assert.NoError(t, err)
			rebasing, err := repoB.InitRebase("origin", branchName)
			if assert.NoError(t, err) {
				status, err := rebasing.Status()
				assert.NoError(t, err)
				assert.Equal(t, RebaseHaveConflicts, status)
			}

			files, err := rebasing.ConflictingFiles()
			assert.NoError(t, err)
			if assert.Len(t, files, 1) {
				assert.Equal(t, "a.txt", files[0])
				patch, err := rebasing.ConflictDiff(files[0])
				assert.NoError(t, err)
				assert.Equal(t, `diff --git a/a.txt b/a.txt
index 0000000..0000000 81a4
--- a/a.txt
+++ b/a.txt
@@ -1,4 +1,4 @@
-hello world
+hello sturdy in workspace
 hello second
 hello world
 content
@@ -11,6 +11,7 @@ content
 content
 content
 content
+added here, but spicy
 content
 content
 content
`, patch.WorkspacePatch)

				assert.Equal(t, `diff --git a/a.txt b/a.txt
index 0000000..0000000 81a4
--- a/a.txt
+++ b/a.txt
@@ -1,4 +1,4 @@
-hello world
+hello sturdy in trunk
 hello second
 hello world
 content
@@ -11,6 +11,7 @@ content
 content
 content
 content
+added here
 content
 content
 content
`, patch.TrunkPatch)

				// Verify patch with unidiff
				newDiffs := unidiff.ExpandHunks([]string{patch.WorkspacePatch})

				if assert.Len(t, newDiffs, 2) {
					assert.Equal(t, `diff --git a/a.txt b/a.txt
index 0000000..0000000 81a4
--- a/a.txt
+++ b/a.txt
@@ -1,4 +1,4 @@
-hello world
+hello sturdy in workspace
 hello second
 hello world
 content
`, newDiffs[0])
					assert.Equal(t, `diff --git a/a.txt b/a.txt
index 0000000..0000000 81a4
--- a/a.txt
+++ b/a.txt
@@ -11,6 +11,7 @@ content
 content
 content
 content
+added here, but spicy
 content
 content
 content
`, newDiffs[1])
				}
			}

			err = rebasing.ResolveFile("a.txt", tc.resolveVersion)
			assert.NoError(t, err)

			_, err = rebasing.Continue()
			assert.NoError(t, err)

			status, err := rebasing.Status()
			assert.NoError(t, err)
			assert.Equal(t, RebaseCompleted, status)

			logs, err := repoB.LogHead(10)
			assert.NoError(t, err)

			t.Logf("%+v", logs)
			if assert.Len(t, logs, len(tc.expectedCommits)) {
				for i, v := range tc.expectedCommits {
					assert.Equal(t, v, logs[i].Message, "pos=%d", i)
				}
			}

			// Our unsaved changes should have been restored
			if tc.withUnsavedChanges {
				diffs, err := repoB.CurrentDiff()
				assert.NoError(t, err)
				if assert.Len(t, diffs, 1) {
					assert.Equal(t, `diff --git a/a.txt b/a.txt
index ea90cc6..e630172 100644
--- a/a.txt
+++ b/a.txt
@@ -1,4 +1,4 @@
-hello sturdy in workspace
+hello sturdy ~~local~~
 hello second
 hello world
 content
`, diffs[0])
				}
			}
		})
	}
}*/

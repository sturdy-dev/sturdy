package vcs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"

	"mash/pkg/unidiff"
)

func TestRebaseNoConflicts(t *testing.T) {
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
	err = ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\n"), 0666)
	assert.NoError(t, err)
	_, err = repoA.AddAndCommit("Commit 1 (in A)")
	assert.NoError(t, err)
	err = repoA.Push(zap.NewNop(), branchName)
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
	err = repoB.Push(zap.NewNop(), branchName)
	assert.Error(t, err)

	// Ok, here we go. Start the rebasing
	err = repoB.FetchOriginCLI()
	assert.NoError(t, err)
	rebasing, rebasedCommits, err := repoB.InitRebase("origin", branchName)
	assert.NoError(t, err)

	assert.Len(t, rebasedCommits, 2)

	for _, rebasedCommit := range rebasedCommits {
		// Make sure that both versions are the same
		oldCommit, err := repoB.GetCommit(rebasedCommit.OldCommitID)
		assert.NoError(t, err)
		newCommit, err := repoB.GetCommit(rebasedCommit.NewCommitID)
		assert.NoError(t, err)
		assert.Equal(t, oldCommit.Message(), newCommit.Message())
	}

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
	_, err = repoA.AddAndCommit("Commit 1 (in A)")
	assert.NoError(t, err)
	err = repoA.Push(zap.NewNop(), branchName)
	assert.NoError(t, err)

	// Create two commits in B
	// The first one is conflicting with the commit in repo A
	err = repoB.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoB.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy\nhello sturdy\nhello sturdy\n"), 0666)
	assert.NoError(t, err)
	_, err = repoB.AddAndCommit("Commit 1 (in B)")
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/b.txt", []byte("a new file\nupdated\n"), 0666)
	assert.NoError(t, err)
	_, err = repoB.AddAndCommit("Commit 2 (in B)")
	assert.NoError(t, err)

	// Just for the sake of it, make sure that the push fails (non-fast-forward)
	err = repoB.Push(zap.NewNop(), branchName)
	assert.Error(t, err)

	// Ok, here we go. Start the rebasing
	err = repoB.FetchOriginCLI()
	assert.NoError(t, err)
	rebasing, _, err := repoB.InitRebase("origin", branchName)
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

	err = rebasing.ResolveFiles([]SturdyRebaseResolve{{"a.txt", "workspace"}})
	assert.NoError(t, err)

	_, _, err = rebasing.Continue()
	assert.NoError(t, err)

	status, err := rebasing.Status()
	assert.NoError(t, err)
	assert.Equal(t, RebaseCompleted, status)

	logs, err := repoB.LogHead(10)
	assert.NoError(t, err)

	t.Logf("%+v", logs)
	if assert.Len(t, logs, 4) { // Root + 1 + 2
		assert.Equal(t, "Commit 2 (in B)", logs[0].RawCommitMessage)
		assert.Equal(t, "Commit 1 (in B)", logs[1].RawCommitMessage)
		assert.Equal(t, "Commit 1 (in A)", logs[2].RawCommitMessage)
		assert.Equal(t, "Root Commit", logs[3].RawCommitMessage)
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
	_, err = repoA.AddAndCommit("Commit 1 (in A)")
	assert.NoError(t, err)
	err = repoA.Push(zap.NewNop(), branchName)
	assert.NoError(t, err)

	// Create two commits in B
	// The first one is conflicting with the commit in repo A
	err = repoB.CreateNewBranchOnHEAD(branchName)
	assert.NoError(t, err)
	err = repoB.CheckoutBranchWithForce(branchName)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy\nhello sturdy\nhello sturdy\n"), 0666)
	assert.NoError(t, err)
	_, err = repoB.AddAndCommit("Commit 1 (in B)")
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientB+"/b.txt", []byte("a new file\nupdated\n"), 0666)
	assert.NoError(t, err)
	_, err = repoB.AddAndCommit("Commit 2 (in B)")
	assert.NoError(t, err)

	// Just for the sake of it, make sure that the push fails (non-fast-forward)
	err = repoB.Push(zap.NewNop(), branchName)
	assert.Error(t, err)

	// Ok, here we go. Start the rebasing
	// Creating a new scope (for the rebasing-variable), to make sure that it's not accidentally used afterwards
	{
		err = repoB.FetchOriginCLI()
		assert.NoError(t, err)
		rebasing, _, err := repoB.InitRebase("origin", branchName)
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

	err = rebasing2.ResolveFiles([]SturdyRebaseResolve{{"a.txt", "workspace"}})
	assert.NoError(t, err)

	stoppedFromConflicts, _, err := rebasing2.Continue()
	assert.NoError(t, err)
	assert.False(t, stoppedFromConflicts)

	status, err := rebasing2.Status()
	assert.NoError(t, err)
	assert.Equal(t, RebaseCompleted, status)

	logs, err := repoB.LogHead(10)
	assert.NoError(t, err)

	t.Logf("%+v", logs)
	if assert.Len(t, logs, 4) { // Root + 1 + 2
		assert.Equal(t, "Commit 2 (in B)", logs[0].RawCommitMessage)
		assert.Equal(t, "Commit 1 (in B)", logs[1].RawCommitMessage)
		assert.Equal(t, "Commit 1 (in A)", logs[2].RawCommitMessage)
		assert.Equal(t, "Root Commit", logs[3].RawCommitMessage)
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

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", idx, tc.resolveVersion), func(t *testing.T) {

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
			_, err = repoA.AddAndCommit("Commit 1 (shared history)")
			assert.NoError(t, err)
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\nhello second\nhello world\n"+multilinesToBreakDiffContext+multilinesToBreakDiffContext), 0666))
			_, err = repoA.AddAndCommit("Commit 2 (shared history)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), branchName))

			// Setup Repo B
			repoB, err := CloneRepo(pathBase, clientB)
			if err != nil {
				panic(err)
			}

			// Create one commit in A
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("hello sturdy in trunk\nhello second\nhello world\n"+multilinesToBreakDiffContext+"added here (A)\n"+multilinesToBreakDiffContext), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in A)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), branchName))

			// Create a conflicting commit in B
			assert.NoError(t, repoB.CreateBranchTrackingUpstream(branchName))
			assert.NoError(t, repoB.CheckoutBranchWithForce(branchName))
			assert.NoError(t, ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy in workspace\nhello second\nhello world\n"+multilinesToBreakDiffContext+"added here (B)\n"+multilinesToBreakDiffContext), 0666))
			_, err = repoB.AddAndCommit("Commit 1 (in B)")
			assert.NoError(t, err)

			// Just for the sake of it, make sure that the push fails (non-fast-forward)
			err = repoB.Push(zap.NewNop(), branchName)
			assert.Error(t, err)

			// Make a local change that is not commited
			if tc.withUnsavedChanges {
				assert.NoError(t, ioutil.WriteFile(clientB+"/a.txt", []byte("hello sturdy ~~local~~\nhello second\nhello world\n"+multilinesToBreakDiffContext+"added here (B)\n"+multilinesToBreakDiffContext), 0666))
			}

			// Ok, here we go. Start the rebasing
			err = repoB.FetchOriginCLI()
			assert.NoError(t, err)
			rebasing, _, err := repoB.InitRebase("origin", branchName)
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
+added here (B)
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
+added here (A)
 content
 content
 content
`, patch.TrunkPatch)

				// Verify patch with unidiff

				newDiffs, err := unidiff.NewUnidiff(unidiff.NewStringsPatchReader([]string{patch.WorkspacePatch}), zap.NewNop()).WithExpandedHunks().DecorateSingle()
				assert.NoError(t, err)

				if assert.Len(t, newDiffs.Hunks, 2) {
					assert.Equal(t, `diff --git "a/a.txt" "b/a.txt"
index 0000000..0000000 81a4
--- "a/a.txt"
+++ "b/a.txt"
@@ -1,4 +1,4 @@
-hello world
+hello sturdy in workspace
 hello second
 hello world
 content
`, newDiffs.Hunks[0].Patch)
					assert.Equal(t, `diff --git "a/a.txt" "b/a.txt"
index 0000000..0000000 81a4
--- "a/a.txt"
+++ "b/a.txt"
@@ -11,6 +11,7 @@ content
 content
 content
 content
+added here (B)
 content
 content
 content
`, newDiffs.Hunks[1].Patch)
				}
			}

			err = rebasing.ResolveFiles([]SturdyRebaseResolve{{"a.txt", tc.resolveVersion}})
			assert.NoError(t, err)

			data, err := ioutil.ReadFile(path.Join(clientB, "a.txt"))
			assert.NoError(t, err)
			t.Logf("a.txt: %s", string(data))

			if tc.resolveVersion == "workspace" {
				assert.True(t, strings.Contains(string(data), "added here (B)"))
			} else {
				assert.True(t, strings.Contains(string(data), "added here (A)"))
			}

			_, _, err = rebasing.Continue()
			assert.NoError(t, err)

			status, err := rebasing.Status()
			assert.NoError(t, err)
			assert.Equal(t, RebaseCompleted, status)

			logs, err := repoB.LogHead(10)
			assert.NoError(t, err)

			t.Logf("%+v", logs)
			if assert.Len(t, logs, len(tc.expectedCommits)) {
				for i, v := range tc.expectedCommits {
					assert.Equal(t, v, logs[i].RawCommitMessage, "pos=%d", i)
				}
			}

			// Our unsaved changes should have been restored
			if tc.withUnsavedChanges {
				gdiff, err := repoB.CurrentDiff()
				assert.NoError(t, err)
				diffs, err := unidiff.NewUnidiff(unidiff.NewGitPatchReader(gdiff), zap.NewNop()).Patches()
				assert.NoError(t, err)

				if assert.Len(t, diffs, 1) {
					assert.Equal(t, `diff --git "a/a.txt" "b/a.txt"
index e19b89e..3cf0405 100644
--- "a/a.txt"
+++ "b/a.txt"
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
}

func TestConflictInMultipleFiles(t *testing.T) {
	cases := []struct {
		pick                         string
		expectConflictOnSecondCommit bool
		dropFirstCommit              bool
	}{
		{pick: "workspace", expectConflictOnSecondCommit: false},
		{pick: "trunk", expectConflictOnSecondCommit: true},
		{pick: "workspace", dropFirstCommit: true, expectConflictOnSecondCommit: false},
	}

	for _, tc := range cases {

		t.Run(tc.pick, func(t *testing.T) {

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

			t.Log(pathBase)
			t.Log(clientA)

			// Create one commit in A
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("a1\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/b.txt", []byte("b1\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in A)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), "sturdytrunk"))

			assert.NoError(t, repoA.CreateNewBranchOnHEAD("workspace-1"))
			assert.NoError(t, repoA.CreateNewBranchOnHEAD("workspace-2"))

			assert.NoError(t, repoA.CheckoutBranchWithForce("workspace-1"))
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("w1+2\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/b.txt", []byte("w1+2\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in Workspace 1)")
			assert.NoError(t, err)
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("w1+2\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/b.txt", []byte("w1+2\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 2 (in Workspace 1)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), "workspace-1"))

			assert.NoError(t, repoA.CheckoutBranchWithForce("workspace-2"))
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("w2+1\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/b.txt", []byte("w2+1\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/ws2-commit-1.txt", []byte("w2+1\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in Workspace 2)")
			assert.NoError(t, err)
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("w2+2\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/b.txt", []byte("w2+2\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/ws2-commit-2.txt", []byte("w2+1\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 2 (in Workspace 2)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), "workspace-2"))

			// Land workspace-1
			assert.NoError(t, repoA.MoveBranch("sturdytrunk", "workspace-1"))
			assert.NoError(t, repoA.Push(zap.NewNop(), "sturdytrunk"))

			// Rebase workspace-2
			assert.NoError(t, repoA.CheckoutBranchWithForce("workspace-2"))
			rebasing, _, err := repoA.InitRebase("origin", "sturdytrunk")
			if assert.NoError(t, err) {
				status, err := rebasing.Status()
				assert.NoError(t, err)
				assert.Equal(t, RebaseHaveConflicts, status)
			}

			files, err := rebasing.ConflictingFiles()
			assert.NoError(t, err)
			assert.Equal(t, []string{"a.txt", "b.txt"}, files)

			// First commit
			t.Log("resolve a.txt in first")
			assert.NoError(t, rebasing.ResolveFiles([]SturdyRebaseResolve{
				{Path: "a.txt", Version: tc.pick},
				{Path: "b.txt", Version: tc.pick},
			}))

			status, err := rebasing.Status()
			assert.NoError(t, err)
			assert.Equal(t, RebaseCanContinue, status)

			conflicts, _, err := rebasing.Continue()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectConflictOnSecondCommit, conflicts)

			if tc.expectConflictOnSecondCommit {
				// Second commit
				t.Log("resolve a.txt in second")
				assert.NoError(t, rebasing.ResolveFiles([]SturdyRebaseResolve{
					{Path: "a.txt", Version: tc.pick},
					{Path: "b.txt", Version: tc.pick},
				}))

				status, err = rebasing.Status()
				assert.NoError(t, err)
				assert.Equal(t, RebaseCanContinue, status)

				// No more conflicts
				conflicts, _, err = rebasing.Continue()
				assert.NoError(t, err)
				assert.False(t, conflicts)
			}

			status, err = rebasing.Status()
			assert.NoError(t, err)
			assert.Equal(t, RebaseCompleted, status)
		})
	}
}

func TestConflictNoCommonAncestor(t *testing.T) {
	cases := []struct {
		pick string
	}{
		{pick: "workspace"},
		{pick: "trunk"},
	}

	for _, tc := range cases {
		t.Run(tc.pick, func(t *testing.T) {

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

			t.Log(pathBase)
			t.Log(clientA)

			// Create one commit in A
			assert.NoError(t, ioutil.WriteFile(clientA+"/a.txt", []byte("a1\n"), 0666))
			assert.NoError(t, ioutil.WriteFile(clientA+"/b.txt", []byte("b1\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in A)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), "sturdytrunk"))

			assert.NoError(t, repoA.CreateNewBranchOnHEAD("workspace-1"))
			assert.NoError(t, repoA.CreateNewBranchOnHEAD("workspace-2"))

			assert.NoError(t, repoA.CheckoutBranchWithForce("workspace-1"))
			assert.NoError(t, ioutil.WriteFile(clientA+"/new.txt", []byte("in-1\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in Workspace 1)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), "workspace-1"))

			assert.NoError(t, repoA.CheckoutBranchWithForce("workspace-2"))
			assert.NoError(t, ioutil.WriteFile(clientA+"/new.txt", []byte("in-2\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in Workspace 2)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), "workspace-2"))

			// Land workspace-1
			assert.NoError(t, repoA.MoveBranch("sturdytrunk", "workspace-1"))
			assert.NoError(t, repoA.Push(zap.NewNop(), "sturdytrunk"))

			// Rebase workspace-2
			assert.NoError(t, repoA.CheckoutBranchWithForce("workspace-2"))
			rebasing, _, err := repoA.InitRebase("origin", "sturdytrunk")
			if assert.NoError(t, err) {
				status, err := rebasing.Status()
				assert.NoError(t, err)
				assert.Equal(t, RebaseHaveConflicts, status)
			}

			files, err := rebasing.ConflictingFiles()
			assert.NoError(t, err)
			assert.Equal(t, []string{"new.txt"}, files)

			conflictDiffs, err := rebasing.ConflictDiff("new.txt")
			assert.NoError(t, err)

			assert.Equal(t, "diff --git /dev/null b/new.txt\nindex 0000000..0000000 81a4\n--- /dev/null\n+++ b/new.txt\n@@ -0,0 +1 @@\n+in-2\n", conflictDiffs.WorkspacePatch)
			assert.Equal(t, "diff --git /dev/null b/new.txt\nindex 0000000..0000000 81a4\n--- /dev/null\n+++ b/new.txt\n@@ -0,0 +1 @@\n+in-1\n", conflictDiffs.TrunkPatch)

			assert.NoError(t, rebasing.ResolveFiles([]SturdyRebaseResolve{
				{Path: "new.txt", Version: tc.pick},
			}))

			status, err := rebasing.Status()
			assert.NoError(t, err)
			assert.Equal(t, RebaseCanContinue, status)

			conflicts, _, err := rebasing.Continue()
			assert.NoError(t, err)
			assert.False(t, conflicts)

			status, err = rebasing.Status()
			assert.NoError(t, err)
			assert.Equal(t, RebaseCompleted, status)
		})
	}
}

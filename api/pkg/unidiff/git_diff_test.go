package unidiff

import (
	"fmt"
	"getsturdy.com/api/vcs"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCurrentDiffs(t *testing.T) {
	cases := []struct {
		before, after string
		hunkify       bool
		expectedDiff  string
	}{
		{
			before:       "hello world\n",
			after:        "hello WORLD\n",
			expectedDiff: "diff --git \"a/a.txt\" \"b/a.txt\"\nindex 3b18e51..efc781c 100644\n--- \"a/a.txt\"\n+++ \"b/a.txt\"\n@@ -1,1 +1,1 @@\n-hello world\n+hello WORLD\n",
		},
		{
			before:       "hello world\r\n",
			after:        "hello WORLD\r\n",
			expectedDiff: "diff --git \"a/a.txt\" \"b/a.txt\"\nindex f35d3e6..4187f1c 100644\n--- \"a/a.txt\"\n+++ \"b/a.txt\"\n@@ -1,1 +1,1 @@\n-hello world\r\n+hello WORLD\r\n",
		},
		{
			before:       "hello\r\nworld\r\n",
			after:        "hello\r\nworld",
			expectedDiff: "diff --git \"a/a.txt\" \"b/a.txt\"\nindex 23eb407..2930a14 100644\n--- \"a/a.txt\"\n+++ \"b/a.txt\"\n@@ -1,2 +1,2 @@\n hello\r\n-world\r\n+world\n\\ No newline at end of file\n",
		},
		{
			before:       "hello world\n",
			after:        "hello WORLD\n",
			expectedDiff: "diff --git \"a/a.txt\" \"b/a.txt\"\nindex 3b18e51..efc781c 100644\n--- \"a/a.txt\"\n+++ \"b/a.txt\"\n@@ -1,1 +1,1 @@\n-hello world\n+hello WORLD\n",
			hunkify:      true,
		},
		{
			before:       "hello world\r\n",
			after:        "hello WORLD\r\n",
			expectedDiff: "diff --git \"a/a.txt\" \"b/a.txt\"\nindex f35d3e6..4187f1c 100644\n--- \"a/a.txt\"\n+++ \"b/a.txt\"\n@@ -1,1 +1,1 @@\n-hello world\r\n+hello WORLD\r\n",
			hunkify:      true,
		},
		{
			before:       "hello\r\nworld\r\n",
			after:        "hello\r\nworld",
			expectedDiff: "diff --git \"a/a.txt\" \"b/a.txt\"\nindex 23eb407..2930a14 100644\n--- \"a/a.txt\"\n+++ \"b/a.txt\"\n@@ -1,2 +1,2 @@\n hello\r\n-world\r\n+world\n\\ No newline at end of file\n",
			hunkify:      true,
		},
	}

	for idx, tc := range cases {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
			assert.NoError(t, err)

			pathBase := tmpBase + "base"
			repoPath := tmpBase + "client-a"
			_, err = vcs.CreateBareRepoWithRootCommit(pathBase)
			if err != nil {
				panic(err)
			}
			repo, err := vcs.CloneRepo(pathBase, repoPath)
			if err != nil {
				panic(err)
			}

			err = ioutil.WriteFile(path.Join(repoPath, "a.txt"), []byte(tc.before), 0o644)
			assert.NoError(t, err)

			_, err = repo.AddAndCommit("before")
			assert.NoError(t, err)

			err = ioutil.WriteFile(path.Join(repoPath, "a.txt"), []byte(tc.after), 0o644)
			assert.NoError(t, err)

			gdiff, err := repo.CurrentDiff()
			assert.NoError(t, err)
			defer gdiff.Free()
			ud := NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop())
			if tc.hunkify {
				ud = ud.WithExpandedHunks()
			}
			diffs, err := ud.Patches()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedDiff, diffs[0])
		})
	}
}

func TestCurrentDiff(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	_, err = vcs.CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}
	repoA, err := vcs.CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}

	gdiff, err := repoA.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()
	diffNoChanges, err := NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)
	assert.Empty(t, diffNoChanges)

	// Add a file
	err = ioutil.WriteFile(clientA+"/a.txt", []byte("hello world\n"), 0666)
	assert.NoError(t, err)

	gdiff, err = repoA.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()
	diff, err := NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)

	assert.Len(t, diff, 1)
	assert.Equal(t, "diff --git /dev/null \"b/a.txt\"\n"+
		"new file mode 100644\n"+
		"index 0000000..3b18e51\n"+
		"--- /dev/null\n"+
		"+++ \"b/a.txt\"\n"+
		"@@ -0,0 +1,1 @@\n"+
		"+hello world\n", diff[0])

	// Add another file
	err = ioutil.WriteFile(clientA+"/b.txt", []byte("b\nbb\nbbb\n"), 0666)
	assert.NoError(t, err)

	gdiff, err = repoA.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()

	diff, err = NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)

	assert.Equal(t, []string{
		"diff --git /dev/null \"b/a.txt\"\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ \"b/a.txt\"\n@@ -0,0 +1,1 @@\n+hello world\n",
		"diff --git /dev/null \"b/b.txt\"\nnew file mode 100644\nindex 0000000..c2e5d6d\n--- /dev/null\n+++ \"b/b.txt\"\n@@ -0,0 +1,3 @@\n+b\n+bb\n+bbb\n",
	}, diff)

	_, err = repoA.AddAndCommit("Committing to get rid of it")
	assert.NoError(t, err)

	// Add a directory and a file!
	err = os.Mkdir(clientA+"/new-dir", 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/new-dir/in-dir.txt", []byte("file i a dir!\n"), 0666)
	assert.NoError(t, err)

	gdiff, err = repoA.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()

	diff, err = NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)

	assert.Equal(t, "diff --git /dev/null \"b/new-dir/in-dir.txt\"\nnew file mode 100644\nindex 0000000..7f8c1f2\n--- /dev/null\n+++ \"b/new-dir/in-dir.txt\"\n@@ -0,0 +1,1 @@\n+file i a dir!\n", diff[0])
}

func TestCurrentDiffWithNestedNewDirs(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	_, err = vcs.CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}
	repoA, err := vcs.CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}

	gdiff, err := repoA.CurrentDiff()
	assert.NoError(t, err)
	diffNoChanges, err := NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)
	assert.Empty(t, diffNoChanges)

	err = os.Mkdir(clientA+"/new-dir", 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/new-dir/in-dir-a.txt", []byte("file i a dir!\n"), 0666)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/new-dir/in-dir-b.txt", []byte("file i a dir!\n"), 0666)
	assert.NoError(t, err)

	err = os.Mkdir(clientA+"/new-dir/suba", 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/new-dir/suba/a.txt", []byte("file i a dir!\n"), 0666)
	assert.NoError(t, err)

	err = os.Mkdir(clientA+"/new-dir/subb", 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(clientA+"/new-dir/subb/a.txt", []byte("file i a dir!\n"), 0666)
	assert.NoError(t, err)

	gdiff, err = repoA.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()
	diff, err := NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)

	assert.Equal(t,
		[]string{
			"diff --git /dev/null \"b/new-dir/in-dir-a.txt\"\nnew file mode 100644\nindex 0000000..7f8c1f2\n--- /dev/null\n+++ \"b/new-dir/in-dir-a.txt\"\n@@ -0,0 +1,1 @@\n+file i a dir!\n",
			"diff --git /dev/null \"b/new-dir/in-dir-b.txt\"\nnew file mode 100644\nindex 0000000..7f8c1f2\n--- /dev/null\n+++ \"b/new-dir/in-dir-b.txt\"\n@@ -0,0 +1,1 @@\n+file i a dir!\n",
			"diff --git /dev/null \"b/new-dir/suba/a.txt\"\nnew file mode 100644\nindex 0000000..7f8c1f2\n--- /dev/null\n+++ \"b/new-dir/suba/a.txt\"\n@@ -0,0 +1,1 @@\n+file i a dir!\n",
			"diff --git /dev/null \"b/new-dir/subb/a.txt\"\nnew file mode 100644\nindex 0000000..7f8c1f2\n--- /dev/null\n+++ \"b/new-dir/subb/a.txt\"\n@@ -0,0 +1,1 @@\n+file i a dir!\n",
		}, diff)

	// Add a .gitignore (the changes to .gitignore should take affect immediately)
	err = ioutil.WriteFile(clientA+"/.gitignore", []byte("new-dir/suba\nnew-dir/subb\n"), 0666)
	assert.NoError(t, err)

	gdiff, err = repoA.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()
	diff, err = NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)

	expecteds := []string{
		"diff --git /dev/null \"b/.gitignore\"\nnew file mode 100644\nindex 0000000..f796ee4\n--- /dev/null\n+++ \"b/.gitignore\"\n@@ -0,0 +1,2 @@\n+new-dir/suba\n+new-dir/subb\n",
		"diff --git /dev/null \"b/new-dir/in-dir-a.txt\"\nnew file mode 100644\nindex 0000000..7f8c1f2\n--- /dev/null\n+++ \"b/new-dir/in-dir-a.txt\"\n@@ -0,0 +1,1 @@\n+file i a dir!\n",
		"diff --git /dev/null \"b/new-dir/in-dir-b.txt\"\nnew file mode 100644\nindex 0000000..7f8c1f2\n--- /dev/null\n+++ \"b/new-dir/in-dir-b.txt\"\n@@ -0,0 +1,1 @@\n+file i a dir!\n",
	}
	if assert.Len(t, diff, len(expecteds)) {
		for i, d := range diff {
			assert.Equal(t, expecteds[i], d)
		}
	}
}

func TestCurrentDiffDeletedFile(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	_, err = vcs.CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}
	repoA, err := vcs.CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}

	// Add a file, and commit it
	err = ioutil.WriteFile(clientA+"/a.txt", []byte("aaaa!\n"), 0666)
	assert.NoError(t, err)
	_, err = repoA.AddAndCommit("add")
	assert.NoError(t, err)

	// Delete the file
	err = os.Remove(clientA + "/a.txt")
	assert.NoError(t, err)

	gdiff, err := repoA.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()
	diff, err := NewUnidiff(NewGitPatchReader(gdiff), zap.NewNop()).Patches()
	assert.NoError(t, err)

	assert.Len(t, diff, 1)
	assert.Equal(t,
		"diff --git \"a/a.txt\" /dev/null\ndeleted file mode 100644\nindex bbd4a99..0000000\n--- \"a/a.txt\"\n+++ /dev/null\n@@ -1,1 +0,0 @@\n-aaaa!\n",
		string(diff[0]))
}

package vcs

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileContentsAtCommit(t *testing.T) {
	tmpBase := t.TempDir()

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	_, err := CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}
	repoA, err := CloneRepo(pathBase, clientA)
	if err != nil {
		panic(err)
	}

	type expectedContentsAt struct {
		commitID        string
		expectedContent string
	}
	var expecteds []expectedContentsAt

	for _, content := range []string{"1111", "2222", "3333"} {
		err := ioutil.WriteFile(path.Join(clientA, "a.txt"), []byte(content), 0o666)
		assert.NoError(t, err)
		commitID, err := repoA.AddAndCommit(content)
		assert.NoError(t, err)
		expecteds = append(expecteds, expectedContentsAt{commitID: commitID, expectedContent: content})
	}

	for _, expected := range expecteds {
		contents, err := repoA.FileContentsAtCommit(expected.commitID, "a.txt")
		assert.NoError(t, err)
		assert.Equal(t, expected.expectedContent, string(contents))
	}
}

func TestDirectoryChildrenAtCommit(t *testing.T) {
	tmpBase := t.TempDir()

	pathBase := tmpBase + "base"
	clientA := tmpBase + "client-a"
	_, err := CreateBareRepoWithRootCommit(pathBase)
	assert.NoError(t, err)
	repoA, err := CloneRepo(pathBase, clientA)
	assert.NoError(t, err)

	subDir := "subdir"
	subSubDir := path.Join(subDir, "subsubdir")

	err = os.MkdirAll(path.Join(clientA, subSubDir), 0o755)
	assert.NoError(t, err)

	fileA := path.Join(subDir, "fileA")
	err = ioutil.WriteFile(path.Join(clientA, fileA), []byte{0, 1, 2}, 0o666)
	assert.NoError(t, err)

	fileB := path.Join(subSubDir, "fileB")
	err = ioutil.WriteFile(path.Join(clientA, fileB), []byte{0, 1, 2}, 0o666)
	assert.NoError(t, err)

	commitID, err := repoA.AddAndCommit("add files")
	assert.NoError(t, err)

	rootChildren, err := repoA.DirectoryChildrenAtCommit(commitID, "/")
	assert.NoError(t, err)
	assert.Equal(t, rootChildren, []string{subDir})

	subDirChildren, err := repoA.DirectoryChildrenAtCommit(commitID, subDir)
	assert.NoError(t, err)
	assert.Equal(t, subDirChildren, []string{fileA, subSubDir})

	subSubDirChildren, err := repoA.DirectoryChildrenAtCommit(commitID, subSubDir)
	assert.NoError(t, err)
	assert.Equal(t, subSubDirChildren, []string{fileB})
}

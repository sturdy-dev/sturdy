package vcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCommitWithFiles(t *testing.T) {
	repoPath := t.TempDir()

	repo, err := CreateBareRepoWithRootCommit(repoPath)
	assert.NoError(t, err)

	commitID, err := repo.CreateCommitWithFiles([]FileContents{
		{"README.md", []byte("# Hello World!")},
	}, "new-branch-name")
	assert.NoError(t, err)

	contents, err := repo.FileContentsAtCommit(commitID, "README.md")
	assert.NoError(t, err)
	assert.Equal(t, "# Hello World!", string(contents))
}

package vcs

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/unidiff"
	"io/ioutil"
	"path"
	"testing"
)

func TestCanApplyPatch(t *testing.T) {
	tmpBase := t.TempDir()

	pathBase := path.Join(tmpBase, "trunk")
	pathCheckout := path.Join(tmpBase, "checkout")
	_, err := CreateBareRepoWithRootCommit(pathBase)
	if err != nil {
		panic(err)
	}
	repo, err := CloneRepo(pathBase, pathCheckout)
	if err != nil {
		panic(err)
	}
	t.Log(pathCheckout)

	logger := zap.NewNop()

	err = ioutil.WriteFile(path.Join(pathCheckout, "a.txt"), []byte("one\ntwo\nthree\n"), 0o666)
	assert.NoError(t, err)

	_, err = repo.AddAndCommit("A")
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(pathCheckout, "a.txt"), []byte("one\nyaya\nthree\n"), 0o666)
	assert.NoError(t, err)

	gdiff, err := repo.CurrentDiff()
	assert.NoError(t, err)
	defer gdiff.Free()
	diffYaya, err := unidiff.NewUnidiff(unidiff.NewGitPatchReader(gdiff), logger).PatchesBytes()
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(pathCheckout, "a.txt"), []byte("foo\nbar\nbaz\n"), 0o666)
	assert.NoError(t, err)

	canApply, err := repo.CanApplyPatch(diffYaya[0])
	assert.False(t, canApply)
	assert.NoError(t, err)

	// reset file to previous state
	err = ioutil.WriteFile(path.Join(pathCheckout, "a.txt"), []byte("one\ntwo\nthree\n"), 0o666)
	assert.NoError(t, err)

	// the patch should now apply
	canApply, err = repo.CanApplyPatch(diffYaya[0])
	assert.True(t, canApply)
	assert.NoError(t, err)

	// working dir should be unchanged
	workdirContents, err := ioutil.ReadFile(path.Join(pathCheckout, "a.txt"))
	assert.Equal(t, []byte("one\ntwo\nthree\n"), workdirContents)
	assert.NoError(t, err)
}

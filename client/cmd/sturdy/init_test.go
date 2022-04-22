package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathIsInGitRepo_IsOutside(t *testing.T) {
	tmpDir := t.TempDir()
	is, gitPath := pathIsInGitRepo(tmpDir)
	assert.False(t, is)
	assert.Empty(t, gitPath)

	// Test with a dir that does not exist
	is, gitPath = pathIsInGitRepo(path.Join(tmpDir, "does-not", "exist"))
	assert.False(t, is)
	assert.Empty(t, gitPath)
}

func TestPathIsInGitRepo_IsInside(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	err := cmd.Run()
	assert.NoError(t, err)

	is, gitPath := pathIsInGitRepo(tmpDir)
	assert.True(t, is)
	assert.Equal(t, path.Join(tmpDir, ".git"), gitPath)

	// Dir inside git repo, .git is not in the same dir
	nested := path.Join(tmpDir, "foo", "bar", "nested")
	err = os.MkdirAll(nested, 0o777)
	assert.NoError(t, err)

	is, gitPath = pathIsInGitRepo(nested)
	assert.True(t, is)

	// for testing on mac, "/var/" and "/private/var/" are the same dir
	gitPath = strings.ReplaceAll(gitPath, "/private/var/", "/var/")

	assert.Equal(t, path.Join(tmpDir, ".git"), gitPath)

	// Test with dir that does not exist (inside git repo)
	is, gitPath = pathIsInGitRepo(path.Join(tmpDir, "does", "not", "exist"))
	assert.True(t, is)
	assert.Equal(t, path.Join(tmpDir, ".git"), gitPath)
}

func TestPathIsSub(t *testing.T) {
	cases := []struct {
		a, b     string
		expected bool
	}{
		{"/foo/bar", "/foo", true},
		{"/foo", "/foo/bar", true},
		{"/foo/aaaa", "/foo/bbbb", false},
		{"/foo/la/da", "/foo/la/di", false},
		{"/foo/la/da", "/foo/la/daba", false},
		{"/foo/la/da", "/foo/la/da/ba", true},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s --> %s", tc.a, tc.b), func(t *testing.T) {
			isSub, err := pathIsSub(tc.a, tc.b)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, isSub)
		})
	}
}

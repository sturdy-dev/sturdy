package ignore

import (
	"os"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/assert"
)

func TestIgnore(t *testing.T) {
	res, err := FindIgnore(os.DirFS("testdata/ignores"))
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"hello-*",
		"/this/that/foobar.txt",
		"/this/that/*.swp",
	}, res)
}

func TestNoRecursionInIgnored(t *testing.T) {
	fs := memfs.New()

	err := fs.MkdirAll("foo/bar", 0o777)
	assert.NoError(t, err)
	err = fs.WriteFile("foo/bar/.gitignore", []byte("in-nested-ignored.txt\n"), 0o644)
	assert.NoError(t, err)

	err = fs.MkdirAll("nested/other", 0o777)
	assert.NoError(t, err)
	err = fs.WriteFile("nested/other/.gitignore", []byte("in-nested-not-ignored.txt\n"), 0o644)
	assert.NoError(t, err)

	err = fs.MkdirAll("nested/f.tmp", 0o777)
	assert.NoError(t, err)
	err = fs.WriteFile("nested/f.tmp/.gitignore", []byte("in-nested-dot-tmp-should-be-ignored.txt\n"), 0o644)
	assert.NoError(t, err)

	err = fs.WriteFile(".gitignore", []byte("in-root.txt\n/foo\n*.swp\n/.DS_Store\n*.tmp\n"), 0o644)
	assert.NoError(t, err)

	res, err := FindIgnore(fs)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"in-root.txt",
		"/foo",
		"*.swp",
		"/.DS_Store",
		"*.tmp",
		"/nested/other/in-nested-not-ignored.txt",
	}, res)
}

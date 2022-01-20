package unidiff_test

import (
	"fmt"
	"testing"

	"getsturdy.com/api/pkg/unidiff"

	"github.com/stretchr/testify/assert"
)

type allowerTestValue struct {
	path      string
	directory bool
	expected  bool
}

type allowerTestCase struct {
	allows []string
	tests  []allowerTestValue
}

func (c *allowerTestCase) run(t *testing.T) {
	allower, err := unidiff.NewAllower(c.allows...)
	if err != nil {
		t.Fatal("unable to create ignorer:", err)
	}

	for _, p := range c.tests {
		if p.directory {
			t.Run(fmt.Sprintf("file: %s", p.path), func(t *testing.T) {
				assert.Equal(t, p.expected, allower.IsAllowed(p.path, p.directory))
			})
		} else {
			t.Run(fmt.Sprintf("dir: %s", p.path), func(t *testing.T) {
				assert.Equal(t, p.expected, allower.IsAllowed(p.path, p.directory))
			})
		}
	}
}

func TestAllower_None(t *testing.T) {
	test := &allowerTestCase{
		allows: nil,
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"something", false, false},
			{"something", true, false},
			{"some/path", false, false},
			{"some/path", true, false},
		},
	}
	test.run(t)
}

func TestAllower_Basic(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"something",
			"otherthing",
			"!something",
			"somedir/",
		},
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"something", false, false},
			{"something", true, false},
			{"subpath/something", false, false},
			{"subpath/something", true, false},
			{"otherthing", false, true},
			{"otherthing", true, true},
			{"subpath/otherthing", false, true},
			{"subpath/otherthing", true, true},
			{"random", false, false},
			{"random", true, false},
			{"subpath/random", false, false},
			{"subpath/random", true, false},
			{"somedir", false, false},
			{"somedir", true, true},
			{"subpath/somedir", false, false},
			{"subpath/somedir", true, true},
		},
	}
	test.run(t)
}

func TestAllower_Group(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"*.py[cod]",
			"*.dir[cod]/",
		},
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"run.py", false, false},
			{"run.pyc", false, true},
			{"run.pyc", true, true},
			{"subpath/run.pyd", false, true},
			{"subpath/run.pyd", true, true},
			{"run.dir", false, false},
			{"run.dir", true, false},
			{"run.dirc", false, false},
			{"run.dirc", true, true},
			{"subpath/run.dird", false, false},
			{"subpath/run.dird", true, true},
		},
	}
	test.run(t)
}

func TestAllower_RootRelative(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"/abspath",
			"/absdir/",
			"/name",
			"!*/**/name",
		},
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"abspath", false, true},
			{"abspath", true, true},
			{"subpath/abspath", false, false},
			{"subpath/abspath", true, false},
			{"absdir", false, false},
			{"absdir", true, true},
			{"subpath/absdir", false, false},
			{"subpath/absdir", true, false},
			{"name", false, true},
			{"name", true, true},
			{"subpath/name", false, false},
			{"subpath/name", true, false},
		},
	}
	test.run(t)
}

func TestAllower_Doublestar(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"some/*",
			"some/**/*",
			"!some/other",
		},
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"something", false, false},
			{"some", false, false},
			{"some/path", false, true},
			{"some/other", false, false},
			{"some/other/path", false, true},
		},
	}
	test.run(t)
}

func TestAllower_NegateOrdering(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"!something",
			"otherthing",
			"something",
		},
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"something", false, true},
			{"something/other", false, false},
			{"otherthing", false, true},
			{"some/path", false, false},
		},
	}
	test.run(t)
}

func TestAllower_Wildcard(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"some*",
			"!someone",
		},
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"som", false, false},
			{"some", false, true},
			{"something", false, true},
			{"someone", false, false},
			{"some/path", false, false},
		},
	}
	test.run(t)
}

func TestAllower_PathWildcard(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"some/*",
			"some/**/*",
			"!some/other",
		},
		tests: []allowerTestValue{
			{"", false, false},
			{"", true, false},
			{"something", false, false},
			{"some", false, false},
			{"some/path", false, true},
			{"some/other", false, false},
			{"some/other/path", false, true},
			{"subdir/some/other/path", false, false},
		},
	}
	test.run(t)
}

func TestAllower_GitIsAlwaysHidden(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"*",
		},
		tests: []allowerTestValue{
			{".git", true, false},
			{".git", false, false},
			{".git/file", false, false},
			{".git/dir", true, false},
			{".git/dir/file", false, false},
			{"/.git", true, false},
			{"../.git", true, false},
			{"dir/.git", true, false},
			{"dir/path/.git", true, false},
		},
	}
	test.run(t)
}

func TestAllower_GitIsAlwaysHidden_evenIfAllowed(t *testing.T) {
	test := &allowerTestCase{
		allows: []string{
			"*",
			".git",
			".git/**/*",
		},
		tests: []allowerTestValue{
			{".git", true, false},
			{".git", false, false},
			{".git/file", false, false},
			{".git/dir", true, false},
			{".git/dir/file", false, false},
			{"/.git", true, false},
			{"../.git", true, false},
			{"dir/.git", true, false},
			{"dir/path/.git", true, false},
		},
	}
	test.run(t)
}

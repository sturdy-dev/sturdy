package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToKnownHosts(t *testing.T) {
	cases := []struct {
		existing      string
		addRows       []string
		lineSeparator string
		expected      string
	}{
		{
			existing:      "",
			addRows:       []string{"foo", "bar"},
			lineSeparator: "\n",
			expected:      "foo\nbar\n",
		},
		{
			existing:      "something else\n",
			addRows:       []string{"foo", "bar"},
			lineSeparator: "\n",
			expected:      "something else\nfoo\nbar\n",
		},
		{
			existing:      "",
			addRows:       []string{"foo", "bar"},
			lineSeparator: "\r\n",
			expected:      "foo\r\nbar\r\n",
		},
		{
			existing:      "",
			addRows:       []string{"foo\r\n", "bar\r\n"},
			lineSeparator: "\r\n",
			expected:      "foo\r\nbar\r\n",
		},
		{
			existing:      "",
			addRows:       []string{"foo\n", "bar\n"},
			lineSeparator: "\r\n",
			expected:      "foo\r\nbar\r\n",
		},
		{
			existing:      "something else\r\n",
			addRows:       []string{"foo", "bar"},
			lineSeparator: "\r\n",
			expected:      "something else\r\nfoo\r\nbar\r\n",
		},
		{
			existing:      "something else\r\nwhat\r\n\r\n\r\n\r\n",
			addRows:       []string{"foo", "bar"},
			lineSeparator: "\r\n",
			expected:      "something else\r\nwhat\r\nfoo\r\nbar\r\n",
		},
	}

	for idx, tc := range cases {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			res := addToKnownHosts(tc.addRows, []byte(tc.existing), tc.lineSeparator)
			assert.Equal(t, tc.expected, string(res))
		})
	}
}

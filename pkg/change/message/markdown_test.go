package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownToHtml(t *testing.T) {
	res, err := MarkdownToHtml(`* here
* we
* go!

This is **fun**!`)
	assert.NoError(t, err)
	assert.Equal(t, "<ul>\n<li>here</li>\n<li>we</li>\n<li>go!</li>\n</ul>\n<p>This is <strong>fun</strong>!</p>\n", res)
}

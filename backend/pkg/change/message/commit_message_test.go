package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitMessage(t *testing.T) {
	input := "<p>Sturdy is epic, riiiiiight?!</p><h3>Third size</h3><ul><li><p>XO XO</p></li><li><p><strong>HO</strong> HO</p></li><li><p>YO <strong><em>NO</em></strong></p></li></ul><p><mark>Highlight</mark></p>"
	expected := "Sturdy is epic, riiiiiight?!\r\n\r\nThird size\r\n\r\n* XO XO\r\n* HO HO\r\n* YO NO\r\nHighlight"

	out := CommitMessage(input)
	assert.Equal(t, expected, out)
}

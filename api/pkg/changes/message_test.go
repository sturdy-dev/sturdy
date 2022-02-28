package changes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	cm := ChangeMetadata{
		Description: "hello everyone",
		UserID:      "user-123",
		ChangeID:    "change-123",
		ViewID:      "view-123",
	}

	message := cm.ToCommitMessage()
	assert.Equal(t, "hello everyone\r\n\r\n--- sturdy ---\r\nchange_id: change-123\r\nuser_id: user-123\r\nview_id: view-123", message)

	// parse from message
	parsed := ParseCommitMessage(message)

	// Should equal the same result
	assert.Equal(t, cm, parsed)
}

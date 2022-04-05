package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	tests := map[string]string{
		"test":         "test",
		":emoji: test": "test",
	}
	for in, out := range tests {
		assert.Equal(t, out, sanitize(in))
	}
}

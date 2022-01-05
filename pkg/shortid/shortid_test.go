package shortid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	for i := 0; i < 100; i++ {
		res := New()
		assert.True(t, len(res) == 7)
	}
}

package img

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThumb(t *testing.T) {
	cases := []struct {
		name string
	}{
		{name: "400x400.jpg"},
		{name: "avatar-0.png"},
		{name: "pattern.png"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := os.CreateTemp("", "sturdy-img-test")
			assert.NoError(t, err)
			defer output.Close()

			input, err := os.Open("testdata/" + tc.name)
			assert.NoError(t, err)
			defer input.Close()

			err = Thumbnail(100, input, output)
			assert.NoError(t, err)

			t.Log(output.Name())
		})
	}
}

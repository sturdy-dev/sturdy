package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"getsturdy.com/api/pkg/file"
)

func TestFileType(t *testing.T) {
	cases := []struct {
		path     string
		expected file.Type
	}{
		{"./testdata/image.png", file.ImageType},
		{"./testdata/text.txt", file.BinaryType},
	}

	for _, tc := range cases {
		s := New(nil)
		fp, err := os.OpenFile(tc.path, os.O_RDONLY, 0o644)
		assert.NoError(t, err)
		res, err := s.detectFileType(fp)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, res)
	}
}

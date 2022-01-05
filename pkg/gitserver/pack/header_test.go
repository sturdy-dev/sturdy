package pack

import (
	"io/ioutil"
	"testing"

	"github.com/bmizerany/assert"
)

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    Header
		expectedErr error
	}{
		{
			name:        "default",
			input:       []byte("00a80000000000000000000000000000000000000000 2ab8b0433111e6d5602a71049e40902c1e5a556c refs/heads/my-branch\x00 report-status side-band-64k agent=git/2.24.3.(Apple.Git-128)0000PACK�;ب��j�\\�<�>�"),
			expected:    Header{Branch: "my-branch"},
			expectedErr: nil,
		},
		{
			name:        "default-2",
			input:       []byte("00c3e424d72b9db65aca594f00a39e61a53dbb767ea4 7b9d96bd14e7a41a93d7b232e81d9c7a5ec87563 refs/heads/953966be-1741-44fe-abd1-bc5d63c26b73 report-status side-band-64k agent=git/2.24.3.(Apple.Git-128)0000PACK�x���"),
			expected:    Header{Branch: "953966be-1741-44fe-abd1-bc5d63c26b73"},
			expectedErr: nil,
		},
		{
			name: "normal-as-raw-bytes",
			input: []byte{
				// 00c33dde561cb61878fb44f9f526ad48b0b25d804558 [SPACE]
				48, 48, 99, 51, 51, 100, 100, 101, 53, 54, 49, 99, 98, 54, 49, 56, 55, 56, 102, 98, 52, 52, 102, 57, 102, 53, 50, 54, 97, 100, 52, 56, 98, 48, 98, 50, 53, 100, 56, 48, 52, 53, 53, 56, 32,
				// 23fce78d0105cec482cd29abc6af9dc3d95d45c5 [SPACE]
				98, 97, 54, 52, 98, 57, 48, 54, 101, 98, 98, 100, 54, 100, 53, 102, 101, 51, 102, 100, 56, 100, 99, 57, 55, 50, 50, 48, 49, 100, 102, 98, 52, 49, 57, 51, 100, 50, 51, 55, 32,
				// refs/heads/953966be-1741-44fe-abd1-bc5d63c26b73 [NUL][SPACE]
				114, 101, 102, 115, 47, 104, 101, 97, 100, 115, 47, 57, 53, 51, 57, 54, 54, 98, 101, 45, 49, 55, 52, 49, 45, 52, 52, 102, 101, 45, 97, 98, 100, 49, 45, 98, 99, 53, 100, 54, 51, 99, 50, 54, 98, 55, 51, 0, 32,
				// report-status side-band-64k agent=git/2.24.3.(Apple.Git-128)0000PAC....
				// more...
				114, 101, 112, 111, 114, 116, 45, 115, 116, 97, 116, 117, 115, 32, 115, 105, 100, 101, 45, 98, 97, 110, 100, 45, 54, 52, 107, 32, 97, 103, 101, 110, 116, 61, 103, 105, 116, 47, 50, 46, 50, 52, 46, 51, 46, 40, 65, 112, 112, 108, 101, 46, 71, 105, 116, 45, 49, 50, 56, 41, 48, 48, 48, 48, 80, 65, 67, 75, 0, 0, 0, 2, 0, 0, 0, 7, 146, 14, 120, 156, 157, 140, 65,
			},
			expected:    Header{Branch: "953966be-1741-44fe-abd1-bc5d63c26b73"},
			expectedErr: nil,
		},
		{
			name:        "libgit2-import",
			input:       mustReadFile("testdata/libgit2.bin"),
			expected:    Header{Branch: "sturdytrunk"},
			expectedErr: nil,
		},
		{
			name:        "libgit2-fresh-import",
			input:       mustReadFile("testdata/libgit2-fresh.bin"),
			expected:    Header{Branch: "sturdytrunk"},
			expectedErr: nil,
		},
		{
			name:        "kube-score-import",
			input:       mustReadFile("testdata/kube-score.bin"),
			expected:    Header{Branch: "sturdytrunk"},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHeader(tt.input)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func mustReadFile(name string) []byte {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return data
}

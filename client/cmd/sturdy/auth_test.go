package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"

	"getsturdy.com/client/cmd/sturdy/config"

	"github.com/stretchr/testify/assert"
)

func TestReadUntilValidTokenInput(t *testing.T) {
	type args struct {
		termReadWriter io.ReadWriter
		validationFunc validateTokenFunc
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "valid on first try",
			args: args{termReadWriter: bytes.NewBufferString("xoxo\n\n"), validationFunc: func(conf *config.Config, checkToken string) error {
				if checkToken == "xoxo" {
					return nil
				}
				return fmt.Errorf("invalid")
			}},
			want: "xoxo",
		},
		{
			name: "valid on second try",
			args: args{
				termReadWriter: &scheduledReader{strs: []string{"xoxo\r\n", "bobo\r\n", "hobo\r\n"}},
				validationFunc: func(conf *config.Config, checkToken string) error {
					log.Println("Validation func", checkToken)
					if checkToken == "bobo" {
						return nil
					}
					return fmt.Errorf("invalid")
				},
			},
			want: "bobo",
		},
		{
			name: "out of attempts",
			args: args{
				termReadWriter: &scheduledReader{strs: []string{
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
					"\r\n", "\r\n", "\r\n", "\r\n", "\r\n",
				}},
				validationFunc: func(conf *config.Config, checkToken string) error {
					return fmt.Errorf("invalid")
				},
			},
			wantErr: errOutOfTries,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf config.Config
			got, err := readUntilValidTokenInput(&conf, tt.args.termReadWriter, tt.args.validationFunc)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

type scheduledReader struct {
	i    int
	strs []string
}

func (r *scheduledReader) Read(p []byte) (n int, err error) {
	if r.i >= len(r.strs) {
		return 0, io.EOF
	}

	d := []byte(r.strs[r.i])
	l := copy(p, d)
	r.i++
	return l, nil
}

func (r *scheduledReader) Write(p []byte) (n int, err error) {
	return len(p), nil
}

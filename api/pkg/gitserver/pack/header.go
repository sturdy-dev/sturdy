package pack

import (
	"bytes"
	"errors"
	"strings"
)

type Header struct {
	Branch string
}

var ErrFail = errors.New("failed to parse git header")
var ErrFailRefs = errors.New("failed to parse git header refs")

func ParseHeader(header []byte) (Header, error) {
	// Examples
	// 00a80000000000000000000000000000000000000000 2ab8b0433111e6d5602a71049e40902c1e5a556c refs/heads/my-branch report-status side-band-64k agent=git/2.24.3.(Apple.Git-128)0000PACK
	// 00c3e424d72b9db65aca594f00a39e61a53dbb767ea4 7b9d96bd14e7a41a93d7b232e81d9c7a5ec87563 refs/heads/953966be-1741-44fe-abd1-bc5d63c26b73 report-status side-band-64k agent=git/2.24.3.(Apple.Git-128)0000PACK
	// ?                                            "Last known commit by server"            branch
	parts := bytes.Split(header, []byte(" "))
	if len(parts) < 3 {
		return Header{}, ErrFail
	}

	refs := parts[2]
	if !bytes.HasPrefix(refs, []byte("refs/heads/")) {
		return Header{}, ErrFailRefs
	}

	// Trim trailing NUL from refs
	branch := string(refs)[len("refs/heads/"):]
	branch = strings.Trim(branch, "\x00")

	return Header{
		Branch: branch,
	}, nil
}

package uploader

import (
	"context"
	"io"

	"getsturdy.com/api/pkg/users/avatars"
)

type Uploader interface {
	Upload(context.Context, string, io.Reader) (*avatars.Avatar, error)
}

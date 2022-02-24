package uploader

import (
	"context"
	"io"

	"getsturdy.com/api/pkg/configuration/flags"
	"getsturdy.com/api/pkg/users/avatars"
)

type Uploader interface {
	Upload(context.Context, string, io.Reader) (*avatars.Avatar, error)
}

type Configuration struct {
	URL flags.URL `long:"url" description:"Avatars base url" default:"http://127.0.0.1:3000/"`
}

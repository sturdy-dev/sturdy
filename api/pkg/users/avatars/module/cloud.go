//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/users/avatars/uploader"
)

func Module(c *di.Container) {
	c.Register(uploader.NewS3, new(uploader.Uploader))
}

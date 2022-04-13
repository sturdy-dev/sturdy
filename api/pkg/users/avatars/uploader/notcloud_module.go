//go:build !cloud
// +build !cloud

package uploader

import (
	service_blobs "getsturdy.com/api/pkg/blobs/service"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(configuration.Module)
	c.Import(service_blobs.Module)
	c.Register(NewBlobs, new(Uploader))
}

//go:build cloud
// +build cloud

package uploader

import (
	"getsturdy.com/api/pkg/aws"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(aws.Module)
	c.Register(NewS3, new(Uploader))
}

//go:build !cloud
// +build !cloud

package aws

import (
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {}

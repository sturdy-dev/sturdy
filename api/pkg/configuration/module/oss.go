//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(configuration.New)
}

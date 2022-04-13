//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(New)
}

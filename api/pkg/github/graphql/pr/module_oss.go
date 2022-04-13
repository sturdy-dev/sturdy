//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package pr

import (
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(NewResolver)
}

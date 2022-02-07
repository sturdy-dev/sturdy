//go:build !cloud
// +build !cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/emails"
)

func Module(c *di.Container) {
	c.Register(emails.NewDisabled)
}

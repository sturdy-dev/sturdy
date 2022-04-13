//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/emails/transactional"
)

func Module(c *di.Container) {
	c.Import(transactional.Module)
	c.Register(func(t *transactional.Sender) transactional.EmailSender {
		return t
	})
}

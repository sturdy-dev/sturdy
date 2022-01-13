//go:build !enterprise
// +build !enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/emails/transactional"
)

func Module(c *di.Container) {
	c.Register(transactional.New, new(transactional.EmailSender))
}

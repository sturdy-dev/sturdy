//go:build !enterprise && !cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/integrations/providers"
)

func Module(c *di.Container) {
	c.Register(func() providers.Providers {
		return providers.Providers{}
	})
}

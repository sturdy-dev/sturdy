//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/configuration/enterprise/selfhosted"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(selfhosted.New)
}

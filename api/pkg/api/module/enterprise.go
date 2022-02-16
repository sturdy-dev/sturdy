//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/api/enterprise/selfhosted"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	common(c)
	c.Register(api.ProvideAPI)
	c.Import(selfhosted.Module)
}

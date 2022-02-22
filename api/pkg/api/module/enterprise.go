//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/api/enterprise/selfhosted"
	"getsturdy.com/api/pkg/di"
	module_queue "getsturdy.com/api/pkg/queue/module"
)

func Module(c *di.Container) {
	common(c)
	c.Import(module_queue.Module)

	c.Register(api.ProvideAPI)
	c.Import(selfhosted.Module)
}

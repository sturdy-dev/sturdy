//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/api/enterprise/selfhosted"
	"getsturdy.com/api/pkg/di"
	module_queue "getsturdy.com/api/pkg/queue/module"
	module_remote "getsturdy.com/api/pkg/remote/module"
)

func Module(c *di.Container) {
	common(c)
	c.Import(module_queue.Module)
	c.Import(module_remote.Module)

	c.Register(api.ProvideAPI)
	c.Import(selfhosted.Module)
}

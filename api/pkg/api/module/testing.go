package module

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/di"
	module_queue "getsturdy.com/api/pkg/queue/module"
)

func TestingModule(c *di.Container) {
	common(c)
	c.Import(module_queue.TestingModule)
	c.Register(api.ProvideAPI, new(api.Starter))
}

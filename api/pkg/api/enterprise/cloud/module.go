package cloud

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(api.Module)
	c.Register(ProvideAPI, new(api.Starter))
}

package http

import (
	"getsturdy.com/api/pkg/di"
	handler "getsturdy.com/api/pkg/http/handler/module"
)

func Module(c *di.Container) {
	c.Import(handler.Module)
	c.Register(ProvideServer)
}

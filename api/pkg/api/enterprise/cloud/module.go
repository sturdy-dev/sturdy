package cloud

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(ProvideAPI, new(api.Starter))
}

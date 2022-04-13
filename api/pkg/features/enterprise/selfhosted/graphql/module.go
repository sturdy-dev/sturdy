package graphql

import (
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(configuration.Module)
	c.Register(NewFeaturesRootResolver)
}

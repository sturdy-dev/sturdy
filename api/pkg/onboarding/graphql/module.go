package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/onboarding/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(logger.Module)
	c.Import(events.Module)
	c.Register(NewRootResolver)
}

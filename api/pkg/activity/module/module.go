package module

import (
	"getsturdy.com/api/pkg/activity/db"
	"getsturdy.com/api/pkg/activity/graphql"
	"getsturdy.com/api/pkg/activity/sender"
	"getsturdy.com/api/pkg/activity/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
	c.Import(sender.Module)
}

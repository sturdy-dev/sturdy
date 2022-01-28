package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspace/activity/db"
	"getsturdy.com/api/pkg/workspace/activity/graphql"
	"getsturdy.com/api/pkg/workspace/activity/sender"
	"getsturdy.com/api/pkg/workspace/activity/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
	c.Import(sender.Module)
}

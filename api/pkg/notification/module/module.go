package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/notification/db"
	"getsturdy.com/api/pkg/notification/graphql"
	"getsturdy.com/api/pkg/notification/sender"
	"getsturdy.com/api/pkg/notification/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(sender.Module)
	c.Import(service.Module)
}

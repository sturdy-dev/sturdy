package sender

import (
	db_activity "getsturdy.com/api/pkg/activity/db"
	service_activity "getsturdy.com/api/pkg/activity/service"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
)

func Module(c *di.Container) {
	c.Import(db_codebases.Module)
	c.Import(db_activity.Module)
	c.Import(events.Module)
	c.Import(service_activity.Module)
	c.Register(NewActivitySender)
}

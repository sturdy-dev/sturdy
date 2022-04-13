package sender

import (
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	transactional "getsturdy.com/api/pkg/emails/transactional/module"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/logger"
	db_notifications "getsturdy.com/api/pkg/notification/db"
	db_users "getsturdy.com/api/pkg/users/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_codebases.Module)
	c.Import(db_notifications.Module)
	c.Import(db_users.Module)
	c.Import(events.Module)
	c.Import(transactional.Module)
	c.Register(NewNotificationSender)
}

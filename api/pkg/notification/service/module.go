package service

import (
	"getsturdy.com/api/pkg/di"
	db_notificaton "getsturdy.com/api/pkg/notification/db"
)

func Module(c *di.Container) {
	c.Import(db_notificaton.Module)
	c.Register(NewPreferences)
}

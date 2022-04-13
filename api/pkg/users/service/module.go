package service

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	db_user "getsturdy.com/api/pkg/users/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_user.Module)
	c.Import(service_analytics.Module)
	c.Register(New)
}

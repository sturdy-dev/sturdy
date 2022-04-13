package service

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/di"
	transactional "getsturdy.com/api/pkg/emails/transactional/module"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/logger"
	service_onetime "getsturdy.com/api/pkg/onetime/service"
	service_organization "getsturdy.com/api/pkg/organization/service"
	db_user "getsturdy.com/api/pkg/users/db"
	service_users "getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(service_users.Module)
	c.Import(logger.Module)
	c.Import(db_user.Module)
	c.Import(service_jwt.Module)
	c.Import(transactional.Module)
	c.Import(service_onetime.Module)
	c.Import(service_analytics.Module)
	c.Import(service_organization.Module)
	c.Register(New)
	c.Register(func(s *Service) service_users.Service {
		return s
	})
}

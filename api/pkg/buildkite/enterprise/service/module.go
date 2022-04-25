package service

import (
	db_buildkite "getsturdy.com/api/pkg/buildkite/enterprise/db"
	"getsturdy.com/api/pkg/buildkite/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db_buildkite.Module)
	c.Register(New)
	c.Register(func(svc *Service) service.Service { return svc })
}

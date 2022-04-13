package service

import (
	"getsturdy.com/api/pkg/di"
	db_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/db"
)

func Module(c *di.Container) {
	c.Import(db_buildkite.Module)
	c.Register(New)
}

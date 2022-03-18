package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/integrations/db"
	"getsturdy.com/api/pkg/integrations/graphql"
	module_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/module"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(module_buildkite.Module)
}

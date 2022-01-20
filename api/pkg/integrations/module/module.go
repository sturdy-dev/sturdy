package module

import (
	"mash/pkg/di"
	module_buildkite "mash/pkg/integrations/buildkite/module"
	"mash/pkg/integrations/db"
	"mash/pkg/integrations/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(module_buildkite.Module)
}

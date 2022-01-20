package enterprise

import (
	"mash/pkg/di"
	"mash/pkg/integrations/buildkite/enterprise/db"
	"mash/pkg/integrations/buildkite/enterprise/graphql"
	"mash/pkg/integrations/buildkite/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}

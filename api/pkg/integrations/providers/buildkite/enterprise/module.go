package enterprise

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/db"
	"getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/graphql"
	"getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}

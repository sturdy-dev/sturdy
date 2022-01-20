package enterprise

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/integrations/buildkite/enterprise/db"
	"getsturdy.com/api/pkg/integrations/buildkite/enterprise/graphql"
	"getsturdy.com/api/pkg/integrations/buildkite/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}

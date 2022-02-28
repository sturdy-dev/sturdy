package cloud

import (
	"getsturdy.com/api/pkg/changes/downloads/enterprise/cloud/graphql"
	"getsturdy.com/api/pkg/changes/downloads/enterprise/cloud/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(graphql.Module)
}

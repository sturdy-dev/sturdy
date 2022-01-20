package module

import (
	"getsturdy.com/api/pkg/codebase/acl/db"
	"getsturdy.com/api/pkg/codebase/acl/graphql"
	"getsturdy.com/api/pkg/codebase/acl/provider"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(provider.Module)
	c.Import(graphql.Module)
	c.Import(db.Module)
}

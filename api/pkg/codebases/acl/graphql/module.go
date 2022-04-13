package graphql

import (
	provider_acl "getsturdy.com/api/pkg/codebases/acl/provider"
	"getsturdy.com/api/pkg/di"
	db_user "getsturdy.com/api/pkg/users/db"
)

func Module(c *di.Container) {
	c.Import(db_user.Module)
	c.Import(provider_acl.Module)
	c.Register(NewResolver)
}

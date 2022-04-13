package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/pki/db"
	graphql_users "getsturdy.com/api/pkg/users/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql_users.Module)
	c.Register(NewResolver)
}

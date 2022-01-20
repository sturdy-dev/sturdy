package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/pki/db"
	"getsturdy.com/api/pkg/pki/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
}

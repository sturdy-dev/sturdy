package module

import (
	"getsturdy.com/api/pkg/crypto/db"
	"getsturdy.com/api/pkg/crypto/graphql"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
}

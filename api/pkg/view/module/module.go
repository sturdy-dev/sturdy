package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/view/graphql"
	"getsturdy.com/api/pkg/view/meta"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(meta.Module)
}

package migrate

import (
	"getsturdy.com/api/pkg/db/migrate/data"
	"getsturdy.com/api/pkg/db/migrate/schema"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(data.Module)
	c.Import(schema.Module)
	c.Register(New)
}

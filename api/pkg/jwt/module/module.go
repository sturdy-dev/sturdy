package module

import (
	"getsturdy.com/api/pkg/di"
	keys_db "getsturdy.com/api/pkg/jwt/keys/db"
	"getsturdy.com/api/pkg/jwt/service"
)

func Module(c *di.Container) {
	c.Import(keys_db.Module)
	c.Import(service.Module)
}

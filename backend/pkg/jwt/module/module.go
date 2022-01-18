package module

import (
	"mash/pkg/di"
	keys_db "mash/pkg/jwt/keys/db"
	"mash/pkg/jwt/service"
)

func Module(c *di.Container) {
	c.Import(keys_db.Module)
	c.Import(service.Module)
}

package graphql

import (
	db_crypto "getsturdy.com/api/pkg/crypto/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db_crypto.Module)
	c.Register(New)
}

package db

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/organization/db"

	"github.com/jmoiron/sqlx"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(func(db *sqlx.DB) Repository {
		return NewCache(New(db))
	})
}

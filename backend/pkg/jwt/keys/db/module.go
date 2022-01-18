package db

import (
	"mash/pkg/di"

	"github.com/jmoiron/sqlx"
)

func Module(c *di.Container) {
	c.Register(func(db *sqlx.DB) Repository {
		return NewCache(New(db))
	})
}

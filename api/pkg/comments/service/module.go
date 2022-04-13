package service

import (
	db_comments "getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db_comments.Module)
	c.Register(New)
}

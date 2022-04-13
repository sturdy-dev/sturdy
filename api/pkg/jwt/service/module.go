package service

import (
	"getsturdy.com/api/pkg/di"
	db_keys "getsturdy.com/api/pkg/jwt/keys/db"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_keys.Module)
	c.Register(NewService)
}

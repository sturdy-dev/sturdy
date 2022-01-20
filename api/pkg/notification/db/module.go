package db

import "getsturdy.com/api/pkg/di"

func Module(c *di.Container) {
	c.Register(NewRepository)
	c.Register(NewPeferenceRepository)
}

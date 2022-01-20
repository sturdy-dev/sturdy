package db

import "mash/pkg/di"

func Module(c *di.Container) {
	c.Register(NewDatabase)
}

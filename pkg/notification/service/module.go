package service

import "mash/pkg/di"

func Module(c *di.Container) {
	c.Register(NewPreferences)
}

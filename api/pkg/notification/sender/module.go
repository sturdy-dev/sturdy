package sender

import "mash/pkg/di"

func Module(c *di.Container) {
	c.Register(NewNotificationSender)
}

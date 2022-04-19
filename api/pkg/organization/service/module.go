package service

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	sender_notifications "getsturdy.com/api/pkg/notification/sender"
	db_organization "getsturdy.com/api/pkg/organization/db"
)

func Module(c *di.Container) {
	c.Import(events.Module)
	c.Import(db_organization.Module)
	c.Import(service_analytics.Module)
	c.Import(sender_notifications.Module)
	c.Register(New)
}

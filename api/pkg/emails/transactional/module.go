package transactional

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	db_comments "getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/di"
	emails "getsturdy.com/api/pkg/emails/module"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/logger"
	db_newsletter "getsturdy.com/api/pkg/newsletter/db"
	service_notification "getsturdy.com/api/pkg/notification/service"
	db_organizations "getsturdy.com/api/pkg/organization/db"
	db_review "getsturdy.com/api/pkg/review/db"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	db_users "getsturdy.com/api/pkg/users/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(emails.Module)
	c.Import(db_users.Module)
	c.Import(db_codebases.Module)
	c.Import(db_comments.Module)
	c.Import(db_suggestions.Module)
	c.Import(db_review.Module)
	c.Import(db_newsletter.Module)
	c.Import(db_workspaces.Module)
	c.Import(service_jwt.Module)
	c.Import(service_change.Module)
	c.Import(service_notification.Module)
	c.Import(service_analytics.Module)
	c.Import(db_organizations.Module)
	c.Register(New)
}

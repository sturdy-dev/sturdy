package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	graphql_codebases "getsturdy.com/api/pkg/codebases/graphql"
	graphql_comments "getsturdy.com/api/pkg/comments/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	graphql_github "getsturdy.com/api/pkg/github/graphql"
	"getsturdy.com/api/pkg/logger"
	db_notification "getsturdy.com/api/pkg/notification/db"
	service_notification "getsturdy.com/api/pkg/notification/service"
	graphql_review "getsturdy.com/api/pkg/review/graphql"
	graphql_suggestions "getsturdy.com/api/pkg/suggestions/graphql"
	graphql_workspaces "getsturdy.com/api/pkg/workspaces/graphql"
)

func Module(c *di.Container) {
	c.Import(db_notification.Module)
	c.Import(db_codebases.Module)
	c.Import(service_notification.Module)
	c.Import(service_auth.Module)
	c.Import(events.Module)
	c.Import(logger.Module)
	c.Import(graphql_comments.Module)
	c.Import(graphql_codebases.Module)
	c.Import(graphql_author.Module)
	c.Import(graphql_workspaces.Module)
	c.Import(graphql_review.Module)
	c.Import(graphql_suggestions.Module)
	c.Import(graphql_github.Module)
	c.Register(NewResolver)
}

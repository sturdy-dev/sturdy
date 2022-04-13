package graphql

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	graphql_changes "getsturdy.com/api/pkg/changes/graphql"
	graphql_acl "getsturdy.com/api/pkg/codebases/acl/graphql"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	graphql_file "getsturdy.com/api/pkg/file/graphql"
	graphql_github "getsturdy.com/api/pkg/github/graphql"
	"getsturdy.com/api/pkg/graphql/resolvers"
	graphql_integrations "getsturdy.com/api/pkg/integrations/graphql"
	"getsturdy.com/api/pkg/logger"
	service_organization "getsturdy.com/api/pkg/organization/service"
	graphql_remote "getsturdy.com/api/pkg/remote/graphql/module"
	service_remote "getsturdy.com/api/pkg/remote/service/module"
	db_user "getsturdy.com/api/pkg/users/db"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(db_codebases.Module)
	c.Import(db_view.Module)
	c.Import(db_workspaces.Module)
	c.Import(db_user.Module)
	c.Import(logger.Module)
	c.Import(events.Module)
	c.Import(service_analytics.Module)
	c.Import(executor.Module)
	c.Import(service_auth.Module)
	c.Import(service_codebase.Module)
	c.Import(service_organization.Module)
	c.Import(service_remote.Module)
	c.Import(resolvers.Module)
	c.Import(graphql_author.Module)
	c.Import(graphql_acl.Module)
	c.Import(graphql_changes.Module)
	c.Import(graphql_file.Module)
	c.Import(graphql_integrations.Module)
	c.Import(graphql_github.Module)
	c.Import(graphql_remote.Module)
	c.Register(NewCodebaseRootResolver)

	// populate cyclic resolver
	c.Decorate(func(rv *resolvers.CodebaseRootResolver, rp resolvers.CodebaseRootResolver) *resolvers.CodebaseRootResolver {
		*rv = rp
		return &rp
	})
}

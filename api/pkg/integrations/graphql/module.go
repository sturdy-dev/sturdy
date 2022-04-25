package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/di"
	graphql_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/graphql"
	integration_providers "getsturdy.com/api/pkg/integrations/providers/module"
	graphql_statuses "getsturdy.com/api/pkg/statuses/graphql/module"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(service_ci.Module)
	c.Import(service_auth.Module)
	c.Import(service_change.Module)
	c.Import(service_workspaces.Module)
	c.Import(graphql_buildkite.Module)
	c.Import(graphql_statuses.Module)
	c.Import(integration_providers.Module)
	c.Register(NewRootResolver)
}

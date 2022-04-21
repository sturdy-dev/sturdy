package graphql

import (
	graphql_activity "getsturdy.com/api/pkg/activity/graphql"
	graphql_changes "getsturdy.com/api/pkg/changes/graphql"
	graphql_acl "getsturdy.com/api/pkg/codebases/acl/graphql"
	graphql_codebases "getsturdy.com/api/pkg/codebases/graphql"
	graphql_comments "getsturdy.com/api/pkg/comments/graphql"
	graphql_crypto "getsturdy.com/api/pkg/crypto/graphql"
	"getsturdy.com/api/pkg/di"
	graphql_features "getsturdy.com/api/pkg/features/graphql"
	graphql_github "getsturdy.com/api/pkg/github/graphql"
	graphql_installations "getsturdy.com/api/pkg/installations/graphql/module"
	graphql_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/graphql"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	graphql_land "getsturdy.com/api/pkg/land/graphql"
	graphql_licenses "getsturdy.com/api/pkg/licenses/graphql"
	"getsturdy.com/api/pkg/logger"
	graphql_notification "getsturdy.com/api/pkg/notification/graphql"
	graphql_onboarding "getsturdy.com/api/pkg/onboarding/graphql"
	graphql_organizations "getsturdy.com/api/pkg/organization/graphql"
	graphql_pki "getsturdy.com/api/pkg/pki/graphql"
	graphql_servicetokens "getsturdy.com/api/pkg/servicetokens/graphql"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_jwt.Module)
	c.Import(graphql_acl.Module)
	c.Import(graphql_activity.Module)
	c.Import(graphql_buildkite.Module)
	c.Import(graphql_changes.Module)
	c.Import(graphql_github.Module)
	c.Import(graphql_codebases.Module)
	c.Import(graphql_comments.Module)
	c.Import(graphql_crypto.Module)
	c.Import(graphql_features.Module)
	c.Import(graphql_licenses.Module)
	c.Import(graphql_notification.Module)
	c.Import(graphql_onboarding.Module)
	c.Import(graphql_organizations.Module)
	c.Import(graphql_pki.Module)
	c.Import(graphql_installations.Module)
	c.Import(graphql_servicetokens.Module)
	c.Import(graphql_land.Module)
	c.Register(NewRootResolver)
}

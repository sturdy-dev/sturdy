//go:build enterprise || cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/integrations/providers"
	service_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(service_buildkite.Module)

	c.Register(func(buildkite *service_buildkite.Service) providers.Providers {
		return providers.Providers{
			buildkite.ProviderName(): buildkite,
		}
	})
}

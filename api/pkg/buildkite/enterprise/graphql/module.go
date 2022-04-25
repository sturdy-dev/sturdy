package graphql

import (
	"log"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

func Module(c *di.Container) {
	log.Println("ZEGL --- ENTERPRISE BUILDKITE")

	c.Import(service_ci.Module)
	c.Import(service_auth.Module)
	c.Import(resolvers.Module)
	c.Register(New)
}

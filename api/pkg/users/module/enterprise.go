//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/users/enterprise/selfhosted"
	"getsturdy.com/api/pkg/users/graphql"
	selfhosted_oss_service "getsturdy.com/api/pkg/users/oss/selfhosted/service"
	"getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Register(service.New)
	c.Import(selfhosted.Module)
	c.Import(selfhosted_oss_service.Module)
}

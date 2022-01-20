//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/license/enterprise/client"
	"getsturdy.com/api/pkg/license/enterprise/db"
	"getsturdy.com/api/pkg/license/enterprise/graphql"
	"getsturdy.com/api/pkg/license/enterprise/service"
	"getsturdy.com/api/pkg/license/enterprise/validator"
)

func Module(c *di.Container) {
	c.Register(db.New)
	c.Register(db.NewValidationRepository)
	c.Register(service.New)
	c.Register(client.New)
	c.Register(validator.New)
	c.Register(graphql.New)
}

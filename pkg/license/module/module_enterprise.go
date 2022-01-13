//go:build enterprise
// +build enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/license/enterprise/client"
	"mash/pkg/license/enterprise/db"
	"mash/pkg/license/enterprise/graphql"
	"mash/pkg/license/enterprise/service"
	"mash/pkg/license/enterprise/validator"
)

func Module(c *di.Container) {
	c.Register(db.New)
	c.Register(db.NewValidationRepository)
	c.Register(service.New)
	c.Register(client.New)
	c.Register(validator.New)
	c.Register(graphql.New)
}

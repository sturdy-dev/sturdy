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

var Module = di.NewModule(
	di.Provides(db.New),
	di.Provides(db.NewValidationRepository),
	di.Provides(service.New),
	di.Provides(client.New),
	di.Provides(validator.New),
	di.ProvidesCycle(graphql.New),
)

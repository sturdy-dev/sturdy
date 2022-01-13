//go:build !cloud
// +build !cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/serverstatus/selfhosted/graphql"
	"mash/pkg/serverstatus/selfhosted/service"
)

var Module = di.NewModule(
	di.ProvidesCycle(graphql.New),
	di.ProvidesCycle(service.New),
)

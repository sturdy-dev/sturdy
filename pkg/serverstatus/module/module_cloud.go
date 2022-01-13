//go:build cloud
// +build cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/serverstatus/cloud/graphql"
)

var Module = di.NewModule(
	di.ProvidesCycle(graphql.New),
)

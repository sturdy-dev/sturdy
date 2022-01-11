//go:build enterprise
// +build enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/workspace/service"
	"mash/pkg/workspace/service/enterprise"
)

var Module = di.NewModule(
	di.Provides(service.New),
	di.Provides(enterprise.New, new(service.Service)),
)

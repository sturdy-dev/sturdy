//go:build !enterprise
// +build !enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/workspace/service"
)

var Module = di.NewModule(
	di.Provides(service.New, new(service.Service)),
)

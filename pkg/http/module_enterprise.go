//go:build enterprise
// +build enterprise

package http

import (
	"mash/pkg/di"
	"mash/pkg/http/enterprise"
	"mash/pkg/http/oss"
)

var Module = di.NewModule(
	di.Provides(oss.ProvideHandler),
	di.Provides(enterprise.ProvideHandler),
	di.Provides(ProvideServer),
)

package http

import "mash/pkg/di"

var Module = di.NewModule(
	di.Provides(ProvideServer),
	di.Provides(ProvideHandler),
)

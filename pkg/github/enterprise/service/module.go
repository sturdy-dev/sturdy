package service

import (
	"mash/pkg/di"
)

var Module = di.NewModule(
	di.Provides(New),
)

package api

import (
	"mash/pkg/di"
)

var Module = di.NewModule(
	di.Provides(ProvideAPI),
)

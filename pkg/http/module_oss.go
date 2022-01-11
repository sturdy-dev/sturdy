//go:build !enterprise
// +build !enterprise

package http

import (
	"mash/pkg/di"
	"mash/pkg/http/oss"
	"net/http"
)

var Module = di.NewModule(
	di.Provides(oss.ProvideHandler, new(http.Handler)),
	di.Provides(ProvideServer),
)

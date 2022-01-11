//go:build enterprise
// +build enterprise

package api

import (
	"mash/pkg/api"
	"mash/pkg/api/enterprise"
	"mash/pkg/api/oss"
	"mash/pkg/di"
)

var Module = di.NewModule(
	di.Provides(oss.ProvideAPI),
	di.Provides(enterprise.ProvideAPI, new(api.API)),
)

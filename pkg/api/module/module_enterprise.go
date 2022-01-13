//go:build enterprise
// +build enterprise

package api

import (
	"mash/pkg/api"
	"mash/pkg/api/enterprise"
	"mash/pkg/api/oss"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideAPI)
	c.Register(enterprise.ProvideAPI, new(api.API))
}

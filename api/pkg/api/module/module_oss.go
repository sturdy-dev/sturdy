//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package api

import (
	"mash/pkg/api"
	"mash/pkg/api/oss"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideAPI, new(api.API))
}

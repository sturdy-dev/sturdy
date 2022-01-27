//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package api

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(api.ProvideAPI, new(api.Starter))
}

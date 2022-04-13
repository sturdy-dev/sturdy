//go:build !cloud && !enterprise
// +build !cloud,!enterprise

package api

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(api.Module)
	c.Register(func(s *api.API) api.Starter { return s })
}

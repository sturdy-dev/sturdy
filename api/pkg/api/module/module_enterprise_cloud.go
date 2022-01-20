//go:build enterprise || cloud
// +build enterprise cloud

package api

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/api/enterprise"
	"getsturdy.com/api/pkg/api/oss"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideAPI)
	c.Register(enterprise.ProvideAPI, new(api.API))
}

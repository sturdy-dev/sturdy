//go:build enterprise
// +build enterprise

package api

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/api/enterprise/selfhosted"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(api.ProvideAPI)
	c.Import(selfhosted.Module)
}

//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	service_buildkite "getsturdy.com/api/pkg/buildkite/enterprise/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(service_buildkite.Module)
}

//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/di"
	enterprise "getsturdy.com/api/pkg/users/enterprise/selfhosted/service"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
}

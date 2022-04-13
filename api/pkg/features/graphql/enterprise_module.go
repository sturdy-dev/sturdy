//go:build enterprise
// +build enterprise

package graphql

import (
	"getsturdy.com/api/pkg/di"
	enterprise "getsturdy.com/api/pkg/features/enterprise/selfhosted/graphql"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
}

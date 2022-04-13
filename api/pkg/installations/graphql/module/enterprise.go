//go:build enterprise
// +build enterprise

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/enterprise/selfhosted/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}

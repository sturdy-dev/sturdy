//go:build !enterprise
// +build !enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/integrations/buildkite/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}

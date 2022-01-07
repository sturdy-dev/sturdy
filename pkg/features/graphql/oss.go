// +build !enterprise

package graphql

import (
	"mash/pkg/features/graphql/oss"
)

var Resolver = oss.Resolver

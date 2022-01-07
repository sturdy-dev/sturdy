// +build enterprise

package schema

import (
	_ "embed"
)

var (
	//go:embed enterprise.graphql
	enterpriseSchema string
	//go:embed oss.graphql
	ossSchema string
	//go:embed schema.graphql
	schema string
)

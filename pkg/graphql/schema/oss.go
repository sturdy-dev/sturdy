// +build !enterprise

package schema

import (
	_ "embed"
)

var (
	enterpriseSchema string = ""
	//go:embed oss.graphql
	ossSchema string
	//go:embed schema.graphql
	schema string
)

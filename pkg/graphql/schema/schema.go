package schema

import (
	_ "embed"
	"strings"
)

var (
	//go:embed cloud.graphql
	cloudSchema string
	//go:embed enterprise.graphql
	enterpriseSchema string
	//go:embed oss.graphql
	ossSchema string
	//go:embed schema.graphql
	schema string
)

var String = strings.Join([]string{
	schema,
	ossSchema,
	enterpriseSchema,
	cloudSchema,
}, "\n")

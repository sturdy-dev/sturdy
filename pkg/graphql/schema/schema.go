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

// String contains the full GraphQL Schema of Sturdy
//
// The schema is the same in all versions of Sturdy (OSS, Enterprise) to make it easier for clients to
// consume the API for both the OSS and Enterprise versions of Sturdy (via for example dynamic @include and @skip GraphQL
// annotations in queries).
//
// The APIs defined in enterpriseSchema and cloudSchema have multiple implementations of the resolvers. A stubbed
// implementation that's included in the OSS version, and a Enterprise implementation that's only included in Enterprise
// builds.
//
// The full schema is licensed under the same license as the OSS version of Sturdy.
var String = strings.Join([]string{
	schema,
	ossSchema,
	enterpriseSchema,
	cloudSchema,
}, "\n")

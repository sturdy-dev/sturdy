package schema

import (
	"strings"
)

var String = strings.Join([]string{
	schema,
	ossSchema,
	enterpriseSchema,
}, "\n")

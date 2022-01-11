//go:build enterprise
// +build enterprise

package enterprise

import (
	"mash/pkg/di"
	"mash/pkg/github/enterprise/workers"
)

var Module = di.NewModule(
	di.Needs(workers.Module),
)

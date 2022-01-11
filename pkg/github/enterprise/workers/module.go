//go:build enterprise
// +build enterprise

package workers

import (
	"mash/pkg/di"
)

var Module = di.NewModule(
	di.Provides(NewImporterQueue),
	di.Provides(NewClonerQueue),
)

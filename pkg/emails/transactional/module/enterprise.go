//go:build enterprise
// +build enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/emails/transactional"
	"mash/pkg/emails/transactional/enterprise"
)

var Module = di.NewModule(
	di.Provides(transactional.New),
	di.Provides(enterprise.New, new(transactional.EmailSender)),
)

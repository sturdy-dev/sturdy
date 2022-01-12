//go:build !enterprise
// +build !enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/emails/transactional"
)

var Module = di.NewModule(
	di.Provides(transactional.New, new(transactional.EmailSender)),
)

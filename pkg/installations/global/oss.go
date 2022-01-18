//go:build !enterprise
// +build !enterprise

package global

import "mash/pkg/installations"

var installationType = installations.TypeOSS

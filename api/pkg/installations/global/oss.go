//go:build !enterprise
// +build !enterprise

package global

import "getsturdy.com/api/pkg/installations"

var installationType = installations.TypeOSS

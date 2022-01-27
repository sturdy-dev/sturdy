package installations

import (
	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/version"
)

type Type uint

// Installation represents a selfhosted installation of Sturdy.
type Installation struct {
	ID         string                   `db:"id"`
	Type       version.DistributionType `db:"-"`
	Version    string                   `db:"-"`
	LicenseKey *string                  `db:"license_key"`
	License    *licenses.License        `db:"-"`
}

package validations

import (
	"time"

	"getsturdy.com/api/pkg/licenses"
)

type Validation struct {
	ID        string          `db:"id"`
	LicenseID licenses.ID     `db:"license_id"`
	Timestamp time.Time       `db:"timestamp"`
	Status    licenses.Status `db:"status"`
}

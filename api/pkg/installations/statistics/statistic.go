package statistics

import (
	"time"
)

type Statistic struct {
	InstallationID string    `json:"installation_id" db:"installation_id" validate:"required"`
	LicenseKey     *string   `json:"license_key" db:"license_key"`
	Version        string    `json:"version" db:"version" validate:"required"`
	IP             *string   `json:"-" db:"ip"`
	RecordedAt     time.Time `json:"recorded_at" db:"recorded_at"`
	ReceivedAt     time.Time `json:"-" db:"received_at"`
	UsersCount     uint64    `json:"users_count" db:"users_count"`
	CodebasesCount uint64    `json:"codebases_count" db:"codebases_count"`
	FirstUserEmail *string   `json:"first_user_email" db:"first_user_email"`
}

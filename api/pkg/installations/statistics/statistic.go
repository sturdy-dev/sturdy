package statistics

import (
	"net"
	"time"
)

type Statistic struct {
	InstallationID string    `json:"installation_id" db:"installation_id"`
	LicenseKey     *string   `json:"license_key" db:"license_key"`
	Version        string    `json:"version" db:"version"`
	IP             *net.IP   `json:"-" db:"ip"`
	RecordedAt     time.Time `json:"recorded_at" db:"recorded_at"`
	ReceivedAt     time.Time `json:"-" db:"received_at"`
	UsersCount     uint64    `json:"users_count" db:"users_count"`
	CodebasesCount uint64    `json:"codebases_count" db:"codebases_count"`
}

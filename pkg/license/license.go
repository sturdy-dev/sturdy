package license

import (
	"time"
)

type SelfHostedLicense struct {
	ID                  string    `db:"id"`
	CloudOrganizationID string    `db:"cloud_organization_id"`
	Seats               int       `db:"seats"`
	CreatedAt           time.Time `db:"created_at"`
	Active              bool      `db:"active"`
}

type SelfHostedLicenseValidation struct {
	ID                    string    `db:"id"`
	SelfHostedLicenseID   string    `db:"self_hosted_license_id"`
	ValidatedAt           time.Time `db:"validated_at"`
	Status                bool      `db:"status"`
	ReportedVersion       string    `db:"reported_version"`
	ReportedBootedAt      time.Time `db:"reported_booted_at"`
	ReportedUserCount     int       `db:"reported_user_count"`
	ReportedCodebaseCount int       `db:"reported_codebase_count"`
	FromIPAddr            string    `db:"from_ip_addr"`
}

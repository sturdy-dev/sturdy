package graphql

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"mash/pkg/graphql/resolvers"
	service_license "mash/pkg/license/enterprise/service"
)

type licenseValidation struct {
	status  error
	message string
}

func (l *licenseValidation) ID() graphql.ID {
	return "license"
}

func (l *licenseValidation) Status() resolvers.LicenseValidationStatus {
	switch {
	case l.status == nil:
		return resolvers.LicenseValidationStatusOk
	case errors.Is(l.status, service_license.ErrExpired):
		return resolvers.LicenseValidationStatusExpired
	case errors.Is(l.status, service_license.ErrTooManyUsers):
		return resolvers.LicenseValidationStatusInvalid
	default:
		return resolvers.LicenseValidationStatusUnknown
	}
}

func (l *licenseValidation) Message() *string {
	switch {
	case l.status == nil:
		return nil
	case errors.Is(l.status, service_license.ErrExpired):
		return str("This license has expired, visit getsturdy.com to renew it.")
	case errors.Is(l.status, service_license.ErrTooManyUsers):
		return str("This license is not valid for this many seats. Visit getsturdy.com to upgrade to a new license.")
	default:
		return str("This license is could not be validated. Contact support@getsturdy.com for more information.")
	}
}

func str(s string) *string {
	return &s
}

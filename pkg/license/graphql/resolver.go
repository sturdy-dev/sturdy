package graphql

import (
	"context"
	"errors"
	"time"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"

	"mash/pkg/graphql/resolvers"
	"mash/pkg/ip"
	"mash/pkg/license"
	service_license "mash/pkg/license/service"
)

type licenseRootResovler struct {
	service *service_license.Service
	logger  *zap.Logger
}

func New(
	service *service_license.Service,
	logger *zap.Logger,
) resolvers.LicenseRootResolver {
	return &licenseRootResovler{
		service: service,
		logger:  logger,
	}
}

func (r *licenseRootResovler) ValidateLicense(ctx context.Context, args resolvers.ValidateLicenseArgs) (resolvers.LicenseValidation, error) {
	val := license.SelfHostedLicenseValidation{
		// Populate with user provided values
		ReportedUserCount:     int(args.Input.UserCount),
		ReportedCodebaseCount: int(args.Input.CodebaseCount),
		ReportedBootedAt:      time.Unix(int64(args.Input.BootedAt), 0),
		ReportedVersion:       args.Input.Version,
	}

	if remoteIP, ok := ip.FromContext(ctx); ok {
		val.FromIPAddr = remoteIP.String()
	}

	status := r.service.Validate(ctx, args.Input.Key, val)
	if status != nil {
		r.logger.Warn("failed to validate license", zap.Error(status), zap.Any("input", args.Input))
	}

	return &licenseValidation{status: status}, nil
}

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

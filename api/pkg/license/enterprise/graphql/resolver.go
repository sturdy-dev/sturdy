package graphql

import (
	"context"
	"time"

	"go.uber.org/zap"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/ip"
	"getsturdy.com/api/pkg/license/enterprise/license"
	service_license "getsturdy.com/api/pkg/license/enterprise/service"
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

func (r *licenseRootResovler) InternalListForOrganization(ctx context.Context, id string) ([]resolvers.LicenseResolver, error) {
	licenses, err := r.service.ListOrganizationSubscriptions(ctx, id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.LicenseResolver
	for _, l := range licenses {
		res = append(res, &licenseResolver{license: l})
	}
	return res, nil
}

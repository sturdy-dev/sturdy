package graphql

import (
	"context"
	"fmt"
	"sync"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/licenses"

	"github.com/graph-gophers/graphql-go"
)

type licenseResolver struct {
	root    *licenseRootResovler
	license *licenses.License

	validateOnce  sync.Once
	validateError error
}

func (r *licenseResolver) ID() graphql.ID {
	return graphql.ID(r.license.ID)
}

func (r *licenseResolver) Key() string {
	return r.license.Key
}

func (r *licenseResolver) Status(ctx context.Context) (resolvers.LicenseStatus, error) {
	r.validateOnce.Do(func() {
		r.validateError = r.root.svc.Validate(ctx, r.license)
	})
	if r.validateError != nil {
		return resolvers.LicenseStatusUnknown, r.validateError
	}
	switch r.license.Status {
	case licenses.StatusValid:
		return resolvers.LicenseStatusValid, nil
	case licenses.StatusInvalid:
		return resolvers.LicenseStatusInvalid, nil
	default:
		return resolvers.LicenseStatusUnknown, fmt.Errorf("unknown status: %s", r.license.Status)
	}
}

func (r *licenseResolver) Messages(ctx context.Context) ([]resolvers.LicenseMessageResolver, error) {
	r.validateOnce.Do(func() {
		r.validateError = r.root.svc.Validate(ctx, r.license)
	})
	if r.validateError != nil {
		return nil, gqlerrors.Error(r.validateError)
	}

	rr := make([]resolvers.LicenseMessageResolver, 0, len(r.license.Messages))
	for _, message := range r.license.Messages {
		rr = append(rr, &messageResolver{
			msg: message,
		})
	}

	return rr, nil
}

func (r *licenseResolver) CreatedAt() (int32, error) {
	return int32(r.license.CreatedAt.Unix()), nil
}

func (r *licenseResolver) ExpiresAt() (int32, error) {
	return int32(r.license.ExpiresAt.Unix()), nil
}

func (r *licenseResolver) Seats() int32 {
	return int32(r.license.Seats)
}

type messageResolver struct {
	msg *licenses.Message
}

func (r *messageResolver) Text() string {
	return r.msg.Text
}

func (r *messageResolver) Type() resolvers.LicenseMessageType {
	switch r.msg.Type {
	case licenses.TypeBanner:
		return resolvers.LicenseMessageTypeBanner
	case licenses.TypeFullscreen:
		return resolvers.LicenseMessageTypeFullscreen
	case licenses.TypeNotification:
		return resolvers.LicenseMessageTypeNotification
	default:
		return resolvers.LicenseMessageTypeUnknown
	}
}

func (r *messageResolver) Level() resolvers.LicenseMessageLevel {
	switch r.msg.Level {
	case licenses.LevelInfo:
		return resolvers.LicenseMessageLevelInfo
	case licenses.LevelWarning:
		return resolvers.LicenseMessageLevelWarning
	case licenses.LevelError:
		return resolvers.LicenseMessageLevelError
	default:
		return resolvers.LicenseMessageLevelUnknown
	}
}

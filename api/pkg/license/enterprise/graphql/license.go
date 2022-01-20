package graphql

import (
	"time"

	"github.com/graph-gophers/graphql-go"

	"mash/pkg/license/enterprise/license"
)

type licenseResolver struct {
	license *license.SelfHostedLicense
}

func (l *licenseResolver) ID() graphql.ID {
	return graphql.ID(l.license.ID)
}

func (l *licenseResolver) Seats() int32 {
	return int32(l.license.Seats)
}

func (l *licenseResolver) UsedSeats() int32 {
	// TODO
	return 0
}

func (l *licenseResolver) ExpiresAt() int32 {
	return int32(l.license.CreatedAt.Add(time.Hour * 24 * 365).Unix())
}

func (l *licenseResolver) LicenseKey() string {
	return l.license.ID
}

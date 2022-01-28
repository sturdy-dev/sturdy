package service

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/installations/statistics"
)

type Service struct {
	installation *installations.Installation
}

func New(
	installation *installations.Installation,
) *Service {
	return &Service{
		installation: installation,
	}
}

func (s *Service) Get(ctx context.Context) (*statistics.Statistic, error) {
	stat := &statistics.Statistic{
		InstallationID: s.installation.ID,
		LicenseKey:     s.installation.LicenseKey,
		Version:        s.installation.Version,
		RecordedAt:     time.Now(),
	}
	return stat, nil
}

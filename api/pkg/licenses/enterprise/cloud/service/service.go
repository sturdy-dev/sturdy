package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	service_installation_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/db"
	service_license_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"

	"github.com/mergestat/timediff"
)

type Service struct {
	db db.Repository

	statisticsService  *service_installation_statistics.Service
	validationsService *service_license_validations.Service
}

func NewService(
	db db.Repository,
	statisticsService *service_installation_statistics.Service,
	validationsService *service_license_validations.Service,
) *Service {
	return &Service{
		db:                 db,
		statisticsService:  statisticsService,
		validationsService: validationsService,
	}
}

func (s *Service) ValidateByKey(ctx context.Context, key string) (*licenses.License, error) {
	license, err := s.db.GetByKey(ctx, key)
	if errors.Is(err, db.ErrNotFound) {
		return &licenses.License{
			Status: licenses.StatusInvalid,
			Messages: []*licenses.Message{
				{
					Type:  licenses.TypeBanner,
					Text:  "No license found",
					Level: licenses.LevelInfo,
				},
			},
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get license: %w", err)
	}

	if err := s.Validate(ctx, license); err != nil {
		return nil, fmt.Errorf("failed to validate license: %w", err)
	}

	return license, nil
}

func (s *Service) ListByOrganizationID(ctx context.Context, orgID string) ([]*licenses.License, error) {
	ll, err := s.db.ListByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list licenses: %w", err)
	}
	for _, l := range ll {
		if err := s.Validate(ctx, l); err != nil {
			return nil, fmt.Errorf("failed to validate license: %w", err)
		}
	}
	return ll, nil
}

func (s *Service) Validate(ctx context.Context, license *licenses.License) error {
	status, messages, err := s.validate(ctx, license)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	license.Status = status
	license.Messages = messages
	return nil
}

var (
	oneDay = time.Hour * 24

	expiryLeeway = 3 * oneDay

	// TODO: Start using
	// seatsLeeway        uint64 = 5

	statisticsLeeway   = 3 * time.Hour
	statisticsDeadline = oneDay
)

func (s *Service) validate(ctx context.Context, license *licenses.License) (licenses.Status, []*licenses.Message, error) {
	if license.ExpiresAt.Before(time.Now()) {
		return licenses.StatusInvalid, []*licenses.Message{
			{
				Type:  licenses.TypeFullscreen,
				Level: licenses.LevelError,
				Text:  fmt.Sprintf("The license expired %s", timediff.TimeDiff(license.ExpiresAt)),
			},
		}, nil
	}

	messages := []*licenses.Message{}
	untilExpiration := time.Until(license.ExpiresAt)
	if untilExpiration < expiryLeeway {
		messages = append(messages, &licenses.Message{
			Type:  licenses.TypeBanner,
			Level: licenses.LevelWarning,
			Text:  fmt.Sprintf("The license will expire in %s", timediff.TimeDiff(license.ExpiresAt)),
		})
	}

	validations, err := s.validationsService.ListLatest(ctx, license.ID)
	if err != nil {
		return licenses.StatusUnknown, nil, fmt.Errorf("failed to list validations: %w", err)
	}

	statistics, err := s.statisticsService.GetByLicenseKey(ctx, license.Key)
	if errors.Is(err, service_installation_statistics.ErrNotFound) {
		if len(validations) == 0 {
			return licenses.StatusValid, []*licenses.Message{
				{
					Type:  licenses.TypeBanner,
					Level: licenses.LevelWarning,
					Text:  "Sturdy hasn't heard from this installation yet. Please make sure that the installation can connect to the cloud, otherwise, this license will be considered invalid.",
				},
			}, nil
		} else if len(validations) > 3 {
			return licenses.StatusInvalid, []*licenses.Message{
				{
					Type:  licenses.TypeBanner,
					Level: licenses.LevelWarning,
					Text:  "Sturdy has not heard from this license yet.",
				},
			}, nil

		} else {
			return licenses.StatusValid, messages, nil
		}
	} else if err != nil {
		return licenses.StatusUnknown, nil, fmt.Errorf("failed to get license statistics: %w", err)
	}

	if statistics.RecordedAt.Add(statisticsLeeway).Before(time.Now()) {
		messages = append(messages, &licenses.Message{
			Type:  licenses.TypeBanner,
			Level: licenses.LevelWarning,
			Text:  "Sturdy has not heard from this installation in the last 3 hours. Please make sure that the installation sends statistics, otherwise, this license will be considered invalid.",
		})
	}

	if statistics.RecordedAt.Add(statisticsDeadline).Before(time.Now()) {
		return licenses.StatusInvalid, append(messages, &licenses.Message{
			Type:  licenses.TypeFullscreen,
			Level: licenses.LevelError,
			Text:  "Sturdy has not heard from this installation in the last 24 hours.",
		}), nil
	}

	if license.Seats < statistics.UsersCount {
		return licenses.StatusInvalid, []*licenses.Message{
			{
				Type:  licenses.TypeBanner,
				Level: licenses.LevelError,
				Text:  fmt.Sprintf("Maximum number of users exceeded. %d users are allowed, but %d users are currently using this license.", license.Seats, statistics.UsersCount),
			},
		}, nil
	}

	return licenses.StatusValid, messages, nil
}

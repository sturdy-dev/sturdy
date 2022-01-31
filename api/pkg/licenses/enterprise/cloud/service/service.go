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

	expiryLeeway              = 3 * oneDay
	seatsLeeway        uint64 = 5
	statisticsLeeway          = 3 * time.Hour
	statisticsDeadline        = oneDay
)

func (s *Service) validate(ctx context.Context, license *licenses.License) (licenses.Status, []*licenses.Message, error) {
	if license.ExpiresAt.Before(time.Now()) {
		return licenses.StatusInvalid, []*licenses.Message{
			{
				Type:  licenses.TypeBanner,
				Level: licenses.LevelError,
				Text:  "License has expired",
			},
		}, nil
	}

	messages := []*licenses.Message{}
	untilExpiration := time.Until(license.ExpiresAt)
	if untilExpiration < expiryLeeway {
		messages = append(messages, &licenses.Message{
			Type:  licenses.TypeBanner,
			Level: licenses.LevelWarning,
			Text:  "License expires in less than three days",
		})
	}

	validations, err := s.validationsService.ListLatest(ctx, license.ID)
	if err != nil {
		return licenses.StatusUnknown, nil, fmt.Errorf("failed to list validations: %w", err)
	}

	statistics, err := s.statisticsService.GetByLicenseKey(ctx, license.Key)
	if errors.Is(err, service_installation_statistics.ErrNotFound) {
		if len(validations) >= 1 {
			return licenses.StatusInvalid, append(messages, &licenses.Message{
				Type:  licenses.TypeBanner,
				Level: licenses.LevelWarning,
				Text:  "We didn't receive any statistics for this license yet. License is considered invalid, until we receive statistics.",
			}), nil
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
			Text:  "We didn't receive any statistics for this license in the last 3 hours.",
		})
	}

	if statistics.RecordedAt.Add(statisticsDeadline).Before(time.Now()) {
		return licenses.StatusInvalid, append(messages, &licenses.Message{
			Type:  licenses.TypeBanner,
			Level: licenses.LevelError,
			Text:  "We didn't receive any statistics for this license in the last 24 hours.",
		}), nil
	}

	if license.Seats < statistics.UsersCount {
		return licenses.StatusInvalid, []*licenses.Message{
			{
				Type:  licenses.TypeBanner,
				Level: licenses.LevelError,
				Text:  "Maximum number of users exceeded",
			},
		}, nil
	}

	return licenses.StatusValid, messages, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/db"
)

type Service struct {
	db db.Repository
}

func NewService(
	db db.Repository,
) *Service {
	return &Service{
		db: db,
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
	oneDay       = time.Hour * 24
	expiryLeeway = 3 * oneDay
)

func (s *Service) validate(ctx context.Context, license *licenses.License) (licenses.Status, []*licenses.Message, error) {
	if license.ExpiresAt.Before(time.Now()) {
		return licenses.StatusInvalid, []*licenses.Message{
			{
				Type:  licenses.TypeBanner,
				Level: licenses.LevelError,
				Text:  "license expired",
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

	// TODO: more validations based on statistics

	return licenses.StatusValid, messages, nil
}

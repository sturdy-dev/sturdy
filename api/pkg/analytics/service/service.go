package service

import (
	"context"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/users"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	client analytics.Client
}

func New(
	logger *zap.Logger,
	client analytics.Client,
) *Service {
	return &Service{
		logger: logger,
		client: client,
	}
}

func (s *Service) IdentifyOrganization(ctx context.Context, org *organization.Organization) {
	if err := s.client.Enqueue(analytics.GroupIdentify{
		Type: "organization", // this should match other event's property key
		Key:  org.ID,
		Properties: map[string]interface{}{
			"name": org.Name,
		},
	}); err != nil {
		s.logger.Error("failed to identify codebase", zap.Error(err))
	}
}

func (s *Service) IdentifyCodebase(ctx context.Context, cb *codebase.Codebase) {
	if err := s.client.Enqueue(analytics.GroupIdentify{
		Type: "codebase_id", // this should match other event's property key
		Key:  cb.ID,
		Properties: map[string]interface{}{
			"name":      cb.Name,
			"is_public": cb.IsPublic,
		},
	}); err != nil {
		s.logger.Error("failed to identify codebase", zap.Error(err))
	}
}

func (s *Service) IdentifyUser(ctx context.Context, user *users.User) {
	if err := s.client.Enqueue(analytics.Identify{
		DistinctId: user.ID,
		Properties: analytics.NewProperties().
			Set("name", user.Name).
			Set("email", user.Email),
	}); err != nil {
		s.logger.Error("failed to identify user", zap.Error(err))
	}
}

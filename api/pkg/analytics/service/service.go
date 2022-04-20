package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/users"

	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	client posthog.Client
}

func New(
	logger *zap.Logger,
	client posthog.Client,
) *Service {
	return &Service{
		logger: logger,
		client: client,
	}
}

func (s *Service) Capture(ctx context.Context, event string, oo ...analytics.CaptureOption) {
	userID, _ := auth.UserID(ctx)
	s.CaptureUser(userID, event, oo...)
}

func (s *Service) CaptureUser(userID users.ID, event string, oo ...analytics.CaptureOption) {
	options := &analytics.CaptureOptions{}
	for _, o := range oo {
		o(options)
	}

	capture := posthog.Capture{
		DistinctId: userID.String(),
		Properties: options.Properties,
		Event:      event,
		Groups:     options.Groups,
	}

	if err := s.client.Enqueue(capture); err != nil {
		s.logger.Error("failed to capture user event", zap.Error(err), zap.Any("capture", capture))
	}
}

func (s *Service) IdentifyOrganization(ctx context.Context, org *organization.Organization) {
	if err := s.client.Enqueue(posthog.GroupIdentify{
		Type: "organization", // this should match other event's property key
		Key:  org.ID,
		Properties: map[string]any{
			"name": org.Name,
		},
	}); err != nil {
		s.logger.Error("failed to identify codebase", zap.Error(err))
	}
}

func (s *Service) IdentifyCodebase(ctx context.Context, cb *codebases.Codebase) {
	if err := s.client.Enqueue(posthog.GroupIdentify{
		Type: "codebase",
		Key:  cb.ID.String(),
		Properties: map[string]any{
			"name":      cb.Name,
			"is_public": cb.IsPublic,
		},
	}); err != nil {
		s.logger.Error("failed to identify codebase", zap.Error(err))
	}
}

func (s *Service) IdentifyUser(ctx context.Context, user *users.User) {
	id := user.ID
	if user.Is != nil {
		id = *user.Is
	}
	if err := s.client.Enqueue(posthog.Identify{
		DistinctId: id.String(),
		Properties: map[string]any{
			"name":   user.Name,
			"email":  user.Email,
			"status": user.Status,
		},
	}); err != nil {
		s.logger.Error("failed to identify user", zap.Error(err))
	}
}

func (s *Service) IdentifyGitHubInstallation(ctx context.Context, installationID int64, accountLogin, accountEmail string) {
	if err := s.client.Enqueue(posthog.Identify{
		DistinctId: fmt.Sprintf("%d", installationID), // Using the installation ID as a person?
		Properties: map[string]any{
			"installation_org":        accountLogin,
			"email":                   accountEmail,
			"github_app_installation": true,
		},
	}); err != nil {
		s.logger.Error("failed to identify github installation", zap.Error(err))
	}
}

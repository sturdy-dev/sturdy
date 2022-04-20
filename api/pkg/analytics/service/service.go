package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/version"

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
	s.CaptureUser(ctx, userID, event, oo...)
}

func (s *Service) CaptureUser(ctx context.Context, userID users.ID, event string, oo ...analytics.CaptureOption) {
	options := &analytics.CaptureOptions{Properties: map[string]interface{}{}}
	for _, o := range oo {
		o(options)
	}

	capture := posthog.Capture{
		DistinctId: userID.String(),
		Properties: options.Properties,
		Event:      event,
		Groups:     options.Groups,
	}

	s.capture(ctx, capture)
}

func (s *Service) IdentifyOrganization(ctx context.Context, org *organization.Organization) {
	s.groupIdentify(ctx, posthog.GroupIdentify{
		Type: "organization", // this should match other event's property key
		Key:  org.ID,
		Properties: map[string]any{
			"name": org.Name,
		},
	})
}

func (s *Service) IdentifyCodebase(ctx context.Context, cb *codebases.Codebase) {
	s.groupIdentify(ctx, posthog.GroupIdentify{
		Type: "codebase",
		Key:  cb.ID.String(),
		Properties: map[string]any{
			"name":      cb.Name,
			"is_public": cb.IsPublic,
		},
	})
}

func (s *Service) IdentifyUser(ctx context.Context, user *users.User) {
	id := user.ID
	if user.Is != nil {
		id = *user.Is
	}
	s.identify(ctx, posthog.Identify{
		DistinctId: id.String(),
		Properties: map[string]any{
			"name":   user.Name,
			"email":  user.Email,
			"status": user.Status,
		},
	})
}

func (s *Service) IdentifyGitHubInstallation(ctx context.Context, installationID int64, accountLogin, accountEmail string) {
	s.identify(ctx, posthog.Identify{
		DistinctId: fmt.Sprintf("%d", installationID), // Using the installation ID as a person?
		Properties: map[string]any{
			"installation_org":        accountLogin,
			"email":                   accountEmail,
			"github_app_installation": true,
		},
	})
}

func (s *Service) groupIdentify(ctx context.Context, groupIdentify posthog.GroupIdentify) {
	groupIdentify.Properties.Set("environment", version.Type.String())
	if err := s.client.Enqueue(groupIdentify); err != nil {
		s.logger.Error("failed to group identify", zap.Error(err), zap.Any("groupIdentify", groupIdentify))
	}
}

func (s *Service) identify(ctx context.Context, identify posthog.Identify) {
	identify.Properties.Set("environment", version.Type.String())
	if err := s.client.Enqueue(identify); err != nil {
		s.logger.Error("failed to identify", zap.Error(err), zap.Any("identify", identify))
	}
}

func (s *Service) capture(ctx context.Context, capture posthog.Capture) {
	capture.Properties.Set("environment", version.Type.String())
	if err := s.client.Enqueue(capture); err != nil {
		s.logger.Error("failed to capture", zap.Error(err), zap.Any("capture", capture))
	}
}

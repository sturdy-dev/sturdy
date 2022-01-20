package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"mash/pkg/analytics"
	"mash/pkg/auth"
	"mash/pkg/codebase"
	db_codebase "mash/pkg/codebase/db"
	"mash/pkg/codebase/vcs"
	"mash/pkg/shortid"
	"mash/pkg/view/events"
	service_workspace "mash/pkg/workspace/service"
	"mash/vcs/executor"
	"mash/vcs/provider"
)

type Service struct {
	repo             db_codebase.CodebaseRepository
	codebaseUserRepo db_codebase.CodebaseUserRepository

	workspaceService service_workspace.Service

	logger           *zap.Logger
	executorProvider executor.Provider
	analyticsClient  analytics.Client
	eventsSender     events.EventSender
}

func New(
	repo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,

	workspaceService service_workspace.Service,

	logger *zap.Logger,
	executorProvider executor.Provider,
	analyticsClient analytics.Client,
	eventsSender events.EventSender,
) *Service {
	return &Service{
		repo:             repo,
		codebaseUserRepo: codebaseUserRepo,

		workspaceService: workspaceService,

		logger:           logger,
		executorProvider: executorProvider,
		analyticsClient:  analyticsClient,
		eventsSender:     eventsSender,
	}
}

func (s *Service) GetByID(_ context.Context, id string) (*codebase.Codebase, error) {
	return s.repo.Get(id)
}

func (s *Service) GetByShortID(_ context.Context, shortID string) (*codebase.Codebase, error) {
	return s.repo.GetByShortID(shortID)
}

func (s *Service) CanAccess(_ context.Context, userID string, codebaseID string) (bool, error) {
	_, err := s.codebaseUserRepo.GetByUserAndCodebase(userID, codebaseID)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("failed to check user %s access to codebase %s: %w", userID, codebaseID, err)
	}
}

func (s *Service) ListByOrganization(ctx context.Context, organizationID string) ([]*codebase.Codebase, error) {
	res, err := s.repo.ListByOrganization(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("could not ListByOrganization: %w", err)
	}
	return res, nil
}

func (svc *Service) Create(ctx context.Context, name string, organizationID *string) (*codebase.Codebase, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	codebaseID := uuid.NewString()
	t := time.Now()

	cb := codebase.Codebase{
		ID:              codebaseID,
		ShortCodebaseID: codebase.ShortCodebaseID(shortid.New()),
		Name:            name,
		Description:     "",
		Emoji:           "",
		CreatedAt:       &t,
		IsReady:         true,           // No additional setup needed
		OrganizationID:  organizationID, // Optional
	}

	// Create codebase in database
	if err := svc.repo.Create(cb); err != nil {
		return nil, fmt.Errorf("failed to create codebase: %w", err)
	}

	if err := svc.executorProvider.New().
		AllowRebasingState(). // allowed because the repo does not exist yet
		Schedule(func(trunkProvider provider.RepoProvider) error {
			return vcs.Create(trunkProvider, cb.ID)
		}).ExecTrunk(cb.ID, "createCodebase"); err != nil {

		return nil, fmt.Errorf("failed to create codebase on disk: %w", err)
	}

	// Add user
	err = svc.codebaseUserRepo.Create(codebase.CodebaseUser{
		ID:         uuid.New().String(),
		UserID:     userID,
		CodebaseID: cb.ID,
		CreatedAt:  &t,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	err = svc.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: userID,
		Event:      "create codebase",
		Properties: analytics.NewProperties().
			Set("codebase_id", cb.ID).
			Set("name", cb.Name),
	})
	if err != nil {
		svc.logger.Error("analytics failed", zap.Error(err))
	}

	if err := svc.workspaceService.CreateWelcomeWorkspace(cb.ID, userID, cb.Name); err != nil {
		svc.logger.Error("failed to create welcome workspace", zap.Error(err))
		// not a critical error, continue
	}

	// Send events
	if err := svc.eventsSender.Codebase(cb.ID, events.CodebaseUpdated, cb.ID); err != nil {
		return nil, fmt.Errorf("failed to send events: %w", err)
	}

	return &cb, nil
}

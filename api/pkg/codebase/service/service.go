package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebase"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/codebase/vcs"
	"getsturdy.com/api/pkg/shortid"
	"getsturdy.com/api/pkg/events"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
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

func (s *Service) ListByOrganizationAndUser(ctx context.Context, organizationID, userID string) ([]*codebase.Codebase, error) {
	codebases, err := s.repo.ListByOrganization(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("could not ListByOrganization: %w", err)
	}

	var res []*codebase.Codebase

	for _, cb := range codebases {
		_, err := s.codebaseUserRepo.GetByUserAndCodebase(userID, cb.ID)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			continue
		case err != nil:
			return nil, fmt.Errorf("could not codebase user: %w", err)
		case err == nil:
			res = append(res, cb)
		}
	}

	return res, nil
}

// ListOrgsByUser returns a list of organization IDs that the user can _see_ through it's explicit membership
// of one of it's codebases.
func (svc *Service) ListOrgsByUser(ctx context.Context, userID string) ([]string, error) {
	orgIDs, err := svc.orgsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []string
	for k := range orgIDs {
		res = append(res, k)
	}

	return res, nil
}

func (svc *Service) UserIsMemberOfCodebaseInOrganization(ctx context.Context, userID, organizationID string) (bool, error) {
	orgIDs, err := svc.orgsByUser(ctx, userID)
	if err != nil {
		return false, err
	}

	_, ok := orgIDs[organizationID]
	return ok, nil
}

func (svc *Service) orgsByUser(ctx context.Context, userID string) (map[string]struct{}, error) {
	codebaseUsers, err := svc.codebaseUserRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("could not ListByUser: %w", err)
	}

	orgIDs := make(map[string]struct{})

	for _, cu := range codebaseUsers {
		cb, err := svc.repo.Get(cu.CodebaseID)
		if err != nil {
			return nil, fmt.Errorf("could not get codebase: %w", err)
		}
		if cb.OrganizationID != nil {
			orgIDs[*cb.OrganizationID] = struct{}{}
		}
	}

	return orgIDs, nil
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

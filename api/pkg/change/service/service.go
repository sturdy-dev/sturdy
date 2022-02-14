package service

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	acl_provider "getsturdy.com/api/pkg/codebase/acl/provider"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/workspace"

	"github.com/google/uuid"
)

type Service struct {
	aclProvider *acl_provider.Provider
	userRepo    db_user.Repository
	changeRepo  db_change.Repository
}

func New(
	aclProvider *acl_provider.Provider,
	userRepo db_user.Repository,
	changeRepo db_change.Repository,
) *Service {
	return &Service{
		aclProvider: aclProvider,
		userRepo:    userRepo,
		changeRepo:  changeRepo,
	}
}

func (svc *Service) ListChanges(ctx context.Context, ids ...change.ID) ([]*change.Change, error) {
	return svc.changeRepo.ListByIDs(ctx, ids...)
}

func (svc *Service) GetChangeByID(ctx context.Context, id change.ID) (*change.Change, error) {
	ch, err := svc.changeRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (svc *Service) GetByCommitID(ctx context.Context, commitID, codebaseID string) (*change.Change, error) {
	ch, err := svc.changeRepo.GetByCommitID(ctx, commitID, codebaseID)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (s *Service) Create(ctx context.Context, ws *workspace.Workspace, commitID, msg string) (*change.Change, error) {
	changeID := change.ID(uuid.NewString())
	t := time.Now()
	changeChange := change.Change{
		ID:                 changeID,
		CodebaseID:         ws.CodebaseID,
		Title:              &msg,
		UpdatedDescription: ws.DraftDescription,
		UserID:             &ws.UserID,
		CreatedAt:          &t,
		CommitID:           &commitID,
	}
	if err := s.changeRepo.Insert(changeChange); err != nil {
		return nil, fmt.Errorf("failed to insert change: %w", err)
	}

	return &changeChange, nil
}

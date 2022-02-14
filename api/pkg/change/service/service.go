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
	aclProvider      *acl_provider.Provider
	userRepo         db_user.Repository
	changeRepo       db_change.Repository
	commitChangeRepo db_change.CommitRepository
}

func New(
	aclProvider *acl_provider.Provider,
	userRepo db_user.Repository,
	changeRepo db_change.Repository,
	commitChangeRepo db_change.CommitRepository,
) *Service {
	return &Service{
		aclProvider:      aclProvider,
		userRepo:         userRepo,
		changeRepo:       changeRepo,
		commitChangeRepo: commitChangeRepo,
	}
}

func (svc *Service) ListChangeCommits(ctx context.Context, ids ...change.ID) ([]*change.ChangeCommit, error) {
	return svc.commitChangeRepo.ListByChangeIDs(ctx, ids...)
}

func (svc *Service) GetChangeCommitByCommitIDAndCodebaseID(ctx context.Context, commitID, codebaseID string) (*change.ChangeCommit, error) {
	changeCommit, err := svc.commitChangeRepo.GetByCommitID(commitID, codebaseID)
	if err != nil {
		return nil, err
	}
	return &changeCommit, nil
}

func (svc *Service) GetChangeByID(ctx context.Context, id change.ID) (*change.Change, error) {
	ch, err := svc.changeRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (svc *Service) GetChangeCommitOnTrunkByChangeID(ctx context.Context, id change.ID) (*change.ChangeCommit, error) {
	ch, err := svc.commitChangeRepo.GetByChangeIDOnTrunk(id)
	if err != nil {
		return nil, err
	}
	return &ch, nil
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
	}
	if err := s.changeRepo.Insert(changeChange); err != nil {
		return nil, fmt.Errorf("failed to insert change: %w", err)
	}

	if err := s.commitChangeRepo.Insert(change.ChangeCommit{
		ChangeID:   changeID,
		CommitID:   commitID,
		CodebaseID: ws.CodebaseID,
		Trunk:      true,
	}); err != nil {
		return nil, fmt.Errorf("failed to insert change commit: %w", err)
	}

	return &changeChange, nil
}

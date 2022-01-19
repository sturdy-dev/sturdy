package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mash/pkg/codebase"
	db_codebase "mash/pkg/codebase/db"
)

type Service struct {
	repo             db_codebase.CodebaseRepository
	codebaseUserRepo db_codebase.CodebaseUserRepository
}

func New(
	repo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
) *Service {
	return &Service{
		repo:             repo,
		codebaseUserRepo: codebaseUserRepo,
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

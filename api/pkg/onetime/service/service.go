package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"getsturdy.com/api/pkg/onetime"
	db_onetime "getsturdy.com/api/pkg/onetime/db"
	"getsturdy.com/api/pkg/users"
)

type Service struct {
	repo db_onetime.Repository
}

func New(repo db_onetime.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateToken(ctx context.Context, userID string) (*onetime.Token, error) {
	token := onetime.New(userID)

	if err := s.repo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	return token, nil
}

var (
	ErrExpired = fmt.Errorf("token expired")
	ErrReused  = fmt.Errorf("token reused")
	ErrInvalid = fmt.Errorf("token invalid")
)

func (s *Service) Resolve(ctx context.Context, user *users.User, key string) (*onetime.Token, error) {
	key = strings.ToUpper(key)

	token, err := s.repo.Get(ctx, user.ID, key)
	switch {
	case err == nil:
		if token.IsExpired() {
			return nil, ErrExpired
		}
		if token.IsReused() {
			return nil, ErrReused
		}

		token.Clicks++
		if err := s.repo.Update(ctx, token); err != nil {
			return nil, fmt.Errorf("failed to update token: %w", err)
		}

		return token, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrInvalid
	default:
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
}

package service

import (
	"context"
	"fmt"
	"time"

	"mash/pkg/servicetokens"
	db_servicetokens "mash/pkg/servicetokens/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = bcrypt.DefaultCost
)

type Service struct {
	repo db_servicetokens.Repository
}

func New(
	repo db_servicetokens.Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new service token. It returns the created service token in plaintext (now stored) and the
// servicetoken in encrypted form (can be recovered from the database).
func (s *Service) Create(ctx context.Context, codebaseID string, name string) (string, *servicetokens.Token, error) {
	plainTextToken := uuid.New().String()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(plainTextToken), bcryptCost)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token has: %w", err)
	}

	token := &servicetokens.Token{
		ID:         uuid.NewString(),
		CodebaseID: codebaseID,
		Hash:       hashedToken,
		Name:       name,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, token); err != nil {
		return "", nil, fmt.Errorf("failed to create: %w", err)
	}

	return plainTextToken, token, nil
}

// todo: updated last used

func (s *Service) Get(ctx context.Context, id string) (*servicetokens.Token, error) {
	return s.repo.GetByID(ctx, id)
}

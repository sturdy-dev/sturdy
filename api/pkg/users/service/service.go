package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrExceeded = fmt.Errorf("maximum number of users exceeded")
)

type UserSerice struct {
	logger          *zap.Logger
	userRepo        db_user.Repository
	analyticsClient analytics.Client
}

type Service interface {
	CreateWithPassword(ctx context.Context, name, password, email string) (*users.User, error)
	GetByIDs(ctx context.Context, ids ...string) ([]*users.User, error)
	GetByID(_ context.Context, id string) (*users.User, error)
	GetByEmail(_ context.Context, email string) (*users.User, error)
	UsersCount(context.Context) (uint64, error)
	GetFirstUser(ctx context.Context) (*users.User, error)
	GetAsAuthor(ctx context.Context, userID string) (*author.Author, error)
}

func New(
	logger *zap.Logger,
	userRepo db_user.Repository,
	analyticsClient analytics.Client,
) *UserSerice {
	return &UserSerice{
		logger:          logger,
		userRepo:        userRepo,
		analyticsClient: analyticsClient,
	}
}

var (
	ErrExists   = fmt.Errorf("user already exists")
	ErrNotFound = fmt.Errorf("user not found")
)

func (s *UserSerice) GetFirstUser(ctx context.Context) (*users.User, error) {
	uu, err := s.userRepo.List(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	if len(uu) == 0 {
		return nil, ErrNotFound
	}
	return uu[0], nil
}

func (s *UserSerice) CreateWithPassword(ctx context.Context, name, password, email string) (*users.User, error) {
	if _, err := s.userRepo.GetByEmail(email); errors.Is(err, sql.ErrNoRows) {
		// all good
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	} else {
		return nil, ErrExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	t := time.Now()
	newUser := &users.User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    &t,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send events
	if err := s.analyticsClient.Enqueue(analytics.Identify{
		DistinctId: newUser.ID,
		Properties: analytics.NewProperties().
			Set("name", newUser.Name).
			Set("email", newUser.Email),
	}); err != nil {
		s.logger.Error("send to analytics failed", zap.Error(err))
	}

	if err := s.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: newUser.ID,
		Event:      "created account",
	}); err != nil {
		s.logger.Error("send to analytics failed", zap.Error(err))
	}

	return newUser, nil
}

func (s *UserSerice) GetByIDs(ctx context.Context, ids ...string) ([]*users.User, error) {
	return s.userRepo.GetByIDs(ctx, ids...)
}

func (s *UserSerice) GetByID(_ context.Context, id string) (*users.User, error) {
	return s.userRepo.Get(id)
}

func (s *UserSerice) GetByEmail(_ context.Context, email string) (*users.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserSerice) UsersCount(ctx context.Context) (uint64, error) {
	return s.userRepo.Count(ctx)
}

func (s *UserSerice) GetAsAuthor(ctx context.Context, userID string) (*author.Author, error) {
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	return &author.Author{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: emptyIfNull(user.AvatarURL),
	}, nil
}

func emptyIfNull(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

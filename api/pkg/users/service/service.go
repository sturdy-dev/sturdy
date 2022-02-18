package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrExceeded = fmt.Errorf("maximum number of users exceeded")
)

type UserService struct {
	logger           *zap.Logger
	userRepo         db_user.Repository
	analyticsService *service_analytics.Service
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
	analyticsService *service_analytics.Service,
) *UserService {
	return &UserService{
		logger:           logger,
		userRepo:         userRepo,
		analyticsService: analyticsService,
	}
}

var (
	ErrExists   = fmt.Errorf("user already exists")
	ErrNotFound = fmt.Errorf("user not found")
)

func (s *UserService) GetFirstUser(ctx context.Context) (*users.User, error) {
	uu, err := s.userRepo.List(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	if len(uu) == 0 {
		return nil, ErrNotFound
	}
	return uu[0], nil
}

func (s *UserService) CreateWithPassword(ctx context.Context, name, password, email string) (*users.User, error) {
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

	s.analyticsService.IdentifyUser(ctx, newUser)
	s.analyticsService.Capture(ctx, "created account")

	return newUser, nil
}

func (s *UserService) GetByIDs(ctx context.Context, ids ...string) ([]*users.User, error) {
	return s.userRepo.GetByIDs(ctx, ids...)
}

func (s *UserService) GetByID(_ context.Context, id string) (*users.User, error) {
	return s.userRepo.Get(id)
}

func (s *UserService) GetByEmail(_ context.Context, email string) (*users.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserService) UsersCount(ctx context.Context) (uint64, error) {
	return s.userRepo.Count(ctx)
}

func (s *UserService) GetAsAuthor(ctx context.Context, userID string) (*author.Author, error) {
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

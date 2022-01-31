package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/emails/transactional"
	"getsturdy.com/api/pkg/jwt"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_onetime "getsturdy.com/api/pkg/onetime/service"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserSerice struct {
	logger                   *zap.Logger
	userRepo                 db_user.Repository
	jwtService               *service_jwt.Service
	onetimeService           *service_onetime.Service
	transactionalEmailSender transactional.EmailSender
	analyticsClient          analytics.Client
}

type Service interface {
	CreateWithPassword(ctx context.Context, name, password, email string) (*users.User, error)
	Create(ctx context.Context, name, email string) (*users.User, error)
	SendMagicLink(ctx context.Context, user *users.User) error
	GetByIDs(ctx context.Context, ids ...string) ([]*users.User, error)
	GetByID(_ context.Context, id string) (*users.User, error)
	GetByEmail(_ context.Context, email string) (*users.User, error)
	SendEmailVerification(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, userID string, rawToken string) (*users.User, error)
	VerifyMagicLink(ctx context.Context, user *users.User, code string) error
	UsersCount(context.Context) (uint64, error)
}

func New(
	logger *zap.Logger,
	userRepo db_user.Repository,
	jwtService *service_jwt.Service,
	onetimeService *service_onetime.Service,
	transactionalEmailSender transactional.EmailSender,
	analyticsClient analytics.Client,
) *UserSerice {
	return &UserSerice{
		logger:                   logger,
		userRepo:                 userRepo,
		jwtService:               jwtService,
		onetimeService:           onetimeService,
		transactionalEmailSender: transactionalEmailSender,
		analyticsClient:          analyticsClient,
	}
}

var (
	ErrExists = fmt.Errorf("user already exists")
)

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

	// Send emails
	if err := s.transactionalEmailSender.SendWelcome(ctx, newUser); err != nil {
		s.logger.Error("failed to send welcome email", zap.Error(err))
	}

	return newUser, nil
}

func (s *UserSerice) Create(ctx context.Context, name, email string) (*users.User, error) {
	// If user already exists, send OTP
	if existingUser, err := s.userRepo.GetByEmail(email); err == nil {
		if err := s.SendMagicLink(ctx, existingUser); err != nil {
			return nil, fmt.Errorf("failed to send magic link: %w", err)
		}
		return existingUser, nil
	}

	t := time.Now()
	newUser := &users.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: &t,
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

	// Send emails
	if err := s.transactionalEmailSender.SendWelcome(ctx, newUser); err != nil {
		s.logger.Error("failed to send welcome email", zap.Error(err))
	}

	if err := s.SendMagicLink(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to send magic link: %w", err)
	}

	return newUser, nil
}

func (s *UserSerice) VerifyMagicLink(ctx context.Context, user *users.User, code string) error {
	if _, err := s.onetimeService.Resolve(ctx, user, code); err != nil {
		return fmt.Errorf("failed to resolve magic link: %w", err)
	}

	if err := s.setEmailVerified(ctx, user); err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}

	if err := s.analyticsClient.Enqueue(analytics.Identify{
		DistinctId: user.ID,
		Properties: analytics.NewProperties().
			Set("name", user.Name).
			Set("email", user.Email),
	}); err != nil {
		s.logger.Error("send to analytics failed", zap.Error(err))
	}

	if err := s.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: user.ID,
		Event:      "logged in",
		Properties: analytics.NewProperties().
			Set("type", "code"),
	}); err != nil {
		s.logger.Error("send to analytics failed", zap.Error(err))
	}

	return nil
}

func (s *UserSerice) SendMagicLink(ctx context.Context, user *users.User) error {
	token, err := s.onetimeService.CreateToken(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	code := fmt.Sprintf("%s-%s", token.Key[:3], token.Key[3:])

	if err := s.transactionalEmailSender.SendMagicLink(ctx, user, code); err != nil {
		return fmt.Errorf("failed to send magic link: %w", err)
	}

	return nil
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

func (s *UserSerice) SendEmailVerification(ctx context.Context, userID string) error {
	user, err := s.userRepo.Get(userID)
	if err != nil {
		return err
	}

	if user.EmailVerified {
		return nil
	}

	if err := s.transactionalEmailSender.SendConfirmEmail(ctx, user); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *UserSerice) VerifyEmail(ctx context.Context, userID string, rawToken string) (*users.User, error) {
	user, err := s.userRepo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.EmailVerified {
		return user, nil
	}

	emailToken, err := s.jwtService.Verify(ctx, rawToken, jwt.TokenTypeVerifyEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	if emailToken.Subject != userID {
		return nil, fmt.Errorf("invalid token")
	}

	if err := s.setEmailVerified(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to set email verified: %w", err)
	}

	return user, nil
}

func (s *UserSerice) setEmailVerified(ctx context.Context, user *users.User) error {
	if user.EmailVerified {
		return nil
	}

	user.EmailVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserSerice) UsersCount(ctx context.Context) (uint64, error) {
	return s.userRepo.Count(ctx)
}

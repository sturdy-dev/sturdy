package service

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/emails/transactional"
	"getsturdy.com/api/pkg/jwt"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_onetime "getsturdy.com/api/pkg/onetime/service"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/users/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	*service.UserService

	logger                   *zap.Logger
	userRepo                 db_user.Repository
	jwtService               *service_jwt.Service
	transactionalEmailSender transactional.EmailSender
	onetimeService           *service_onetime.Service
	analyticsService         *service_analytics.Service
}

func New(
	userService *service.UserService,
	logger *zap.Logger,
	userRepo db_user.Repository,
	jwtService *service_jwt.Service,
	transactionalEmailSender transactional.EmailSender,
	onetimeService *service_onetime.Service,
	analyticsService *service_analytics.Service,
) *Service {
	return &Service{
		UserService: userService,

		logger:                   logger,
		userRepo:                 userRepo,
		jwtService:               jwtService,
		transactionalEmailSender: transactionalEmailSender,
		onetimeService:           onetimeService,
		analyticsService:         analyticsService,
	}
}

func (s *Service) CreateWithPassword(ctx context.Context, name, password, email string) (*users.User, error) {
	user, err := s.UserService.CreateWithPassword(ctx, name, password, email)
	if err != nil {
		return nil, err
	}

	if err := s.transactionalEmailSender.SendWelcome(ctx, user); err != nil {
		s.logger.Error("failed to send welcome email", zap.Error(err))
	}

	return user, nil
}

func (s *Service) Create(ctx context.Context, name, email string) (*users.User, error) {
	// If user already exists, send OTP
	if existingUser, err := s.userRepo.GetByEmail(email); err == nil {
		if err := s.SendMagicLink(ctx, existingUser); err != nil {
			return nil, fmt.Errorf("failed to send magic link: %w", err)
		}
		return existingUser, nil
	}

	t := time.Now()
	newUser := &users.User{
		ID:        users.ID(uuid.New().String()),
		Name:      name,
		Email:     email,
		CreatedAt: &t,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.analyticsService.IdentifyUser(ctx, newUser)
	s.analyticsService.Capture(ctx, "created account")

	if err := s.transactionalEmailSender.SendWelcome(ctx, newUser); err != nil {
		s.logger.Error("failed to send welcome email", zap.Error(err))
	}

	if err := s.SendMagicLink(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to send magic link: %w", err)
	}

	return newUser, nil
}

func (s *Service) SendMagicLink(ctx context.Context, user *users.User) error {
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

func (s *Service) SendEmailVerification(ctx context.Context, userID users.ID) error {
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

func (s *Service) VerifyMagicLink(ctx context.Context, user *users.User, code string) error {
	if _, err := s.onetimeService.Resolve(ctx, user, code); err != nil {
		return fmt.Errorf("failed to resolve magic link: %w", err)
	}

	if err := s.setEmailVerified(ctx, user); err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}

	s.analyticsService.IdentifyUser(ctx, user)
	s.analyticsService.Capture(ctx, "logged in", analytics.Property("type", "code"))
	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, userID users.ID, rawToken string) (*users.User, error) {
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

	if emailToken.Subject != userID.String() {
		return nil, fmt.Errorf("invalid token")
	}

	if err := s.setEmailVerified(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to set email verified: %w", err)
	}

	return user, nil
}

func (s *Service) setEmailVerified(ctx context.Context, user *users.User) error {
	if user.EmailVerified {
		return nil
	}

	user.EmailVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

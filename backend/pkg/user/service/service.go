package service

import (
	"context"
	"fmt"
	"time"

	"mash/pkg/analytics"
	"mash/pkg/emails/transactional"
	"mash/pkg/jwt"
	service_jwt "mash/pkg/jwt/service"
	service_onetime "mash/pkg/onetime/service"
	"mash/pkg/user"
	db_user "mash/pkg/user/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger                   *zap.Logger
	userRepo                 db_user.Repository
	jwtService               *service_jwt.Service
	onetimeService           *service_onetime.Service
	transactionalEmailSender transactional.EmailSender
	analyticsClient          analytics.Client
}

func New(
	logger *zap.Logger,
	userRepo db_user.Repository,
	jwtService *service_jwt.Service,
	onetimeService *service_onetime.Service,
	transactionalEmailSender transactional.EmailSender,
	analyticsClient analytics.Client,
) *Service {
	return &Service{
		logger:                   logger,
		userRepo:                 userRepo,
		jwtService:               jwtService,
		onetimeService:           onetimeService,
		transactionalEmailSender: transactionalEmailSender,
		analyticsClient:          analyticsClient,
	}
}

func (s *Service) Create(ctx context.Context, name, email string) (*user.User, error) {
	// If user already exists, send OTP
	if existingUser, err := s.userRepo.GetByEmail(email); err == nil {
		if err := s.SendMagicLink(ctx, existingUser); err != nil {
			return nil, fmt.Errorf("failed to send magic link: %w", err)
		}
		return existingUser, nil
	}

	t := time.Now()
	newUser := &user.User{
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

func (s *Service) VerifyMagicLink(ctx context.Context, user *user.User, code string) error {
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

func (s *Service) SendMagicLink(ctx context.Context, user *user.User) error {
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

func (s *Service) GetByIDs(ctx context.Context, ids ...string) ([]*user.User, error) {
	return s.userRepo.GetByIDs(ctx, ids...)
}

func (s *Service) GetByID(_ context.Context, id string) (*user.User, error) {
	return s.userRepo.Get(id)
}

func (s *Service) GetByEmail(_ context.Context, email string) (*user.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *Service) SendEmailVerification(ctx context.Context, userID string) error {
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

func (s *Service) VerifyEmail(ctx context.Context, userID string, rawToken string) (*user.User, error) {
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

func (s *Service) setEmailVerified(ctx context.Context, user *user.User) error {
	if user.EmailVerified {
		return nil
	}

	user.EmailVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *Service) UserCount(ctx context.Context) (int, error) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

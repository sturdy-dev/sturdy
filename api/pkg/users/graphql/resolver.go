package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebases"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/newsletter"
	db_newsletter "getsturdy.com/api/pkg/newsletter/db"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	service_user "getsturdy.com/api/pkg/users/service"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userRootResolver struct {
	userRepo                 db_user.Repository
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository

	userService service_user.Service

	viewRootResolver          resolvers.ViewRootResolver
	notificationRootResolver  resolvers.NotificationRootResolver
	githubAccountRootResolver resolvers.GitHubAccountRootResolver
	analyticsService          *service_analytics.Service
}

func NewResolver(
	userRepo db_user.Repository,
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository,

	userService service_user.Service,

	viewRootResolver resolvers.ViewRootResolver,
	notificationRootResolver resolvers.NotificationRootResolver,
	githubAccountRootResolver resolvers.GitHubAccountRootResolver,

	logger *zap.Logger,
	analyticsService *service_analytics.Service,
) *UserDataloader {
	return NewDataloader(&userRootResolver{
		userRepo:                 userRepo,
		notificationSettingsRepo: notificationSettingsRepo,

		userService: userService,

		viewRootResolver:          viewRootResolver,
		notificationRootResolver:  notificationRootResolver,
		githubAccountRootResolver: githubAccountRootResolver,
		analyticsService:          analyticsService,
	}, logger)
}

func (r *userRootResolver) InternalUser(_ context.Context, id users.ID) (resolvers.UserResolver, error) {
	user, err := r.userRepo.Get(id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &userResolver{root: r, u: user}, nil
}

func (r *userRootResolver) UpdateUser(ctx context.Context, args resolvers.UpdateUserArgs) (resolvers.UserResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	user, err := r.userRepo.Get(userID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if args.Input.Name != nil {
		user.Name = *args.Input.Name
	}
	if args.Input.Email != nil {
		user.Email = *args.Input.Email
	}

	// Update user
	err = r.userRepo.Update(user)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	r.analyticsService.IdentifyUser(ctx, user)

	// User password is updated by a separate method
	if args.Input.Password != nil && len(*args.Input.Password) >= 8 {
		// Salt and hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*args.Input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		user.PasswordHash = string(hashedPassword)

		err = r.userRepo.UpdatePassword(user)
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
	}

	// Update notification settings
	if args.Input.NotificationsReceiveNewsletter != nil {
		// Get existing settings
		settings, err := r.notificationSettingsRepo.GetByUser(user.ID)

		// Create new settings object
		if errors.Is(err, sql.ErrNoRows) {
			err := r.notificationSettingsRepo.Insert(newsletter.NotificationSettings{
				UserID:            user.ID,
				ReceiveNewsletter: *args.Input.NotificationsReceiveNewsletter,
			})
			if err != nil {
				return nil, gqlerrors.Error(err)
			}
		} else if err != nil {
			// could not get settings
			return nil, gqlerrors.Error(err)
		} else {
			// update existing settings
			settings.ReceiveNewsletter = *args.Input.NotificationsReceiveNewsletter
			err := r.notificationSettingsRepo.Update(settings)
			if err != nil {
				return nil, gqlerrors.Error(err)
			}
		}
	}

	return &userResolver{root: r, u: user}, nil
}

func (r *userRootResolver) VerifyEmail(_ context.Context, _ resolvers.VerifyEmailArgs) (resolvers.UserResolver, error) {
	return nil, gqlerrors.Error(gqlerrors.ErrNotImplemented)
}

type userResolver struct {
	root *userRootResolver
	u    *users.User
}

func (r *userResolver) Status() (resolvers.UserStatus, error) {
	switch r.u.Status {
	case users.StatusActive:
		return resolvers.UserStatusActive, nil
	case users.StatusShadow:
		return resolvers.UserStatusShadow, nil
	default:
		return resolvers.UserStatusUndefined, gqlerrors.Error(fmt.Errorf("unknown status: %s", r.u.Status))
	}
}

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(r.u.ID)
}

func (r *userResolver) Name() string {
	return r.u.Name
}

func (r *userResolver) Email() string {
	return r.u.Email
}

func (r *userResolver) EmailVerified() bool {
	return r.u.EmailVerified
}

func (r *userResolver) AvatarUrl() *string {
	return r.u.AvatarURL
}

func (r *userResolver) NotificationPreferences(ctx context.Context) ([]resolvers.NotificationPreferenceResolver, error) {
	return r.root.notificationRootResolver.InternalNotificationPreferences(ctx, r.u.ID)
}

func (r *userResolver) NotificationsReceiveNewsletter() (bool, error) {
	settings, err := r.root.notificationSettingsRepo.GetByUser(r.u.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return true, nil // default to true
	}
	if err != nil {
		return false, gqlerrors.Error(err)
	}
	return settings.ReceiveNewsletter, nil
}

func (r *userResolver) GitHubAccount(ctx context.Context) (resolvers.GitHubAccountResolver, error) {
	if account, err := r.root.githubAccountRootResolver.InteralByID(ctx, r.u.ID); errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	} else {
		return account, nil
	}
}

func (r *userResolver) Views() ([]resolvers.ViewResolver, error) {
	viewsByUser, err := r.root.viewRootResolver.InternalViewsByUser(r.u.ID)
	if errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return viewsByUser, err
}

func (r *userResolver) LastUsedView(ctx context.Context, args resolvers.LastUsedViewArgs) (resolvers.ViewResolver, error) {
	resolver, err := r.root.viewRootResolver.InternalLastUsedViewByUser(ctx, codebases.ID(args.CodebaseID), r.u.ID)
	if errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resolver, err
}

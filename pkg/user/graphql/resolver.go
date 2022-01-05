package graphql

import (
	"context"
	"database/sql"
	"errors"
	"mash/pkg/auth"
	"mash/pkg/github"
	db_github "mash/pkg/github/db"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/newsletter"
	db_newsletter "mash/pkg/newsletter/db"
	"mash/pkg/user"
	db_user "mash/pkg/user/db"
	service_user "mash/pkg/user/service"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userRootResolver struct {
	userRepo                 db_user.Repository
	gitHubUserRepo           db_github.GitHubUserRepo
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository

	userService *service_user.Service

	viewRootResolver         *resolvers.ViewRootResolver
	notificationRootResolver *resolvers.NotificationRootResolver
}

func NewResolver(
	userRepo db_user.Repository,
	gitHubUserRepo db_github.GitHubUserRepo,
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository,

	userService *service_user.Service,

	viewRootResolver *resolvers.ViewRootResolver,
	notificationRootResolver *resolvers.NotificationRootResolver,

	logger *zap.Logger,
) resolvers.UserRootResolver {
	return NewDataloader(&userRootResolver{
		userRepo:                 userRepo,
		gitHubUserRepo:           gitHubUserRepo,
		notificationSettingsRepo: notificationSettingsRepo,

		userService: userService,

		viewRootResolver:         viewRootResolver,
		notificationRootResolver: notificationRootResolver,
	}, logger)
}

func (r *userRootResolver) InternalUser(id string) (resolvers.UserResolver, error) {
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

	// User password is updated by a separate method
	if args.Input.Password != nil && len(*args.Input.Password) >= 8 {
		// Salt and hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*args.Input.Password), 8)
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

func (r *userRootResolver) VerifyEmail(ctx context.Context, args resolvers.VerifyEmailArgs) (resolvers.UserResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	usr, err := r.userService.VerifyEmail(ctx, userID, args.Input.Token)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &userResolver{root: r, u: usr}, nil
}

type userResolver struct {
	root *userRootResolver
	u    *user.User
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
	return (*r.root.notificationRootResolver).InternalNotificationPreferences(ctx, r.u.ID)
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

func (r *userResolver) GitHubAccount() (resolvers.GitHubAccountResolver, error) {
	githubUser, err := r.root.gitHubUserRepo.GetByUserID(r.u.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &gitHubAccountResolver{githubUser: githubUser}, nil
}

func (r *userResolver) Views() ([]resolvers.ViewResolver, error) {
	resolvers, err := (*r.root.viewRootResolver).InternalViewsByUser(r.u.ID)
	if errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resolvers, err
}

func (r *userResolver) LastUsedView(ctx context.Context, args resolvers.LastUsedViewArgs) (resolvers.ViewResolver, error) {
	resolver, err := (*r.root.viewRootResolver).InternalLastUsedViewByUser(ctx, string(args.CodebaseID), r.u.ID)
	if errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resolver, err
}

type gitHubAccountResolver struct {
	githubUser *github.GitHubUser
}

func (r *gitHubAccountResolver) ID() graphql.ID {
	return graphql.ID(r.githubUser.ID)
}

func (r *gitHubAccountResolver) Login() string {
	return r.githubUser.Username
}

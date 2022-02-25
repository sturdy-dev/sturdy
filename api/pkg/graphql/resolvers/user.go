package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/users"
	"github.com/graph-gophers/graphql-go"
)

type UserRootResolver interface {
	// Queries
	User(ctx context.Context) (UserResolver, error)

	// Mutations
	UpdateUser(ctx context.Context, args UpdateUserArgs) (UserResolver, error)
	VerifyEmail(ctx context.Context, args VerifyEmailArgs) (UserResolver, error)

	// Internal
	InternalUser(context.Context, users.ID) (UserResolver, error)
}

type UpdateUserArgs struct {
	Input UpdateUserInput
}

type UpdateUserInput struct {
	Name                           *string
	Email                          *string
	Password                       *string
	NotificationsReceiveNewsletter *bool
}

type VerifyEmailArgs struct {
	Input VerifyEmailInput
}

type VerifyEmailInput struct {
	Token string
}

type UserResolver interface {
	ID() graphql.ID
	Name() string
	Email() string
	EmailVerified() bool
	AvatarUrl() *string
	NotificationPreferences(context.Context) ([]NotificationPreferenceResolver, error)
	GitHubAccount(context.Context) (GitHubAccountResolver, error)
	NotificationsReceiveNewsletter() (bool, error)
	Views() ([]ViewResolver, error)
	LastUsedView(ctx context.Context, args LastUsedViewArgs) (ViewResolver, error)
}

type LastUsedViewArgs struct {
	CodebaseID graphql.ID
}

package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/users"
	"github.com/graph-gophers/graphql-go"
)

type NotificationRootResolver interface {
	Notifications(ctx context.Context) ([]NotificationResolver, error)

	// Mutations
	ArchiveNotifications(ctx context.Context, args ArchiveNotificationsArgs) ([]NotificationResolver, error)
	UpdateNotificationPreference(context.Context, UpdateNotificationPreferenceArgs) (NotificationPreferenceResolver, error)

	// Subscriptions
	UpdatedNotifications(ctx context.Context) (chan NotificationResolver, error)

	// Internal
	InternalNotificationPreferences(context.Context, users.ID) ([]NotificationPreferenceResolver, error)
}

type commonNotificationResolver interface {
	ID() graphql.ID
	Type() (NotificationType, error)
	CreatedAt() int32
	ArchivedAt() *int32
}

type NotificationResolver interface {
	ToCommentNotification() (CommentNotificationResolver, bool)
	ToRequestedReviewNotification() (RequestedReviewNotificationResolver, bool)
	ToReviewNotification() (ReviewNotificationResolver, bool)
	ToNewSuggestionNotification() (NewSuggestionNotificationResolver, bool)
	ToGitHubRepositoryImported() (GitHubRepositoryImportedNotificationResovler, bool)
	ToInvitedToOrganizationNotification() (InvitedToOrganizationNotificationResolver, bool)
	ToInvitedToCodebaseNotification() (InvitedToCodebaseNotificationResolver, bool)

	commonNotificationResolver
}

type CommentNotificationResolver interface {
	commonNotificationResolver
	Comment(ctx context.Context) (CommentResolver, error)
}

type RequestedReviewNotificationResolver interface {
	commonNotificationResolver
	Review(ctx context.Context) (ReviewResolver, error)
}

type ReviewNotificationResolver interface {
	commonNotificationResolver
	Review(ctx context.Context) (ReviewResolver, error)
}

type NewSuggestionNotificationResolver interface {
	commonNotificationResolver
	Suggestion(context.Context) (SuggestionResolver, error)
}

type GitHubRepositoryImportedNotificationResovler interface {
	commonNotificationResolver
	Repository(context.Context) (CodebaseGitHubIntegrationResolver, error)
}

type InvitedToCodebaseNotificationResolver interface {
	commonNotificationResolver
	Codebase(context.Context) (CodebaseResolver, error)
}

type InvitedToOrganizationNotificationResolver interface {
	commonNotificationResolver
	Organization(context.Context) (OrganizationResolver, error)
}

type ArchiveNotificationsArgs struct {
	Input ArchiveNotificationsInput
}

type ArchiveNotificationsInput struct {
	IDs []graphql.ID
}

type UpdateNotificationPreferenceArgs struct {
	Input UpdateNotificationPreferenceInput
}

type UpdateNotificationPreferenceInput struct {
	Type    NotificationType
	Channel NotificationChannel
	Enabled bool
}

type NotificationType string

const (
	NotificationTypeUndefined             NotificationType = ""
	NotificationTypeComment               NotificationType = "Comment"
	NotificationTypeReview                NotificationType = "Review"
	NotificationTypeRequestedReview       NotificationType = "RequestedReview"
	NotificationTypeNewSuggestion         NotificationType = "NewSuggestion"
	NotificationGitHubRepositoryImported  NotificationType = "GitHubRepositoryImported"
	NotificationTypeInvitedToCodebase     NotificationType = "InvitedToCodebase"
	NotificationTypeInvitedToOrganization NotificationType = "InvitedToOrganization"
)

type NotificationChannel string

const (
	NotificationChannelUndefined NotificationChannel = ""
	NotificationChannelWeb       NotificationChannel = "Web"
	NotificationChannelEmail     NotificationChannel = "Email"
)

type NotificationPreferenceResolver interface {
	Type() (NotificationType, error)
	Channel() (NotificationChannel, error)
	Enabled() bool
}

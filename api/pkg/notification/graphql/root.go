package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/notification"
	db_notification "getsturdy.com/api/pkg/notification/db"
	service_notification "getsturdy.com/api/pkg/notification/service"
	db_organizations "getsturdy.com/api/pkg/organization/db"
	"getsturdy.com/api/pkg/suggestions"
	"getsturdy.com/api/pkg/users"

	"github.com/graph-gophers/graphql-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	concurrentUpdatedNotificationsConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "sturdy_graphql_concurrent_subscriptions",
		ConstLabels: prometheus.Labels{"subscription": "updatedNotification"},
	})
)

type notificationRootResolver struct {
	notificationRepository db_notification.Repository
	codebaseUserRepo       db_codebases.CodebaseUserRepository
	organizationUserRepo   db_organizations.MemberRepository
	codebaseRepo           db_codebases.CodebaseRepository

	preferencesService *service_notification.Preferences
	authService        *service_auth.Service

	commentResolver                       resolvers.CommentRootResolver
	codebaseResolver                      resolvers.CodebaseRootResolver
	authorRootResolver                    resolvers.AuthorRootResolver
	workspaceRootResolver                 resolvers.WorkspaceRootResolver
	reviewRootResolver                    resolvers.ReviewRootResolver
	suggestionRootResolver                resolvers.SuggestionRootResolver
	codebaseGitHubIntegrationRootResolver resolvers.CodebaseGitHubIntegrationRootResolver
	organizationResolver                  resolvers.OrganizationRootResolver

	eventsReader events.EventReader
	eventSender  events.EventSender
	logger       *zap.Logger
}

func NewResolver(
	notificationRepository db_notification.Repository,
	codebaseUserRepo db_codebases.CodebaseUserRepository,
	organizationUserRepo db_organizations.MemberRepository,
	codebaseRepo db_codebases.CodebaseRepository,

	preferencesService *service_notification.Preferences,
	authService *service_auth.Service,

	commentResolver resolvers.CommentRootResolver,
	codebaseResolver resolvers.CodebaseRootResolver,
	authorRootResolver resolvers.AuthorRootResolver,
	workspaceRootResolver resolvers.WorkspaceRootResolver,
	reviewRootResolver resolvers.ReviewRootResolver,
	suggestionRootResolver resolvers.SuggestionRootResolver,
	codebaseGitHubIntegrationRootResolver resolvers.CodebaseGitHubIntegrationRootResolver,
	organizationResolver resolvers.OrganizationRootResolver,

	eventsReader events.EventReader,
	eventSender events.EventSender,
	logger *zap.Logger,
) resolvers.NotificationRootResolver {
	return &notificationRootResolver{
		notificationRepository: notificationRepository,
		codebaseUserRepo:       codebaseUserRepo,
		organizationUserRepo:   organizationUserRepo,
		codebaseRepo:           codebaseRepo,

		preferencesService: preferencesService,
		authService:        authService,

		commentResolver:                       commentResolver,
		codebaseResolver:                      codebaseResolver,
		authorRootResolver:                    authorRootResolver,
		workspaceRootResolver:                 workspaceRootResolver,
		reviewRootResolver:                    reviewRootResolver,
		suggestionRootResolver:                suggestionRootResolver,
		codebaseGitHubIntegrationRootResolver: codebaseGitHubIntegrationRootResolver,
		organizationResolver:                  organizationResolver,

		eventsReader: eventsReader,
		eventSender:  eventSender,
		logger:       logger,
	}
}

func (r *notificationRootResolver) Notifications(ctx context.Context) ([]resolvers.NotificationResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		// for anonymous users, we return an empty list
		return nil, nil
	}

	notifications, err := r.notificationRepository.ListByUser(
		userID,
		// TODO: Support pagination
		50,
		0)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make([]resolvers.NotificationResolver, 0, len(notifications))
	for _, notif := range notifications {
		notifResolver := &notificationResolver{notif: notif, root: r}
		// Get the sub-item that this notification is referencing
		// If it can't be resolved, the notification won't be returned
		sub, err := notifResolver.sub(ctx)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, auth.ErrForbidden) || errors.Is(err, ErrUnknownNotificationType) {
			continue
		} else if err != nil {
			r.logger.Error("failed to get sub notification item", zap.Any("notif", notif))
			return nil, gqlerrors.Error(err)
		}
		notifResolver.subItem = sub
		res = append(res, notifResolver)

	}
	return res, nil
}

func convertChannelType(in resolvers.NotificationChannel) (notification.Channel, error) {
	switch in {
	case resolvers.NotificationChannelEmail:
		return notification.ChannelEmail, nil
	case resolvers.NotificationChannelWeb:
		return notification.ChannelWeb, nil
	default:
		return notification.ChannelUndefined, fmt.Errorf("unknown notification channel: %s", in)
	}
}

func convertNotificationType(in resolvers.NotificationType) (notification.NotificationType, error) {
	switch in {
	case resolvers.NotificationTypeComment:
		return notification.CommentNotificationType, nil
	case resolvers.NotificationTypeNewSuggestion:
		return notification.NewSuggestionNotificationType, nil
	case resolvers.NotificationTypeReview:
		return notification.ReviewNotificationType, nil
	case resolvers.NotificationTypeRequestedReview:
		return notification.RequestedReviewNotificationType, nil
	case resolvers.NotificationGitHubRepositoryImported:
		return notification.GitHubRepositoryImported, nil
	case resolvers.NotificationTypeInvitedToCodebase:
		return notification.InvitedToCodebase, nil
	case resolvers.NotificationTypeInvitedToOrganization:
		return notification.InvitedToOrganization, nil
	default:
		return notification.NotificationTypeUndefined, fmt.Errorf("unknown notification type: %s", in)
	}
}

func (r *notificationRootResolver) UpdateNotificationPreference(ctx context.Context, args resolvers.UpdateNotificationPreferenceArgs) (resolvers.NotificationPreferenceResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	typ, err := convertNotificationType(args.Input.Type)
	if err != nil {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "type", err.Error())
	}

	channel, err := convertChannelType(args.Input.Channel)
	if err != nil {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "channel", err.Error())
	}

	pref, err := r.preferencesService.Update(ctx, userID, typ, channel, args.Input.Enabled)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &notificationPreferenceResolver{preference: pref}, nil
}

func (r *notificationRootResolver) internalNotification(ctx context.Context, id string) (resolvers.NotificationResolver, error) {
	notif, err := r.notificationRepository.Get(id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	notifResolver := &notificationResolver{notif: notif, root: r}

	// Get the sub-item that this notification is referencing
	// If it can't be resolved, the notification won't be returned
	sub, err := notifResolver.sub(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	notifResolver.subItem = sub

	return notifResolver, nil
}

func (r *notificationRootResolver) ArchiveNotifications(ctx context.Context, args resolvers.ArchiveNotificationsArgs) ([]resolvers.NotificationResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	var notificationIDs []string
	for _, id := range args.Input.IDs {
		notificationIDs = append(notificationIDs, string(id))
	}
	if len(notificationIDs) == 0 {
		return nil, nil
	}

	if err := r.notificationRepository.ArchiveByUserAndIds(userID, notificationIDs); err != nil {
		return nil, gqlerrors.Error(err)
	}

	notifications, err := r.notificationRepository.ListByUserAndIds(userID, notificationIDs)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.NotificationResolver
	for _, notif := range notifications {
		r.eventSender.User(userID, events.NotificationEvent, notif.ID)
		res = append(res, &notificationResolver{notif: notif, root: r})
	}
	return res, nil
}

func (r *notificationRootResolver) UpdatedNotifications(ctx context.Context) (chan resolvers.NotificationResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	res := make(chan resolvers.NotificationResolver, 100)
	didErrorOut := false

	concurrentUpdatedNotificationsConnections.Inc()

	cancelFunc := r.eventsReader.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		default:
		}

		if eventType == events.NotificationEvent {
			resolver, err := r.internalNotification(ctx, reference)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return events.ErrClientDisconnected
			case res <- resolver:
				if didErrorOut {
					didErrorOut = false
				}
				return nil
			default:
				r.logger.Error("dropped subscription event",
					zap.Stringer("user_id", userID),
					zap.Stringer("event_type", eventType),
					zap.Int("channel_size", len(res)),
				)
				didErrorOut = true
				return nil
			}
		}
		return nil
	})

	go func() {
		<-ctx.Done()
		cancelFunc()
		concurrentUpdatedNotificationsConnections.Dec()
	}()

	return res, nil
}

func (r *notificationRootResolver) InternalNotificationPreferences(ctx context.Context, userID users.ID) ([]resolvers.NotificationPreferenceResolver, error) {
	pp, err := r.preferencesService.ListByUserID(ctx, userID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	rr := make([]resolvers.NotificationPreferenceResolver, 0, len(pp))
	for _, p := range pp {
		rr = append(rr, &notificationPreferenceResolver{preference: p})
	}
	return rr, nil
}

type notificationPreferenceResolver struct {
	preference *notification.Preference
}

func (r *notificationPreferenceResolver) Type() (resolvers.NotificationType, error) {
	return resolveNotificationType(r.preference.Type)
}

func (r *notificationPreferenceResolver) Channel() (resolvers.NotificationChannel, error) {
	switch r.preference.Channel {
	case notification.ChannelEmail:
		return resolvers.NotificationChannelEmail, nil
	case notification.ChannelWeb:
		return resolvers.NotificationChannelWeb, nil
	default:
		return resolvers.NotificationChannelUndefined, fmt.Errorf("unkown notification channel")
	}
}

func (r *notificationPreferenceResolver) Enabled() bool {
	return r.preference.Enabled
}

type notificationResolver struct {
	notif   notification.Notification
	root    *notificationRootResolver
	subItem any
}

func (r *notificationResolver) ID() graphql.ID {
	return graphql.ID(r.notif.ID)
}

func (r *notificationResolver) Type() (resolvers.NotificationType, error) {
	return resolveNotificationType(r.notif.NotificationType)
}

func resolveNotificationType(typ notification.NotificationType) (resolvers.NotificationType, error) {
	switch typ {
	case notification.CommentNotificationType:
		return resolvers.NotificationTypeComment, nil
	case notification.ReviewNotificationType:
		return resolvers.NotificationTypeReview, nil
	case notification.RequestedReviewNotificationType:
		return resolvers.NotificationTypeRequestedReview, nil
	case notification.NewSuggestionNotificationType:
		return resolvers.NotificationTypeNewSuggestion, nil
	case notification.GitHubRepositoryImported:
		return resolvers.NotificationGitHubRepositoryImported, nil
	case notification.InvitedToOrganization:
		return resolvers.NotificationTypeInvitedToOrganization, nil
	case notification.InvitedToCodebase:
		return resolvers.NotificationTypeInvitedToCodebase, nil
	default:
		return resolvers.NotificationTypeUndefined, fmt.Errorf("unknown notification type")
	}
}

func (r *notificationResolver) CreatedAt() int32 {
	return int32(r.notif.CreatedAt.Unix())
}

func (r *notificationResolver) ArchivedAt() *int32 {
	if r.notif.ArchivedAt == nil {
		return nil
	}
	t := int32(r.notif.ArchivedAt.Unix())
	return &t
}

var ErrUnknownNotificationType = fmt.Errorf("unknown notification type")

func (r *notificationResolver) sub(ctx context.Context) (any, error) {
	switch r.notif.NotificationType {
	case notification.CommentNotificationType:
		return r.root.commentResolver.Comment(ctx, resolvers.CommentArgs{ID: graphql.ID(r.notif.ReferenceID)})
	case notification.ReviewNotificationType:
		return r.root.reviewRootResolver.InternalReview(ctx, r.notif.ReferenceID)
	case notification.RequestedReviewNotificationType:
		return r.root.reviewRootResolver.InternalReview(ctx, r.notif.ReferenceID)
	case notification.NewSuggestionNotificationType:
		return r.root.suggestionRootResolver.InternalSuggestionByID(ctx, suggestions.ID(r.notif.ReferenceID))
	case notification.GitHubRepositoryImported:
		return r.root.codebaseGitHubIntegrationRootResolver.InternalGitHubRepositoryByID(r.notif.ReferenceID)
	case notification.InvitedToCodebase:
		member, err := r.root.codebaseUserRepo.GetByID(ctx, r.notif.ReferenceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get codebase member: %w", err)
		}
		id := graphql.ID(member.CodebaseID.String())
		return r.root.codebaseResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
	case notification.InvitedToOrganization:
		member, err := r.root.organizationUserRepo.GetByID(ctx, r.notif.ReferenceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization member: %w", err)
		}
		id := graphql.ID(member.OrganizationID)
		return r.root.organizationResolver.Organization(ctx, resolvers.OrganizationArgs{ID: &id})
	default:
		return resolvers.NotificationTypeUndefined, ErrUnknownNotificationType
	}
}

func (r *notificationResolver) ToInvitedToOrganizationNotification() (resolvers.InvitedToOrganizationNotificationResolver, bool) {
	if r.notif.NotificationType != notification.InvitedToOrganization {
		return nil, false
	}
	return &invitedToOrganizationNotificationResolver{notificationResolver: r}, true
}

func (r *notificationResolver) ToInvitedToCodebaseNotification() (resolvers.InvitedToCodebaseNotificationResolver, bool) {
	if r.notif.NotificationType != notification.InvitedToCodebase {
		return nil, false
	}
	return &invitedToCodebaseNotificationResolver{notificationResolver: r}, true
}

func (r *notificationResolver) ToCommentNotification() (resolvers.CommentNotificationResolver, bool) {
	if r.notif.NotificationType != notification.CommentNotificationType {
		return nil, false
	}

	return &commentNotificationResolver{r}, true
}

func (r *notificationResolver) ToRequestedReviewNotification() (resolvers.RequestedReviewNotificationResolver, bool) {
	if r.notif.NotificationType != notification.RequestedReviewNotificationType {
		return nil, false
	}

	return &requestedReviewNotificationResolver{r}, true
}

func (r *notificationResolver) ToReviewNotification() (resolvers.ReviewNotificationResolver, bool) {
	if r.notif.NotificationType != notification.ReviewNotificationType {
		return nil, false
	}

	return &reviewNotificationResolver{r}, true
}

func (r *notificationResolver) ToGitHubRepositoryImported() (resolvers.GitHubRepositoryImportedNotificationResovler, bool) {
	if r.notif.NotificationType != notification.GitHubRepositoryImported {
		return nil, false
	}
	return &gitHubRepositoryImportedResolver{r}, true
}

func (r *notificationResolver) ToNewSuggestionNotification() (resolvers.NewSuggestionNotificationResolver, bool) {
	if r.notif.NotificationType != notification.NewSuggestionNotificationType {
		return nil, false
	}
	return &newSuggestionNotificationResolver{r}, true
}

type commentNotificationResolver struct {
	*notificationResolver
}

func (r *commentNotificationResolver) Comment(ctx context.Context) (resolvers.CommentResolver, error) {
	if v, ok := r.subItem.(resolvers.CommentResolver); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to get CommentResolver")
}

type requestedReviewNotificationResolver struct {
	*notificationResolver
}

func (r *requestedReviewNotificationResolver) Review(ctx context.Context) (resolvers.ReviewResolver, error) {
	if v, ok := r.subItem.(resolvers.ReviewResolver); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to get ReviewResolver")
}

type reviewNotificationResolver struct {
	*notificationResolver
}

func (r *reviewNotificationResolver) Review(ctx context.Context) (resolvers.ReviewResolver, error) {
	if v, ok := r.subItem.(resolvers.ReviewResolver); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to get ReviewResolver")
}

type newSuggestionNotificationResolver struct {
	*notificationResolver
}

func (r *newSuggestionNotificationResolver) Suggestion(ctx context.Context) (resolvers.SuggestionResolver, error) {
	if v, ok := r.subItem.(resolvers.SuggestionResolver); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to get SuggestionResolver")
}

type gitHubRepositoryImportedResolver struct {
	*notificationResolver
}

func (r *gitHubRepositoryImportedResolver) Repository(ctx context.Context) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	if v, ok := r.subItem.(resolvers.CodebaseGitHubIntegrationResolver); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to get CodebaseGitHubIntegrationResolver")
}

type invitedToCodebaseNotificationResolver struct {
	*notificationResolver
}

func (r *invitedToCodebaseNotificationResolver) Codebase(ctx context.Context) (resolvers.CodebaseResolver, error) {
	if v, ok := r.subItem.(resolvers.CodebaseResolver); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to get CodebaseResolver")
}

type invitedToOrganizationNotificationResolver struct {
	*notificationResolver
}

func (r *invitedToOrganizationNotificationResolver) Organization(ctx context.Context) (resolvers.OrganizationResolver, error) {
	if v, ok := r.subItem.(resolvers.OrganizationResolver); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to get OrganizationResolver")
}

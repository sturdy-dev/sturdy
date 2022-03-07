package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/notification"
	db_notification "getsturdy.com/api/pkg/notification/db"
	service_notification "getsturdy.com/api/pkg/notification/service"
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
	codebaseUserRepo       db_codebase.CodebaseUserRepository
	codebaseRepo           db_codebase.CodebaseRepository

	preferencesService *service_notification.Preferences
	authService        *service_auth.Service

	commentResolver                       resolvers.CommentRootResolver
	codebaseResolver                      resolvers.CodebaseRootResolver
	authorRootResolver                    resolvers.AuthorRootResolver
	workspaceRootResolver                 resolvers.WorkspaceRootResolver
	reviewRootResolver                    resolvers.ReviewRootResolver
	suggestionRootResolver                resolvers.SuggestionRootResolver
	codebaseGitHubIntegrationRootResolver resolvers.CodebaseGitHubIntegrationRootResolver

	eventsReader events.EventReader
	eventSender  events.EventSender
	logger       *zap.Logger
}

func NewResolver(
	notificationRepository db_notification.Repository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	codebaseRepo db_codebase.CodebaseRepository,

	preferencesService *service_notification.Preferences,
	authService *service_auth.Service,

	commentResolver resolvers.CommentRootResolver,
	codebaseResolver resolvers.CodebaseRootResolver,
	authorRootResolver resolvers.AuthorRootResolver,
	workspaceRootResolver resolvers.WorkspaceRootResolver,
	reviewRootResolver resolvers.ReviewRootResolver,
	suggestionRootResolver resolvers.SuggestionRootResolver,
	codebaseGitHubIntegrationRootResolver resolvers.CodebaseGitHubIntegrationRootResolver,

	eventsReader events.EventReader,
	eventSender events.EventSender,
	logger *zap.Logger,
) resolvers.NotificationRootResolver {
	return &notificationRootResolver{
		notificationRepository: notificationRepository,
		codebaseUserRepo:       codebaseUserRepo,
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

	userCodebases, err := r.codebaseUserRepo.GetByUser(userID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	hasAccessToCodebase := map[string]bool{}
	for _, codebase := range userCodebases {
		hasAccessToCodebase[codebase.CodebaseID] = true
	}

	var res []resolvers.NotificationResolver
	for _, notif := range notifications {
		if !hasAccessToCodebase[notif.CodebaseID] {
			continue
		}

		_, codebaseGetErr := r.codebaseRepo.Get(notif.CodebaseID)
		switch {
		case codebaseGetErr == nil:
			notifResolver := &notificationResolver{notif: notif, root: r}

			// Get the sub-item that this notification is referencing
			// If it can't be resolved, the notification won't be returned
			sub, err := notifResolver.sub(ctx)
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else if err != nil {
				r.logger.Error("failed to get sub notification item", zap.Any("notif", notif))
				return nil, gqlerrors.Error(err)
			}

			notifResolver.subItem = sub

			res = append(res, notifResolver)
		case errors.Is(codebaseGetErr, sql.ErrNoRows):
			continue
		default:
			return nil, gqlerrors.Error(err)
		}
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
		close(res)
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
	subItem interface{}
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

func (r *notificationResolver) Codebase(ctx context.Context) (resolvers.CodebaseResolver, error) {
	id := graphql.ID(r.notif.CodebaseID)
	return r.root.codebaseResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
}

func (r *notificationResolver) sub(ctx context.Context) (interface{}, error) {
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
	default:
		return resolvers.NotificationTypeUndefined, fmt.Errorf("unknown notification type")
	}
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

package graphql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"getsturdy.com/api/pkg/activity"
	db_activity "getsturdy.com/api/pkg/activity/db"
	service_activity "getsturdy.com/api/pkg/activity/service"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"go.uber.org/zap"
)

type root struct {
	workspaceActivityRepo      db_activity.ActivityRepository
	workspaceActivityReadsRepo db_activity.ActivityReadsRepository

	authorRootResolver    *resolvers.AuthorRootResolver
	commentRootResolver   *resolvers.CommentRootResolver
	changeRootResolver    *resolvers.ChangeRootResolver
	reviewRootResolver    *resolvers.ReviewRootResolver
	workspaceRootResolver *resolvers.WorkspaceRootResolver

	activityService *service_activity.Service
	authService     *service_auth.Service

	eventsSender events.EventSender
	eventsReader events.EventReader
	logger       *zap.Logger
}

func New(
	workspaceActivityRepo db_activity.ActivityRepository,
	workspaceActivityReadsRepo db_activity.ActivityReadsRepository,

	authorRootResolver *resolvers.AuthorRootResolver,
	commentRootResolver *resolvers.CommentRootResolver,
	changeRootResolver *resolvers.ChangeRootResolver,
	reviewRootResolver *resolvers.ReviewRootResolver,
	workspaceRootResolver *resolvers.WorkspaceRootResolver,

	activityService *service_activity.Service,
	authService *service_auth.Service,

	eventsSender events.EventSender,
	eventsReader events.EventReader,
	logger *zap.Logger,
) resolvers.WorkspaceActivityRootResolver {
	return &root{
		workspaceActivityRepo:      workspaceActivityRepo,
		workspaceActivityReadsRepo: workspaceActivityReadsRepo,

		authorRootResolver:    authorRootResolver,
		commentRootResolver:   commentRootResolver,
		changeRootResolver:    changeRootResolver,
		reviewRootResolver:    reviewRootResolver,
		workspaceRootResolver: workspaceRootResolver,

		activityService: activityService,
		authService:     authService,

		eventsSender: eventsSender,
		eventsReader: eventsReader,
		logger:       logger.Named("workspaceActivityRootResolver"),
	}
}

func (r *root) InternalActivityByWorkspace(ctx context.Context, workspaceID string, args resolvers.WorkspaceActivityArgs) ([]resolvers.WorkspaceActivityResolver, error) {
	unreadOnly := args.Input != nil && args.Input.UnreadOnly != nil && *args.Input.UnreadOnly
	var newerThan time.Time

	if unreadOnly {
		userID, err := auth.UserID(ctx)
		if err != nil {
			// can't filter by unread if not logged in
		} else {
			if read, err := r.workspaceActivityReadsRepo.GetByUserAndWorkspace(ctx, userID, workspaceID); err == nil {
				newerThan = read.LastReadCreatedAt
			} else if !errors.Is(err, sql.ErrNoRows) {
				return nil, gqlerrors.Error(err)
			}
		}
	}

	var activities []*activity.Activity
	var err error

	if newerThan.IsZero() {
		activities, err = r.workspaceActivityRepo.ListByWorkspaceID(ctx, workspaceID)
	} else {
		activities, err = r.workspaceActivityRepo.ListByWorkspaceIDNewerThan(ctx, workspaceID, newerThan)
	}
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.WorkspaceActivityResolver
	for _, activity := range activities {
		res = append(res, &resolver{root: r, activity: activity})
	}
	return res, nil
}

func (r *root) InternalActivity(ctx context.Context, activityID string) (resolvers.WorkspaceActivityResolver, error) {
	activity, err := r.workspaceActivityRepo.Get(ctx, activityID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &resolver{activity: activity, root: r}, nil
}

func (r *root) ReadWorkspaceActivity(ctx context.Context, args resolvers.WorkspaceActivityReadArgs) (resolvers.WorkspaceActivityResolver, error) {
	act, err := r.workspaceActivityRepo.Get(ctx, string(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, act); err != nil {
		return nil, gqlerrors.Error(err)
	}

	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.activityService.MarkAsRead(ctx, userID, act); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &resolver{root: r, activity: act}, nil
}

func (r *root) UpdatedWorkspaceActivity(ctx context.Context) (chan resolvers.WorkspaceActivityResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	c := make(chan resolvers.WorkspaceActivityResolver, 100)
	didErrorOut := false

	cancelFunc := r.eventsReader.SubscribeUser(userID, func(et events.EventType, reference string) error {
		if et != events.WorkspaceUpdatedActivity {
			return nil
		}

		resolver, err := r.InternalActivity(ctx, reference)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return errors.New("disconnected")
		case c <- resolver:
			if didErrorOut {
				didErrorOut = false
			}
			return nil
		default:
			r.logger.Error("dropped subscription event",
				zap.Stringer("user_id", userID),
				zap.Stringer("event_type", et),
				zap.Int("channel_size", len(c)),
			)
			didErrorOut = true
			return nil
		}
	})

	go func() {
		<-ctx.Done()
		cancelFunc()
		close(c)
	}()

	return c, nil
}

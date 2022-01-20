package graphql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"mash/pkg/analytics"
	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	"mash/pkg/change"
	db_change "mash/pkg/change/db"
	db_codebase "mash/pkg/codebase/db"
	"mash/pkg/comments"
	db_comments "mash/pkg/comments/db"
	decorate_comment "mash/pkg/comments/decorate"
	"mash/pkg/comments/live"
	"mash/pkg/comments/vcs"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/notification"
	notification_sender "mash/pkg/notification/sender"
	db_snapshots "mash/pkg/snapshots/db"
	"mash/pkg/user"
	db_user "mash/pkg/user/db"
	db_view "mash/pkg/view/db"
	"mash/pkg/view/events"
	"mash/pkg/workspace"
	"mash/pkg/workspace/activity"
	sender_workspace_activity "mash/pkg/workspace/activity/sender"
	db_workspace "mash/pkg/workspace/db"
	service_workspace_watchers "mash/pkg/workspace/watchers/service"
	"mash/vcs/executor"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	concurrentUpdatedCommentConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "sturdy_graphql_concurrent_subscriptions",
		ConstLabels: prometheus.Labels{"subscription": "updatedComment"},
	})
)

type CommentRootResolver struct {
	executorProvider executor.Provider

	userRepo                 db_user.Repository
	commentsRepo             db_comments.Repository
	snapshotRepo             db_snapshots.Repository
	workspaceReader          db_workspace.WorkspaceReader
	viewRepo                 db_view.Repository
	codebaseUserRepo         db_codebase.CodebaseUserRepository
	changeRepo               db_change.Repository
	workspaceWatchersService *service_workspace_watchers.Service
	authService              *service_auth.Service

	eventsReader       events.EventReader
	eventsSender       events.EventSender
	notificationSender notification_sender.NotificationSender
	activitySender     sender_workspace_activity.ActivitySender

	authorResolver    resolvers.AuthorRootResolver
	workspaceResolver *resolvers.WorkspaceRootResolver
	changeResolver    resolvers.ChangeRootResolver

	logger          *zap.Logger
	analyticsClient analytics.Client
}

func NewResolver(
	userRepo db_user.Repository,
	commentsRepo db_comments.Repository,
	snapshotRepo db_snapshots.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	viewRepo db_view.Repository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	changeRepo db_change.Repository,
	workspaceWatchersService *service_workspace_watchers.Service,
	authService *service_auth.Service,

	eventsSender events.EventSender,
	eventsReader events.EventReader,
	notificationSender notification_sender.NotificationSender,
	activitySender sender_workspace_activity.ActivitySender,

	authorResolver resolvers.AuthorRootResolver,
	workspaceResolver *resolvers.WorkspaceRootResolver,
	changeResolver resolvers.ChangeRootResolver,

	logger *zap.Logger,
	analyticsClient analytics.Client,
	executroProvider executor.Provider,
) resolvers.CommentRootResolver {
	return &CommentRootResolver{
		executorProvider: executroProvider,

		userRepo:                 userRepo,
		commentsRepo:             commentsRepo,
		snapshotRepo:             snapshotRepo,
		workspaceReader:          workspaceReader,
		viewRepo:                 viewRepo,
		codebaseUserRepo:         codebaseUserRepo,
		changeRepo:               changeRepo,
		workspaceWatchersService: workspaceWatchersService,
		authService:              authService,

		eventsSender:       eventsSender,
		eventsReader:       eventsReader,
		notificationSender: notificationSender,
		activitySender:     activitySender,

		authorResolver:    authorResolver,
		workspaceResolver: workspaceResolver,
		changeResolver:    changeResolver,

		logger:          logger,
		analyticsClient: analyticsClient,
	}
}

func (r *CommentRootResolver) Comment(ctx context.Context, args resolvers.CommentArgs) (resolvers.CommentResolver, error) {
	comment, err := r.commentsRepo.Get(comments.ID(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanRead(ctx, comment); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &CommentResolver{comment: comment, root: r}, nil
}

func (r *CommentRootResolver) PreFetchedComment(c comments.Comment) (resolvers.CommentResolver, error) {
	return &CommentResolver{comment: c, root: r}, nil
}

func (r *CommentRootResolver) InternalWorkspaceComments(workspace *workspace.Workspace) ([]resolvers.CommentResolver, error) {
	comms, err := live.GetWorkspaceComments(r.commentsRepo, workspace, r.executorProvider, r.snapshotRepo)
	if err != nil {
		return nil, err
	}

	var res []resolvers.CommentResolver
	for _, c := range comms {
		_ = c
		res = append(res, &CommentResolver{comment: c, root: r})
	}

	return res, nil
}

func (r *CommentRootResolver) UpdatedComment(ctx context.Context, args resolvers.UpdatedCommentArgs) (<-chan resolvers.CommentResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	ws, err := r.workspaceReader.Get(string(args.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// TODO: derive from ws.ViewID instead
	var viewID *string
	if args.ViewID != nil {
		s := string(*args.ViewID)
		viewID = &s
	}

	res := make(chan resolvers.CommentResolver, 100)
	didErrorOut := false

	concurrentUpdatedCommentConnections.Inc()

	cancelFunc := r.eventsReader.SubscribeUser(userID, func(et events.EventType, reference string) error {
		// Get all comments if there is a new comment, or if the diffs have changed
		// This is a rather expensive operation, so ideally it should only be done for the comments that are updated, and not all of them
		workspaceCommentUpdated := et == events.WorkspaceUpdatedComments && reference != ws.ID
		viewUpdated := et == events.ViewUpdated && viewID != nil && reference != *viewID
		if !workspaceCommentUpdated && !viewUpdated {
			return nil
		}

		reloadedWs, err := r.workspaceReader.Get(string(args.WorkspaceID))
		if err != nil {
			return gqlerrors.Error(err)
		}

		allResolvers, err := r.InternalWorkspaceComments(reloadedWs)
		if err != nil {
			return err
		}

		for _, resolver := range allResolvers {
			select {
			case <-ctx.Done():
				return errors.New("disconnected")
			case res <- resolver:
				if didErrorOut {
					didErrorOut = false
				}
				return nil
			default:
				r.logger.Error("dropped subscription event",
					zap.String("user_id", userID),
					zap.String("codebase_id", ws.CodebaseID),
					zap.Stringer("event_type", et),
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
		concurrentUpdatedCommentConnections.Dec()
	}()

	return res, nil
}

func (r *CommentRootResolver) UpdateComment(ctx context.Context, args resolvers.UpdateCommentArgs) (resolvers.CommentResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	comm, err := r.commentsRepo.Get(comments.ID(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	comment := &comm

	if r.authService.CanWrite(ctx, comment) != nil {
		return nil, gqlerrors.Error(err)
	}

	codebaseUsers, err := r.getUsersByCodebaseID(ctx, comment.CodebaseID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	comment.Message = args.Input.Message

	mentions := decorate_comment.ExtractNameMentions(comment.Message, codebaseUsers)
	// replace all mentions with ids
	for mention, user := range mentions {
		comment.Message = strings.ReplaceAll(comment.Message, mention, fmt.Sprintf("@%s", user.ID))
	}

	if err := r.commentsRepo.Update(*comment); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if comment.WorkspaceID != nil {
		if err := r.eventsSender.Codebase(comment.CodebaseID, events.WorkspaceUpdatedComments, *comment.WorkspaceID); err != nil {
			r.logger.Error("failed to send workspace updated comments event", zap.Error(err))
			// do not fail
		}
	}

	if err := r.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: userID,
		Event:      "updated comment",
		Properties: analytics.NewProperties().
			Set("comment_id", comment.ID).
			Set("codebase_id", comment.CodebaseID),
	}); err != nil {
		r.logger.Error("analytics failed", zap.Error(err))
	}

	return &CommentResolver{root: r, comment: *comment}, nil
}

func (r *CommentRootResolver) DeleteComment(ctx context.Context, args resolvers.DeleteCommentArgs) (resolvers.CommentResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	comm, err := r.commentsRepo.Get(comments.ID(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, comm); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: userID,
		Event:      "deleted comment",
		Properties: analytics.NewProperties().
			Set("comment_id", comm.ID).
			Set("codebase_id", comm.CodebaseID),
	}); err != nil {
		r.logger.Error("analytics failed", zap.Error(err))
	}

	t := time.Now()
	comm.DeletedAt = &t
	if err := r.commentsRepo.Update(comm); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &CommentResolver{root: r, comment: comm}, nil
}

func (r *CommentRootResolver) getUsersByCodebaseID(ctx context.Context, codebaseID string) ([]*user.User, error) {
	codebaseUsers, err := r.codebaseUserRepo.GetByCodebase(codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase users: %w", err)
	}
	userIDs := make([]string, 0, len(codebaseUsers))
	for _, codebaseUser := range codebaseUsers {
		userIDs = append(userIDs, codebaseUser.UserID)
	}

	users, err := r.userRepo.GetByIDs(ctx, userIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (r *CommentRootResolver) CreateComment(ctx context.Context, args resolvers.CreateCommentArgs) (resolvers.CommentResolver, error) {
	var comment *comments.Comment
	var err error
	if args.Input.InReplyTo == nil {
		comment, err = r.prepareTopComment(ctx, args)
	} else {
		comment, err = r.prepareReplyComment(ctx, args)
	}
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	codebaseUsers, err := r.getUsersByCodebaseID(ctx, comment.CodebaseID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	mentions := decorate_comment.ExtractNameMentions(comment.Message, codebaseUsers)
	// replace all mentions with ids
	for mention, user := range mentions {
		comment.Message = strings.ReplaceAll(comment.Message, mention, fmt.Sprintf("@%s", user.ID))
	}

	if err := r.commentsRepo.Create(*comment); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if comment.WorkspaceID == nil {
		return &CommentResolver{root: r, comment: *comment}, nil
	}

	// all mentioned users start watching the workspace
	for _, mentionedUser := range mentions {
		if _, err := r.workspaceWatchersService.Watch(ctx, mentionedUser.ID, *comment.WorkspaceID); err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to watch workspace: %w", err))
		}
	}

	// comment author starts watching the workspace
	if _, err := r.workspaceWatchersService.Watch(ctx, comment.UserID, *comment.WorkspaceID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to watch workspace: %w", err))
	}

	// Send events
	if err := r.eventsSender.Codebase(comment.CodebaseID, events.WorkspaceUpdatedComments, *comment.WorkspaceID); err != nil {
		r.logger.Error("failed to send workspace updated comments event", zap.Error(err))
		// do not fail
	}

	watchers, err := r.workspaceWatchersService.ListWatchers(ctx, *comment.WorkspaceID)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to list workspace watchers: %w", err))
	}
	for _, watcher := range watchers {
		// Skip sending notification to the user who created the comment
		if watcher.UserID == comment.UserID {
			continue
		}
		if err := r.notificationSender.User(ctx, watcher.UserID, comment.CodebaseID, notification.CommentNotificationType, string(comment.ID)); err != nil {
			r.logger.Error("failed to send comment notification", zap.Error(err))
			// do not fail
		}
	}

	// Create workspace activity
	if err := r.activitySender.Codebase(ctx, comment.CodebaseID, *comment.WorkspaceID, comment.UserID, activity.WorkspaceActivityTypeComment, string(comment.ID)); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &CommentResolver{root: r, comment: *comment}, nil
}

func (r *CommentRootResolver) prepareTopComment(ctx context.Context, args resolvers.CreateCommentArgs) (*comments.Comment, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}
	// Creating a top level comment
	if args.Input.WorkspaceID == nil && args.Input.ChangeID == nil {
		return nil, fmt.Errorf("either workspaceID or changeID must be set for top level comments")
	}
	if args.Input.WorkspaceID != nil && args.Input.ChangeID != nil {
		return nil, fmt.Errorf("workspaceID and changeID can not be set at the same time")
	}

	// Either all of Path, LineIsNew, LineStart, and LineEnd are set. Or none of them are.
	if !allAreEqual(
		args.Input.Path == nil,
		args.Input.LineIsNew == nil,
		args.Input.LineStart == nil,
		args.Input.LineEnd == nil,
	) {
		return nil, fmt.Errorf("path, lineIsNew, lineStart or lineEnd is not set")
	}

	var codebaseID string
	var workspaceID *string
	var changeID *change.ID

	var optionalContext *string
	var optionalContextStartsAt *int

	// Comment in a workspace
	if args.Input.WorkspaceID != nil {
		wid := string(*args.Input.WorkspaceID)
		workspaceID = &wid
		// get and auth against workspace
		ws, err := r.workspaceReader.Get(wid)
		if err != nil {
			return nil, err
		}

		if err := r.authService.CanWrite(ctx, ws); err != nil {
			return nil, err
		}

		codebaseID = ws.CodebaseID

		// Comment on code
		if args.Input.Path != nil {
			// Build context
			context, contextStartsAt, err := vcs.GetContext(int(*args.Input.LineStart), *args.Input.LineIsNew, *args.Input.Path, ws, r.executorProvider, r.snapshotRepo)
			if err != nil {
				return nil, fmt.Errorf("failed to create context: %w", err)
			}
			optionalContext = &context
			optionalContextStartsAt = &contextStartsAt
		}
	} else {
		// Comment on a change
		cid := change.ID(*args.Input.ChangeID)
		ch, err := r.changeRepo.Get(cid)
		if err != nil {
			return nil, err
		}

		if err := r.authService.CanWrite(ctx, ch); err != nil {
			return nil, err
		}

		changeID = &cid
		codebaseID = ch.CodebaseID
	}

	id := comments.ID(uuid.NewString())

	if err := r.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: userID,
		Event:      "created comment",
		Properties: analytics.NewProperties().
			Set("is_reply", false).
			Set("comment_id", id).
			Set("codebase_id", codebaseID).
			Set("workspace_id", workspaceID).
			Set("change_id", changeID),
	}); err != nil {
		r.logger.Error("analytics failed", zap.Error(err))
	}

	newComm := &comments.Comment{
		ID:                  id,
		UserID:              userID,
		CreatedAt:           time.Now(),
		Message:             args.Input.Message,
		CodebaseID:          codebaseID,
		WorkspaceID:         workspaceID,
		ChangeID:            changeID,
		Context:             optionalContext,
		ContextStartsAtLine: optionalContextStartsAt,
	}

	// Comment on code
	if args.Input.Path != nil {
		newComm.Path = *args.Input.Path
		newComm.LineStart = int(*args.Input.LineStart)
		newComm.LineEnd = int(*args.Input.LineEnd)
		newComm.LineIsNew = *args.Input.LineIsNew
	}

	return newComm, nil
}

func (r *CommentRootResolver) prepareReplyComment(ctx context.Context, args resolvers.CreateCommentArgs) (*comments.Comment, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}
	parentID := comments.ID(*args.Input.InReplyTo)

	// Get more meta from parent comment
	parent, err := r.commentsRepo.Get(parentID)
	if err != nil {
		return nil, err
	}

	// We can only reply to top comments
	if parent.ParentComment != nil {
		return nil, errors.New("can not reply to another reply")
	}

	id := comments.ID(uuid.NewString())
	comment := &comments.Comment{
		ID:            id,
		UserID:        userID,
		CreatedAt:     time.Now(),
		Message:       args.Input.Message,
		ParentComment: &parentID,
		CodebaseID:    parent.CodebaseID,  // Not exposed on the API for reply comments
		WorkspaceID:   parent.WorkspaceID, // Not exposed on the API for reply comments, but is used to generate/route events
	}

	if err := r.authService.CanWrite(ctx, comment); err != nil {
		return nil, err
	}

	if err := r.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: userID,
		Event:      "created comment",
		Properties: analytics.NewProperties().
			Set("is_reply", true).
			Set("comment_id", id).
			Set("codebase_id", parent.CodebaseID).
			Set("workspace_id", parent.WorkspaceID).
			Set("change_id", parent.ChangeID),
	}); err != nil {
		r.logger.Error("analytics failed", zap.Error(err))
	}

	return comment, nil
}

type CommentResolver struct {
	root    *CommentRootResolver
	comment comments.Comment
}

func (r *CommentResolver) ToReplyComment() (resolvers.ReplyCommentResolver, bool) {
	if r.comment.ParentComment != nil {
		return &ReplyCommentResolver{CommentResolver: r}, true
	}
	return nil, false
}

func (r *CommentResolver) ToTopComment() (resolvers.TopCommentResolver, bool) {
	if r.comment.ParentComment == nil {
		return &TopCommentResolver{CommentResolver: r}, true
	}
	return nil, false
}

func (r *CommentResolver) ID() graphql.ID {
	return graphql.ID(r.comment.ID)
}

func (r *CommentResolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	return r.root.authorResolver.Author(ctx, graphql.ID(r.comment.UserID))
}

func (r *CommentResolver) CreatedAt() int32 {
	return int32(r.comment.CreatedAt.Unix())
}

func (r *CommentResolver) DeletedAt() *int32 {
	if r.comment.DeletedAt == nil {
		return nil
	}
	t := int32(r.comment.CreatedAt.Unix())
	return &t
}

func (r *CommentResolver) Message() string {
	return r.comment.Message
}

func allAreEqual(a ...bool) bool {
	if len(a) == 0 {
		return true
	}
	for _, v := range a {
		if v != a[0] {
			return false
		}
	}
	return true
}

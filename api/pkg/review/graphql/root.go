package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/notification"
	"getsturdy.com/api/pkg/notification/sender"
	"getsturdy.com/api/pkg/review"
	db_review "getsturdy.com/api/pkg/review/db"
	"getsturdy.com/api/pkg/workspaces/activity"
	activity_sender "getsturdy.com/api/pkg/workspaces/activity/sender"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace_watchers "getsturdy.com/api/pkg/workspaces/watchers/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type reviewRootResolver struct {
	logger *zap.Logger

	reviewRepo      db_review.ReviewRepository
	workspaceReader db_workspaces.WorkspaceReader
	authService     *service_auth.Service

	authorRootResolver    resolvers.AuthorRootResolver
	workspaceRootResolver *resolvers.WorkspaceRootResolver

	eventsSender       events.EventSender
	eventsReader       events.EventReader
	notificationSender sender.NotificationSender
	activitySender     activity_sender.ActivitySender

	analyticsService *service_analytics.Service

	workspaceWatchersService *service_workspace_watchers.Service
}

func New(
	logger *zap.Logger,
	reviewRepo db_review.ReviewRepository,
	workspaceReader db_workspaces.WorkspaceReader,
	authService *service_auth.Service,

	authorRootResolver resolvers.AuthorRootResolver,
	workspaceRootResolver *resolvers.WorkspaceRootResolver,

	eventsSender events.EventSender,
	eventsReader events.EventReader,
	notificationSender sender.NotificationSender,
	activitySender activity_sender.ActivitySender,

	analyticsService *service_analytics.Service,

	workspaceWatchersService *service_workspace_watchers.Service,
) resolvers.ReviewRootResolver {
	return &reviewRootResolver{
		logger: logger.Named("reviewRootResolver"),

		reviewRepo:      reviewRepo,
		workspaceReader: workspaceReader,
		authService:     authService,

		authorRootResolver:    authorRootResolver,
		workspaceRootResolver: workspaceRootResolver,

		eventsSender:       eventsSender,
		eventsReader:       eventsReader,
		notificationSender: notificationSender,
		activitySender:     activitySender,

		analyticsService: analyticsService,

		workspaceWatchersService: workspaceWatchersService,
	}
}

func (r *reviewRootResolver) InternalReviews(ctx context.Context, workspaceID string) ([]resolvers.ReviewResolver, error) {
	reviews, err := r.reviewRepo.ListLatestByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.ReviewResolver
	for _, rev := range reviews {
		res = append(res, &reviewResolver{root: r, rev: rev})
	}
	return res, nil
}

func (r *reviewRootResolver) InternalReview(ctx context.Context, id string) (resolvers.ReviewResolver, error) {
	rev, err := r.reviewRepo.Get(ctx, id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &reviewResolver{root: r, rev: rev}, nil
}

func (r *reviewRootResolver) CreateOrUpdateReview(ctx context.Context, args resolvers.CreateReviewArgs) (resolvers.ReviewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	workspaceID := string(args.Input.WorkspaceID)

	ws, err := r.workspaceReader.Get(workspaceID)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get workspace: %w", err))
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if ws.UserID == userID {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "cannot review your own workspace")
	}

	// reviewer starts watching the workspace
	if _, err := r.workspaceWatchersService.Watch(ctx, userID, ws.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to watch workspace: %w", err))
	}

	var inputGrade review.ReviewGrade
	switch args.Input.Grade {
	case "Approve":
		inputGrade = review.ReviewGradeApprove
	case "Reject":
		inputGrade = review.ReviewGradeReject
	default:
		return nil, gqlerrors.Error(fmt.Errorf("unexpected grade: '%s'", args.Input.Grade))
	}

	// Mark existing as replaced
	if existing, err := r.reviewRepo.GetLatestByUserAndWorkspace(ctx, userID, workspaceID); err == nil {
		// If this review is the same as the existing one, and the review is not dismissed, don't change anything
		if existing.DismissedAt == nil && existing.Grade == inputGrade {
			return &reviewResolver{root: r, rev: existing}, nil
		}

		existing.IsReplaced = true
		if err := r.reviewRepo.Update(ctx, existing); err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to update review: %w", err))
		}

		// Keep going, create a new review
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get existing review: %w", err))
	}

	// Create new
	rev := review.Review{
		ID:          uuid.NewString(),
		UserID:      userID,
		CodebaseID:  ws.CodebaseID,
		WorkspaceID: workspaceID,
		Grade:       inputGrade,
		CreatedAt:   time.Now(),
	}

	if err := r.reviewRepo.Create(ctx, rev); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to create review: %w", err))
	}

	if err := r.activitySender.Codebase(ctx, ws.CodebaseID, ws.ID, userID, activity.WorkspaceActivityTypeReviewed, rev.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to create activity: %w", err))
	}

	// Send notification to the workspace owner
	if err := r.notificationSender.User(ctx, ws.UserID, ws.CodebaseID, notification.ReviewNotificationType, rev.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to send notification: %w", err))
	}

	// Send events
	if err := r.eventsSender.Codebase(ws.CodebaseID, events.WorkspaceUpdatedReviews, ws.ID); err != nil {
		r.logger.Error("failed to send codebase event", zap.Error(err))
		// do not fail
	}

	if err := r.eventsSender.Workspace(ws.ID, events.ReviewUpdated, rev.ID); err != nil {
		r.logger.Error("failed to send workspace event", zap.Error(err))
		// do not fail
	}

	r.analyticsService.Capture(ctx, "review created",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
		analytics.Property("grade", rev.Grade),
	)

	return &reviewResolver{root: r, rev: &rev}, nil
}

func (r *reviewRootResolver) RequestReview(ctx context.Context, args resolvers.RequestReviewArgs) (resolvers.ReviewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	workspaceID := string(args.Input.WorkspaceID)

	ws, err := r.workspaceReader.Get(workspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// requester starts watching the workspace
	if _, err := r.workspaceWatchersService.Watch(ctx, userID, ws.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to watch workspace: %w", err))
	}

	// user requested review from starts watching the workspace
	if _, err := r.workspaceWatchersService.Watch(ctx, string(args.Input.UserID), ws.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to watch workspace: %w", err))
	}

	if existing, err := r.reviewRepo.GetLatestByUserAndWorkspace(ctx, string(args.Input.UserID), workspaceID); err == nil {
		// Don't request a review if this user already has a approved or rejected review
		if existing.DismissedAt == nil && !existing.IsReplaced {
			return &reviewResolver{root: r, rev: existing}, nil
		}

		// Mark as replaced, and create a new review
		existing.IsReplaced = true
		if err := r.reviewRepo.Update(ctx, existing); err != nil {
			return nil, gqlerrors.Error(err)
		}

		// Keep going
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, gqlerrors.Error(err)
	}

	// Create new
	rev := review.Review{
		ID:          uuid.NewString(),
		UserID:      string(args.Input.UserID),
		CodebaseID:  ws.CodebaseID,
		WorkspaceID: workspaceID,
		Grade:       review.ReviewGradeRequested,
		CreatedAt:   time.Now(),
		RequestedBy: &userID,
	}

	if err := r.reviewRepo.Create(ctx, rev); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.activitySender.Codebase(ctx, ws.CodebaseID, ws.ID, userID, activity.WorkspaceActivityTypeRequestedReview, rev.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to create activity: %w", err))
	}

	// Send notification to the user that the review was requested from
	if err := r.notificationSender.User(ctx, string(args.Input.UserID), ws.CodebaseID, notification.RequestedReviewNotificationType, rev.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to send notification: %w", err))
	}

	// Send events
	if err := r.eventsSender.Codebase(ws.CodebaseID, events.WorkspaceUpdatedReviews, ws.ID); err != nil {
		r.logger.Error("failed to send codebase event", zap.Error(err))
		// do not fail
	}

	if err := r.eventsSender.Workspace(ws.ID, events.ReviewUpdated, rev.ID); err != nil {
		r.logger.Error("failed to send workspace event", zap.Error(err))
		// do not fail
	}
	r.analyticsService.Capture(ctx, "review requested",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
		analytics.Property("user_id", rev.UserID),
	)

	return &reviewResolver{root: r, rev: &rev}, nil
}

func (r *reviewRootResolver) DismissReview(ctx context.Context, args resolvers.DismissReviewArgs) (resolvers.ReviewResolver, error) {
	rev, err := r.reviewRepo.Get(ctx, string(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, rev); err != nil {
		return nil, gqlerrors.Error(err)
	}

	ts := time.Now()
	rev.DismissedAt = &ts
	if err := r.reviewRepo.Update(ctx, rev); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Send events
	if err := r.eventsSender.Codebase(rev.CodebaseID, events.WorkspaceUpdatedReviews, rev.ID); err != nil {
		r.logger.Error("failed to send codebase event", zap.Error(err))
		// do not fail
	}

	if err := r.eventsSender.Workspace(rev.WorkspaceID, events.ReviewUpdated, rev.ID); err != nil {
		r.logger.Error("failed to send workspace event", zap.Error(err))
		// do not fail
	}

	r.analyticsService.Capture(ctx, "review dismissed",
		analytics.CodebaseID(rev.CodebaseID),
		analytics.Property("workspace_id", rev.WorkspaceID),
		analytics.Property("user_id", rev.UserID),
	)

	return &reviewResolver{root: r, rev: rev}, nil
}

func (r *reviewRootResolver) InternalDismissAllInWorkspace(ctx context.Context, workspaceID string) error {
	if err := r.reviewRepo.DismissAllInWorkspace(ctx, workspaceID); err != nil {
		return err
	}
	return nil
}

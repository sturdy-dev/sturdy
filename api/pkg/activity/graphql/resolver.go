package graphql

import (
	"context"
	"database/sql"
	"errors"

	"getsturdy.com/api/pkg/activity"
	"getsturdy.com/api/pkg/auth"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type resolver struct {
	root     *root
	activity *activity.Activity
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.activity.ID)
}

func (r *resolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	return (*r.root.authorRootResolver).Author(ctx, graphql.ID(r.activity.UserID))
}

func (r *resolver) CreatedAt() int32 {
	return int32(r.activity.CreatedAt.Unix())
}

func (r *resolver) IsRead(ctx context.Context) (bool, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		// for anonymous users, we consider them to be read
		return true, nil
	}

	read, err := r.root.workspaceActivityReadsRepo.GetByUserAndWorkspace(ctx, userID, r.activity.WorkspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, gqlerrors.Error(err)
	}

	if read.LastReadCreatedAt.Before(r.activity.CreatedAt) {
		return false, nil
	}
	return true, nil
}

func (r *resolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	t := true
	res, err := (*r.root.workspaceRootResolver).Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(r.activity.WorkspaceID), AllowArchived: &t})
	if errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return res, nil
}

func (r *resolver) ToWorkspaceCommentActivity() (resolvers.WorkspaceCommentActivityResolver, bool) {
	if r.activity.ActivityType != activity.TypeComment {
		return nil, false
	}
	return r, true
}

func (r *resolver) ToWorkspaceCreatedChangeActivity() (resolvers.WorkspaceCreatedChangeActivityResolver, bool) {
	if r.activity.ActivityType != activity.TypeCreatedChange {
		return nil, false
	}
	return r, true
}

func (r *resolver) ToWorkspaceRequestedReviewActivity() (resolvers.WorkspaceRequestedReviewActivityResolver, bool) {
	if r.activity.ActivityType != activity.TypeRequestedReview {
		return nil, false
	}
	return r, true
}

func (r *resolver) ToWorkspaceReviewedActivity() (resolvers.WorkspaceReviewedActivityResolver, bool) {
	if r.activity.ActivityType != activity.TypeReviewed {
		return nil, false
	}
	return r, true
}

func (r *resolver) Comment(ctx context.Context) (resolvers.CommentResolver, error) {
	return (*r.root.commentRootResolver).Comment(ctx, resolvers.CommentArgs{ID: graphql.ID(r.activity.Reference)})
}

func (r *resolver) Change(ctx context.Context) (resolvers.ChangeResolver, error) {
	id := graphql.ID(r.activity.Reference)
	return (*r.root.changeRootResolver).Change(ctx, resolvers.ChangeArgs{ID: &id})
}

func (r *resolver) Review(ctx context.Context) (resolvers.ReviewResolver, error) {
	return (*r.root.reviewRootResolver).InternalReview(ctx, r.activity.Reference)
}

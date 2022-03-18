package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	db_comments "getsturdy.com/api/pkg/comments/db"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/vcs/executor"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type ChangeRootResolver struct {
	svc *service.Service

	commentsRepo db_comments.Repository

	authService *service_auth.Service

	commentResolver   *resolvers.CommentRootResolver
	authorResolver    resolvers.AuthorRootResolver
	statusResovler    *resolvers.StatusesRootResolver
	downloadsResovler resolvers.ContentsDownloadUrlRootResolver
	workspaceResolver *resolvers.WorkspaceRootResolver
	codebaseResolver  *resolvers.CodebaseRootResolver
	activityResovler  resolvers.ActivityRootResolver

	executorProvider executor.Provider

	logger *zap.Logger
}

func NewResolver(
	svc *service.Service,

	commentsRepo db_comments.Repository,

	authService *service_auth.Service,

	commentResolver *resolvers.CommentRootResolver,
	authorResolver resolvers.AuthorRootResolver,
	statusResovler *resolvers.StatusesRootResolver,
	downloadsResovler resolvers.ContentsDownloadUrlRootResolver,
	workspaceResolver *resolvers.WorkspaceRootResolver,
	codebaseResolver *resolvers.CodebaseRootResolver,
	activityResovler resolvers.ActivityRootResolver,

	executorProvider executor.Provider,

	logger *zap.Logger,
) resolvers.ChangeRootResolver {
	return &ChangeRootResolver{
		svc: svc,

		commentsRepo: commentsRepo,

		authService: authService,

		commentResolver:   commentResolver,
		authorResolver:    authorResolver,
		statusResovler:    statusResovler,
		downloadsResovler: downloadsResovler,
		workspaceResolver: workspaceResolver,
		codebaseResolver:  codebaseResolver,
		activityResovler:  activityResovler,

		executorProvider: executorProvider,

		logger: logger,
	}
}

func (r *ChangeRootResolver) IntenralListChanges(ctx context.Context, codebaseID codebases.ID, limit int, before *graphql.ID) ([]resolvers.ChangeResolver, error) {
	var beforeChange *changes.ID
	if before != nil {
		changeID := changes.ID(*before)
		beforeChange = &changeID
	}

	changes, err := r.svc.Changelog(ctx, codebaseID, limit, beforeChange)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.ChangeResolver
	for _, change := range changes {
		res = append(res, &ChangeResolver{root: r, ch: change})
	}

	return res, nil
}

func (r *ChangeRootResolver) Change(ctx context.Context, args resolvers.ChangeArgs) (resolvers.ChangeResolver, error) {
	var ch *changes.Change
	var err error

	if args.ID != nil {
		// Lookup by ChangeID
		ch, err = r.svc.GetChangeByID(ctx, changes.ID(*args.ID))
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to lookup by id: %w", err))
		}
	} else if args.CommitID != nil && args.CodebaseID != nil {
		// Lookup by CommitID and CodebaseID
		ch, err = r.svc.GetByCommitAndCodebase(ctx, string(*args.CommitID), codebases.ID(*args.CodebaseID))
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to lookup by commit id and codebaseid: %w", err))
		}
	} else {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "neither id nor commitID is set")
	}

	if err := r.authService.CanRead(ctx, ch); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &ChangeResolver{root: r, ch: ch}, nil
}

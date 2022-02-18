package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/change"
	"getsturdy.com/api/pkg/change/service"
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

		executorProvider: executorProvider,

		logger: logger,
	}
}

func (r *ChangeRootResolver) Change(ctx context.Context, args resolvers.ChangeArgs) (resolvers.ChangeResolver, error) {
	var ch *change.Change
	var err error

	if args.ID != nil {
		// Lookup by ChangeID
		ch, err = r.svc.GetChangeByID(ctx, change.ID(*args.ID))
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to lookup by id: %w", err))
		}
	} else if args.CommitID != nil && args.CodebaseID != nil {
		// Lookup by CommitID and CodebaseID
		ch, err = r.svc.GetByCommitID(ctx, string(*args.CommitID), string(*args.CodebaseID))
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

type ChangeResolver struct {
	ch   *change.Change
	root *ChangeRootResolver
}

func (r *ChangeResolver) ID() graphql.ID {
	return graphql.ID(r.ch.ID)
}

func (r *ChangeResolver) Comments() ([]resolvers.TopCommentResolver, error) {
	comms, err := r.root.commentsRepo.GetByCodebaseAndChange(r.ch.CodebaseID, r.ch.ID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.TopCommentResolver
	for _, comm := range comms {
		resolver, err := (*r.root.commentResolver).PreFetchedComment(comm)
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		if topCommentResolver, ok := resolver.ToTopComment(); ok {
			res = append(res, topCommentResolver)
		}
	}

	return res, nil
}

func (r *ChangeResolver) Title() string {
	if r.ch.Title == nil {
		return "Untitled" // TODO: Is this a bug?
	}
	return *r.ch.Title
}

func (r *ChangeResolver) Description() string {
	return r.ch.UpdatedDescription
}

func (r *ChangeResolver) TrunkCommitID() (*string, error) {
	return r.ch.CommitID, nil
}

func (r *ChangeResolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	// TODO: fetch this data from Git
	if r.ch.UserID == nil {
		if r.ch.GitCreatorName != nil && r.ch.GitCreatorEmail != nil {
			return r.root.authorResolver.InternalAuthorFromNameAndEmail(ctx, *r.ch.GitCreatorName, *r.ch.GitCreatorEmail), nil
		} else {
			return r.root.authorResolver.InternalAuthorFromNameAndEmail(ctx, "Unknown", "unknown@getsturdy.com"), nil
		}
	}
	author, err := r.root.authorResolver.Author(ctx, graphql.ID(*r.ch.UserID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return author, nil
}

func (r *ChangeResolver) CreatedAt() int32 {
	if r.ch.CreatedAt != nil {
		return int32(r.ch.CreatedAt.Unix())
	}
	if r.ch.GitCreatedAt != nil {
		return int32(r.ch.GitCreatedAt.Unix())
	}
	return 0
}

func (r *ChangeResolver) Diffs(ctx context.Context) ([]resolvers.FileDiffResolver, error) {
	if r.ch.CommitID == nil {
		return nil, gqlerrors.ErrNotFound
	}

	allower, err := r.root.authService.GetAllower(ctx, r.ch)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	diffs, err := r.root.svc.Diffs(ctx, r.ch, allower)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make([]resolvers.FileDiffResolver, len(diffs))
	for k, v := range diffs {
		res[k] = &fileDiffResolver{diff: v}
	}
	return res, nil
}

func (r *ChangeResolver) Statuses(ctx context.Context) ([]resolvers.StatusResolver, error) {
	if r.ch.CommitID == nil {
		return nil, gqlerrors.ErrNotFound
	}
	return (*r.root.statusResovler).InteralStatusesByCodebaseIDAndCommitID(ctx, r.ch.CodebaseID, *r.ch.CommitID)
}

func (r *ChangeResolver) DownloadTarGz(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.root.downloadsResovler.InternalContentsDownloadTarGzUrl(ctx, r.ch)
}

func (r *ChangeResolver) DownloadZip(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.root.downloadsResovler.InternalContentsDownloadZipUrl(ctx, r.ch)
}

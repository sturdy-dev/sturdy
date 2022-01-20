package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	vcsvcs "getsturdy.com/api/vcs"
	"sync"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	"getsturdy.com/api/pkg/change/service"
	"getsturdy.com/api/pkg/change/vcs"
	db_comments "getsturdy.com/api/pkg/comments/db"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/vcs/executor"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type ChangeRootResolver struct {
	svc *service.Service

	changeRepo       db_change.Repository
	changeCommitRepo db_change.CommitRepository
	commentsRepo     db_comments.Repository

	authService *service_auth.Service

	commentResolver *resolvers.CommentRootResolver
	authorResolver  resolvers.AuthorRootResolver
	statusResovler  *resolvers.StatusesRootResolver

	executorProvider executor.Provider

	logger *zap.Logger
}

func NewResolver(
	svc *service.Service,

	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,
	commentsRepo db_comments.Repository,

	authService *service_auth.Service,

	commentResolver *resolvers.CommentRootResolver,
	authorResolver resolvers.AuthorRootResolver,
	statusResovler *resolvers.StatusesRootResolver,

	executorProvider executor.Provider,

	logger *zap.Logger,
) resolvers.ChangeRootResolver {
	return &ChangeRootResolver{
		svc: svc,

		changeRepo:       changeRepo,
		changeCommitRepo: changeCommitRepo,
		commentsRepo:     commentsRepo,

		authService: authService,

		commentResolver: commentResolver,
		authorResolver:  authorResolver,
		statusResovler:  statusResovler,

		executorProvider: executorProvider,

		logger: logger,
	}
}

func (r *ChangeRootResolver) Change(ctx context.Context, args resolvers.ChangeArgs) (resolvers.ChangeResolver, error) {
	var changeID change.ID
	if args.ID != nil {
		// Lookup by ChangeID
		changeID = change.ID(*args.ID)
	} else if args.CommitID != nil && args.CodebaseID != nil {
		// Lookup by CommitID and CodebaseID
		changeCommit, err := r.changeCommitRepo.GetByCommitID(string(*args.CommitID), string(*args.CodebaseID))
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to lookup by commit id and codebaseid: %w", err))
		}
		changeID = changeCommit.ChangeID
	} else {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "neither id nor commitID is set")
	}

	// Get change
	ch, err := r.changeRepo.Get(changeID)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to lookup by id: %w", err))
	}

	if err := r.authService.CanRead(ctx, ch); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &ChangeResolver{root: r, ch: ch}, nil
}

type ChangeResolver struct {
	ch   change.Change
	root *ChangeRootResolver

	changeOnTrunk    change.ChangeCommit
	changeOnTrunkErr error
	getChangeOnTrunk sync.Once
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

func (r *ChangeResolver) loadChangeOnTrunk() {
	r.changeOnTrunk, r.changeOnTrunkErr = r.root.changeCommitRepo.GetByChangeIDOnTrunk(r.ch.ID)
}

func (r *ChangeResolver) TrunkCommitID() (*string, error) {
	r.getChangeOnTrunk.Do(r.loadChangeOnTrunk)
	if errors.Is(r.changeOnTrunkErr, sql.ErrNoRows) {
		return nil, nil
	}
	if r.changeOnTrunkErr != nil {
		return nil, gqlerrors.Error(r.changeOnTrunkErr)
	}
	return &r.changeOnTrunk.CommitID, nil
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
	r.getChangeOnTrunk.Do(r.loadChangeOnTrunk)
	if errors.Is(r.changeOnTrunkErr, sql.ErrNoRows) {
		return nil, nil
	}
	if r.changeOnTrunkErr != nil {
		return nil, gqlerrors.Error(r.changeOnTrunkErr)
	}

	var diffs []string
	err := r.root.executorProvider.New().Git(func(repo vcsvcs.Repo) error {
		var err error
		diffs, err = vcs.GetDiffs(repo, r.changeOnTrunk.CommitID)
		if err != nil {
			return err
		}
		return nil
	}).ExecTrunk(
		r.ch.CodebaseID,
		"changeResolverDiffs",
	)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	allower, err := r.root.authService.GetAllower(ctx, r.ch)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	decoratedDiff, err := unidiff.NewUnidiff(
		unidiff.NewStringsPatchReader(diffs),
		r.root.logger,
	).WithAllower(allower).Decorate()
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make([]resolvers.FileDiffResolver, len(decoratedDiff))
	for k, v := range decoratedDiff {
		res[k] = &fileDiffResolver{diff: v}
	}
	return res, nil
}

func (r *ChangeResolver) Statuses(ctx context.Context) ([]resolvers.StatusResolver, error) {
	changeCommit, err := r.root.changeCommitRepo.GetByChangeIDOnTrunk(r.ch.ID)
	switch {
	case err == nil:
		return (*r.root.statusResovler).InteralStatusesByCodebaseIDAndCommitID(ctx, changeCommit.CodebaseID, changeCommit.CommitID)
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(fmt.Errorf("failed to lookup by commit id and codebaseid: %w", err))
	}
}

func (r *ChangeResolver) DownloadTarGz(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.download(ctx, service.ArchiveFormatTarGz)
}

func (r *ChangeResolver) DownloadZip(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.download(ctx, service.ArchiveFormatZip)
}

func (r *ChangeResolver) download(ctx context.Context, format service.ArchiveFormat) (resolvers.ContentsDownloadUrlResolver, error) {
	r.getChangeOnTrunk.Do(r.loadChangeOnTrunk)
	if errors.Is(r.changeOnTrunkErr, sql.ErrNoRows) {
		return nil, nil
	}
	if r.changeOnTrunkErr != nil {
		return nil, gqlerrors.Error(r.changeOnTrunkErr)
	}

	if err := r.root.authService.CanRead(ctx, r.ch); err != nil {
		return nil, gqlerrors.Error(err)
	}

	allower, err := r.root.authService.GetAllower(ctx, r.ch)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	url, err := r.root.svc.CreateArchive(ctx, allower, r.changeOnTrunk.CodebaseID, r.changeOnTrunk.CommitID, format)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &download{url: url}, nil
}

type download struct {
	url string
}

func (d *download) ID() graphql.ID {
	return graphql.ID(d.url)
}

func (d *download) URL() string {
	return d.url
}

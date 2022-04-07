package graphql

import (
	"context"
	"errors"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/changes/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type ChangeResolver struct {
	ch   *changes.Change
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

func (r *ChangeResolver) Statuses(ctx context.Context) ([]resolvers.ChangeStatusResolver, error) {
	return (*r.root.statusResovler).InternalChangeStatuses(ctx, r.ch)
}

func (r *ChangeResolver) DownloadTarGz(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.root.downloadsResovler.InternalChangeDownloadTarGzUrl(ctx, r.ch)
}

func (r *ChangeResolver) DownloadZip(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.root.downloadsResovler.InternalChangeDownloadZipUrl(ctx, r.ch)
}

func (r *ChangeResolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	if r.ch.WorkspaceID == nil {
		return nil, nil
	}
	yes := true
	return (*r.root.workspaceResolver).Workspace(ctx, resolvers.WorkspaceArgs{
		ID:            graphql.ID(*r.ch.WorkspaceID),
		AllowArchived: &yes,
	})
}

func (r *ChangeResolver) Codebase(ctx context.Context) (resolvers.CodebaseResolver, error) {
	id := graphql.ID(r.ch.CodebaseID)
	return (*r.root.codebaseResolver).Codebase(ctx, resolvers.CodebaseArgs{
		ID: &id,
	})
}

func (r *ChangeResolver) Child(ctx context.Context) (resolvers.ChangeResolver, error) {
	if child, err := r.root.svc.ChildChange(ctx, r.ch); errors.Is(err, service.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	} else {
		return &ChangeResolver{
			ch:   child,
			root: r.root,
		}, nil
	}
}

func (r *ChangeResolver) Parent(ctx context.Context) (resolvers.ChangeResolver, error) {
	if parent, err := r.root.svc.ParentChange(ctx, r.ch); errors.Is(err, service.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	} else {
		return &ChangeResolver{
			ch:   parent,
			root: r.root,
		}, nil
	}
}

func (r *ChangeResolver) Activity(ctx context.Context, args resolvers.ActivityArgs) ([]resolvers.ActivityResolver, error) {
	return r.root.activityResovler.InternalActivityByChangeID(ctx, r.ch.ID, args)
}

package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/unidiff"

	"github.com/graph-gophers/graphql-go"
)

type ChangeRootResolver interface {
	IntenralListChanges(ctx context.Context, codebaseID string, limit int, before *graphql.ID) ([]ChangeResolver, error)

	Change(ctx context.Context, args ChangeArgs) (ChangeResolver, error)
}

type ChangeArgs struct {
	ID         *graphql.ID
	CommitID   *graphql.ID
	CodebaseID *graphql.ID
}

type ChangeResolver interface {
	ID() graphql.ID
	Comments() ([]TopCommentResolver, error)
	Title() string
	Description() string
	TrunkCommitID() (*string, error)
	Author(context.Context) (AuthorResolver, error)
	CreatedAt() int32
	Diffs(context.Context) ([]FileDiffResolver, error)
	Statuses(context.Context) ([]StatusResolver, error)
	Workspace(context.Context) (WorkspaceResolver, error)
	Codebase(context.Context) (CodebaseResolver, error)
	Activity(context.Context, ActivityArgs) ([]ActivityResolver, error)

	Parent(context.Context) (ChangeResolver, error)
	Child(context.Context) (ChangeResolver, error)

	DownloadTarGz(context.Context) (ContentsDownloadUrlResolver, error)
	DownloadZip(context.Context) (ContentsDownloadUrlResolver, error)
}

type FileDiffRootResolver interface {
	// Internal
	InternalFileDiff(prefix string, diff *unidiff.FileDiff) FileDiffResolver
}

type FileDiffResolver interface {
	ID() graphql.ID
	OrigName() string
	NewName() string
	PreferredName() string
	IsDeleted() bool
	IsNew() bool
	IsMoved() bool

	IsLarge() bool
	LargeFileInfo() (LargeFileInfoResolver, error)

	IsHidden() bool

	Hunks() ([]HunkResolver, error)
}

type LargeFileInfoResolver interface {
	ID() graphql.ID
	Size() int32
}

type HunkResolver interface {
	ID() graphql.ID
	Patch() string
	IsOutdated() bool
	IsApplied() bool
	IsDismissed() bool
}

type ContentsDownloadUrlRootResolver interface {
	// Internal
	InternalContentsDownloadTarGzUrl(context.Context, *changes.Change) (ContentsDownloadUrlResolver, error)
	InternalContentsDownloadZipUrl(context.Context, *changes.Change) (ContentsDownloadUrlResolver, error)
}

type ContentsDownloadUrlResolver interface {
	ID() graphql.ID
	URL() string
}

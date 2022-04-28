package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/workspaces"

	"github.com/graph-gophers/graphql-go"
)

type ChangeRootResolver interface {
	InternalListChanges(ctx context.Context, codebaseID codebases.ID, limit int, before *graphql.ID) ([]ChangeResolver, error)

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
	Statuses(context.Context) ([]ChangeStatusResolver, error)
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
	InternalFileDiffWithWorkspace(keyPrefix string, diff *unidiff.FileDiff, workspace *workspaces.Workspace) FileDiffResolver
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

	OldFileInfo() FileInfoResolver
	NewFileInfo() FileInfoResolver
}

type LargeFileInfoResolver interface {
	ID() graphql.ID
	Size() int32
}

type HunkResolver interface {
	ID() graphql.ID
	HunkID() graphql.ID
	Patch() string
	IsOutdated() bool
	IsApplied() bool
	IsDismissed() bool
}

type ContentsDownloadUrlRootResolver interface {
	// Internal
	InternalChangeDownloadTarGzUrl(context.Context, *changes.Change) (ContentsDownloadUrlResolver, error)
	InternalChangeDownloadZipUrl(context.Context, *changes.Change) (ContentsDownloadUrlResolver, error)

	InternalWorkspaceDownloadTarGzUrl(context.Context, *workspaces.Workspace, DownloadArchiveArgs) (ContentsDownloadUrlResolver, error)
	InternalWorkspaceDownloadZipUrl(context.Context, *workspaces.Workspace, DownloadArchiveArgs) (ContentsDownloadUrlResolver, error)
}

type ContentsDownloadUrlResolver interface {
	ID() graphql.ID
	URL() string
}

type FileType string

const (
	FileTypeUnknown FileType = ""
	FileTypeText    FileType = "Text"
	FileTypeBinary  FileType = "Binary"
	FileTypeImage   FileType = "Image"
)

type FileInfoResolver interface {
	ID() graphql.ID
	RawURL(ctx context.Context) *string
	FileType(ctx context.Context) (FileType, error)
}

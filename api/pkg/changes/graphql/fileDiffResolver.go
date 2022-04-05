package graphql

import (
	"fmt"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/workspaces"
)

type fileDiffRootResolver struct {
	fileRootResolver resolvers.FileRootResolver
}

func NewFileDiffRootResolver(fileRootResolver resolvers.FileRootResolver) resolvers.FileDiffRootResolver {
	return &fileDiffRootResolver{
		fileRootResolver: fileRootResolver,
	}
}

func (r *fileDiffRootResolver) InternalFileDiff(keyPrefix string, diff *unidiff.FileDiff) resolvers.FileDiffResolver {
	return &fileDiffResolver{
		root:      r,
		keyPrefix: keyPrefix,
		diff:      *diff,
	}
}

func (r *fileDiffRootResolver) InternalFileDiffWithWorkspace(keyPrefix string, diff *unidiff.FileDiff, workspace *workspaces.Workspace) resolvers.FileDiffResolver {
	return &fileDiffResolver{
		root:      r,
		keyPrefix: keyPrefix,
		diff:      *diff,
		workspace: workspace,
	}
}

type fileDiffResolver struct {
	root      *fileDiffRootResolver
	keyPrefix string
	diff      unidiff.FileDiff
	workspace *workspaces.Workspace
}

func (f *fileDiffResolver) ID() graphql.ID {
	return graphql.ID(f.keyPrefix + "_" + f.diff.PreferredName)
}

func (f *fileDiffResolver) OrigName() string {
	return f.diff.OrigName
}

func (f *fileDiffResolver) NewName() string {
	return f.diff.NewName
}

func (f *fileDiffResolver) PreferredName() string {
	return f.diff.PreferredName
}

func (f *fileDiffResolver) IsDeleted() bool {
	return f.diff.IsDeleted
}

func (f *fileDiffResolver) IsNew() bool {
	return f.diff.IsNew
}

func (f *fileDiffResolver) IsMoved() bool {
	return f.diff.IsMoved
}

func (f *fileDiffResolver) IsLarge() bool {
	return f.diff.IsLarge
}

func (f *fileDiffResolver) LargeFileInfo() (resolvers.LargeFileInfoResolver, error) {
	if f.diff.LargeFileInfo == nil {
		return nil, nil
	}
	return &largeFileInfoResolver{
		id:   f.ID(),
		info: f.diff.LargeFileInfo,
	}, nil
}

func (f *fileDiffResolver) IsHidden() bool {
	return f.diff.IsHidden
}

func (f *fileDiffResolver) Hunks() ([]resolvers.HunkResolver, error) {
	res := make([]resolvers.HunkResolver, len(f.diff.Hunks), len(f.diff.Hunks))
	for k, v := range f.diff.Hunks {
		res[k] = &hunkResolver{
			id:   fmt.Sprintf("%s-%d", f.ID(), k),
			hunk: v,
		}
	}
	return res, nil
}

func (f *fileDiffResolver) NewFileInfo() resolvers.FileInfoResolver {
	// not supported outside of workspaces yet
	if f.workspace == nil {
		return nil
	}
	return f.root.fileRootResolver.InternalFileInfoInWorkspace(f.ID()+"_new", f.diff.NewName, f.workspace, true)
}

func (f *fileDiffResolver) OldFileInfo() resolvers.FileInfoResolver {
	// not supported outside of workspaces yet
	if f.workspace == nil {
		return nil
	}
	return f.root.fileRootResolver.InternalFileInfoInWorkspace(f.ID()+"_old", f.diff.OrigName, f.workspace, false)
}

type hunkResolver struct {
	id   string
	hunk unidiff.Hunk
}

func (h *hunkResolver) ID() graphql.ID {
	return graphql.ID(h.id)
}

func (h *hunkResolver) HunkID() graphql.ID {
	return graphql.ID(h.hunk.ID)
}

func (h *hunkResolver) Patch() string {
	return h.hunk.Patch
}

func (h *hunkResolver) IsOutdated() bool {
	return h.hunk.IsOutdated
}

func (h *hunkResolver) IsApplied() bool {
	return h.hunk.IsApplied
}

func (h *hunkResolver) IsDismissed() bool {
	return h.hunk.IsDismissed
}

type largeFileInfoResolver struct {
	id   graphql.ID
	info *unidiff.LargeFileInfo
}

func (l *largeFileInfoResolver) ID() graphql.ID {
	return l.id
}

func (l *largeFileInfoResolver) Size() int32 {
	return int32(l.info.Size)
}

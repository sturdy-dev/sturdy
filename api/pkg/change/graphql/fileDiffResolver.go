package graphql

import (
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/unidiff"

	"github.com/graph-gophers/graphql-go"
)

type FileDiffRootResolver struct{}

func NewFileDiffRootResolver() resolvers.FileDiffRootResolver {
	return &FileDiffRootResolver{}
}

func (r *FileDiffRootResolver) InternalFileDiff(diff *unidiff.FileDiff) resolvers.FileDiffResolver {
	return &fileDiffResolver{diff: *diff}
}

type fileDiffResolver struct {
	diff unidiff.FileDiff
}

func (f *fileDiffResolver) ID() graphql.ID {
	return graphql.ID(f.diff.PreferredName)
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
		res[k] = &hunkResolver{hunk: v}
	}
	return res, nil
}

type hunkResolver struct {
	hunk unidiff.Hunk
}

func (h *hunkResolver) ID() graphql.ID {
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

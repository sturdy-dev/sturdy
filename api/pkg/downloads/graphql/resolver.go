package graphql

import (
	"context"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/workspaces"
)

type ContentsDownloadURLRootResolver struct{}

func New() resolvers.ContentsDownloadUrlRootResolver {
	return &ContentsDownloadURLRootResolver{}
}

func (*ContentsDownloadURLRootResolver) InternalWorkspaceDownloadTarGzUrl(context.Context, *workspaces.Workspace, resolvers.DownloadArchiveArgs) (resolvers.ContentsDownloadUrlResolver, error) {
	return nil, errors.ErrNotImplemented
}

func (*ContentsDownloadURLRootResolver) InternalWorkspaceDownloadZipUrl(context.Context, *workspaces.Workspace, resolvers.DownloadArchiveArgs) (resolvers.ContentsDownloadUrlResolver, error) {
	return nil, errors.ErrNotImplemented
}

func (*ContentsDownloadURLRootResolver) InternalChangeDownloadTarGzUrl(context.Context, *changes.Change) (resolvers.ContentsDownloadUrlResolver, error) {
	return nil, errors.ErrNotImplemented
}

func (*ContentsDownloadURLRootResolver) InternalChangeDownloadZipUrl(context.Context, *changes.Change) (resolvers.ContentsDownloadUrlResolver, error) {
	return nil, errors.ErrNotImplemented
}

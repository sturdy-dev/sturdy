package graphql

import (
	"context"

	"getsturdy.com/api/pkg/change"
	"getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type ContentsDownloadURLRootResolver struct{}

func New() resolvers.ContentsDownloadUrlRootResolver {
	return &ContentsDownloadURLRootResolver{}
}

func (*ContentsDownloadURLRootResolver) InternalContentsDownloadTarGzUrl(context.Context, *change.Change, *change.ChangeCommit) (resolvers.ContentsDownloadUrlResolver, error) {
	return nil, errors.ErrNotImplemented
}

func (*ContentsDownloadURLRootResolver) InternalContentsDownloadZipUrl(context.Context, *change.Change, *change.ChangeCommit) (resolvers.ContentsDownloadUrlResolver, error) {
	return nil, errors.ErrNotImplemented
}

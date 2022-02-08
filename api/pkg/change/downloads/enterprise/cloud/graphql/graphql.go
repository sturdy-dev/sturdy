package graphql

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/change"
	service_downloads "getsturdy.com/api/pkg/change/downloads/enterprise/cloud/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type ContentsDownloadURLRootResolver struct {
	service     *service_downloads.Service
	authService *service_auth.Service
}

func New(
	service *service_downloads.Service,
	authService *service_auth.Service,
) resolvers.ContentsDownloadUrlRootResolver {
	return &ContentsDownloadURLRootResolver{
		service:     service,
		authService: authService,
	}
}

func (r *ContentsDownloadURLRootResolver) InternalContentsDownloadTarGzUrl(ctx context.Context, change *change.Change, changeCommit *change.ChangeCommit) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.download(ctx, change, changeCommit, service_downloads.ArchiveFormatTarGz)
}

func (r *ContentsDownloadURLRootResolver) InternalContentsDownloadZipUrl(ctx context.Context, change *change.Change, changeCommit *change.ChangeCommit) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.download(ctx, change, changeCommit, service_downloads.ArchiveFormatZip)
}

func (r *ContentsDownloadURLRootResolver) download(ctx context.Context, change *change.Change, changeCommit *change.ChangeCommit, format service_downloads.ArchiveFormat) (resolvers.ContentsDownloadUrlResolver, error) {
	if err := r.authService.CanRead(ctx, change); err != nil {
		return nil, gqlerrors.Error(err)
	}

	allower, err := r.authService.GetAllower(ctx, change)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	url, err := r.service.CreateArchive(ctx, allower, change.CodebaseID, changeCommit.CommitID, format)
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

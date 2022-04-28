package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/changes"
	service_downloads "getsturdy.com/api/pkg/downloads/enterprise/cloud/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/workspaces"

	"github.com/graph-gophers/graphql-go"
)

type ContentsDownloadURLRootResolver struct {
	service     *service_downloads.Service
	authService *service_auth.Service
	snapshotter *service_snapshots.Service
}

func New(
	service *service_downloads.Service,
	authService *service_auth.Service,
	snapshotter *service_snapshots.Service,
) resolvers.ContentsDownloadUrlRootResolver {
	return &ContentsDownloadURLRootResolver{
		service:     service,
		authService: authService,
		snapshotter: snapshotter,
	}
}

func (r *ContentsDownloadURLRootResolver) InternalWorkspaceDownloadTarGzUrl(ctx context.Context, workspace *workspaces.Workspace) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.downloadWorkspace(ctx, workspace, service_downloads.ArchiveFormatTarGz)
}

func (r *ContentsDownloadURLRootResolver) InternalWorkspaceDownloadZipUrl(ctx context.Context, workspace *workspaces.Workspace) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.downloadWorkspace(ctx, workspace, service_downloads.ArchiveFormatZip)
}

func (r *ContentsDownloadURLRootResolver) InternalChangeDownloadTarGzUrl(ctx context.Context, change *changes.Change) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.downloadChange(ctx, change, service_downloads.ArchiveFormatTarGz)
}

func (r *ContentsDownloadURLRootResolver) InternalChangeDownloadZipUrl(ctx context.Context, change *changes.Change) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.downloadChange(ctx, change, service_downloads.ArchiveFormatZip)
}

func (r *ContentsDownloadURLRootResolver) downloadWorkspace(ctx context.Context, workspace *workspaces.Workspace, format service_downloads.ArchiveFormat) (resolvers.ContentsDownloadUrlResolver, error) {
	if err := r.authService.CanRead(ctx, workspace); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if workspace.LatestSnapshotID == nil {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "workspace has no snapshots")
	}

	allower, err := r.authService.GetAllower(ctx, workspace)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	snapshot, err := r.snapshotter.GetByID(ctx, *workspace.LatestSnapshotID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	url, err := r.service.CreateArchive(ctx, allower, snapshot.CodebaseID, fmt.Sprintf("snapshot-%s", snapshot.ID), snapshot.CommitSHA, format)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &download{url: url}, nil
}

func (r *ContentsDownloadURLRootResolver) downloadChange(ctx context.Context, change *changes.Change, format service_downloads.ArchiveFormat) (resolvers.ContentsDownloadUrlResolver, error) {
	if err := r.authService.CanRead(ctx, change); err != nil {
		return nil, gqlerrors.Error(err)
	}

	allower, err := r.authService.GetAllower(ctx, change)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	url, err := r.service.CreateArchive(ctx, allower, change.CodebaseID, "sturdytrunk", *change.CommitID, format)
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

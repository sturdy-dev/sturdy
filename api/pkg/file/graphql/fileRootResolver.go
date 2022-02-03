package graphql

import (
	"context"
	"strings"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebase"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/multierr"
)

type fileRootResolver struct {
	executorProvider executor.Provider
	authService      *service_auth.Service
}

func NewFileRootResolver(
	executorProvider executor.Provider,
	authService *service_auth.Service,
) resolvers.FileRootResolver {
	return &fileRootResolver{
		executorProvider: executorProvider,
		authService:      authService,
	}
}

func (r *fileRootResolver) InternalFile(ctx context.Context, codebaseID string, pathsWithFallback ...string) (resolvers.FileOrDirectoryResolver, error) {
	var resolver resolvers.FileOrDirectoryResolver

	allower, err := r.authService.GetAllower(ctx, &codebase.Codebase{ID: codebaseID})
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	err = r.executorProvider.New().Git(func(repo vcs.Repo) error {
		headCommit, err := repo.HeadCommit()
		if err != nil {
			return multierr.Combine(gqlerrors.ErrNotFound, err)
		}

		for _, path := range pathsWithFallback {
			// Git is case-sensitive, try some alternative variants of the file name if not found on the first attempt
			variants := []string{
				strings.TrimLeft(path, "/"),
				path,
				strings.ToLower(path),
				strings.ToUpper(path),
			}
			for _, variantName := range variants {

				contents, err := repo.FileContentsAtCommit(headCommit.Id().String(), variantName)
				if err == nil && allower.IsAllowed(variantName, false) {
					resolver = &fileResolver{
						codebaseID: codebaseID,
						path:       variantName,
						contents:   contents,
					}
					return nil
				}

				children, err := repo.DirectoryChildrenAtCommit(headCommit.Id().String(), variantName)
				if err == nil {
					resolver = &directoryResolver{
						codebaseID:   codebaseID,
						path:         variantName,
						children:     children,
						rootResolver: r,
					}
					return nil
				}
			}
		}

		return gqlerrors.ErrNotFound
	}).ExecTrunk(codebaseID, "fileRootResolver")
	if err != nil {
		return nil, err
	}

	return resolver, nil
}

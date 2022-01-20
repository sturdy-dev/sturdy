package graphql

import (
	"context"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"strings"

	"go.uber.org/multierr"
)

type fileRootResolver struct {
	executorProvider executor.Provider
}

func NewFileRootResolver(executorProvider executor.Provider) resolvers.FileRootResolver {
	return &fileRootResolver{
		executorProvider: executorProvider,
	}
}

func (r *fileRootResolver) InternalFile(ctx context.Context, codebaseID string, pathsWithFallback ...string) (resolvers.FileOrDirectoryResolver, error) {
	var resolver resolvers.FileOrDirectoryResolver

	err := r.executorProvider.New().Git(func(repo vcs.Repo) error {
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
				if err == nil {
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

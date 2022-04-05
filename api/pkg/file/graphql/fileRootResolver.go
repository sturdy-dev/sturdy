package graphql

import (
	"context"
	"strings"

	"github.com/graph-gophers/graphql-go"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebases"
	service_file "getsturdy.com/api/pkg/file/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/multierr"
)

type fileRootResolver struct {
	executorProvider executor.Provider
	authService      *service_auth.Service
	fileService      *service_file.Service
}

func NewFileRootResolver(
	executorProvider executor.Provider,
	authService *service_auth.Service,
	fileService *service_file.Service,
) resolvers.FileRootResolver {
	return &fileRootResolver{
		executorProvider: executorProvider,
		authService:      authService,
		fileService:      fileService,
	}
}

func (r *fileRootResolver) InternalFile(ctx context.Context, codebase *codebases.Codebase, pathsWithFallback ...string) (resolvers.FileOrDirectoryResolver, error) {
	var resolver resolvers.FileOrDirectoryResolver

	allower, err := r.authService.GetAllower(ctx, codebase)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	err = r.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
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
						codebaseID: codebase.ID,
						path:       variantName,
						contents:   contents,
					}
					return nil
				}

				children, err := repo.DirectoryChildrenAtCommit(headCommit.Id().String(), variantName)
				if err == nil {
					resolver = &directoryResolver{
						codebase:     codebase,
						path:         variantName,
						children:     children,
						rootResolver: r,
					}
					return nil
				}
			}
		}

		return gqlerrors.ErrNotFound
	}).ExecTrunk(codebase.ID, "fileRootResolver")
	if err != nil {
		return nil, err
	}

	return resolver, nil
}

func (r *fileRootResolver) InternalFileInfoInWorkspace(id graphql.ID, filePath string, workspace *workspaces.Workspace, isNew bool) resolvers.FileInfoResolver {
	return &fileInfoResolver{
		root:      r,
		id:        id,
		filePath:  filePath,
		workspace: workspace,
		isNew:     isNew,
	}
}

package graphql

import (
	"context"
	"errors"
	"fmt"

	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/sync"
	"getsturdy.com/api/pkg/sync/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type rootResolver struct {
	fileDiffRootResolver resolvers.FileDiffRootResolver
	executorProvider     executor.Provider
	logger               *zap.Logger
	workspaceService     *service_workspace.Service
}

func NewRootResolver(
	fileDiffRootResolver resolvers.FileDiffRootResolver,
	executorProvider executor.Provider,
	logger *zap.Logger,
	workspaceService *service_workspace.Service,
) resolvers.RebaseStatusRootResolver {
	return &rootResolver{
		fileDiffRootResolver: fileDiffRootResolver,
		executorProvider:     executorProvider,
		logger:               logger,
		workspaceService:     workspaceService,
	}
}

func (r *rootResolver) InternalWorkspaceRebaseStatus(ctx context.Context, workspaceID string) (resolvers.RebaseStatusResolver, error) {
	workspace, _ := r.workspaceService.GetByID(ctx, workspaceID)
	if workspace.ViewID == nil {
		return nil, nil
	}
	var status *sync.RebaseStatusResponse
	if err := r.executorProvider.New().
		AllowRebasingState(). // allowed to be able to get the status if rebasing is in progress
		Write(func(repo vcs.RepoWriter) error {
			rebasing, err := repo.OpenRebase()
			if err != nil {
				if errors.Is(err, vcs.ErrNoRebaseInProgress) {
					status = &sync.RebaseStatusResponse{}
					return nil
				}
				return fmt.Errorf("failed to open rebase: %w", err)
			}
			rebaseStatus, err := service.Status(r.logger, rebasing)
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}
			status = rebaseStatus
			return nil
		}).ExecView(workspace.CodebaseID, *workspace.ViewID, "rebaseStatus"); err != nil {
		return nil, gqlerror.Error(fmt.Errorf("failed to get status: %w", err))
	}
	return &resolver{
		id:                   workspaceID,
		status:               status,
		fileDiffRootResolver: &r.fileDiffRootResolver,
	}, nil
}

type resolver struct {
	id                   string
	status               *sync.RebaseStatusResponse
	fileDiffRootResolver *resolvers.FileDiffRootResolver
}

type ConflictingFileResolver struct {
	id                   string
	conflictingFile      sync.ConflictingFile
	fileDiffRootResolver *resolvers.FileDiffRootResolver
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.id)
}

func (r *resolver) IsRebasing() bool {
	return r.status.IsRebasing
}

func (r *resolver) ConflictingFiles() ([]resolvers.ConflictingFileResolver, error) {
	var elements []resolvers.ConflictingFileResolver
	for _, element := range r.status.ConflictingFiles {
		elements = append(elements, &ConflictingFileResolver{
			id:                   r.id,
			conflictingFile:      element,
			fileDiffRootResolver: r.fileDiffRootResolver,
		})
	}
	return elements, nil
}

func (c *ConflictingFileResolver) ID() graphql.ID {
	return graphql.ID(c.id + c.conflictingFile.Path)
}

func (c *ConflictingFileResolver) Path() string {
	return c.conflictingFile.Path
}

func (c *ConflictingFileResolver) WorkspaceDiff() (resolvers.FileDiffResolver, error) {
	return (*c.fileDiffRootResolver).InternalFileDiff(string(c.ID())+"Workspace", &c.conflictingFile.WorkspaceDiff), nil
}

func (c *ConflictingFileResolver) TrunkDiff() (resolvers.FileDiffResolver, error) {
	return (*c.fileDiffRootResolver).InternalFileDiff(string(c.ID())+"Trunk", &c.conflictingFile.TrunkDiff), nil
}

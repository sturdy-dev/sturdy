package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/changes"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_remote "getsturdy.com/api/pkg/remote/enterprise/service"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
)

type Service struct {
	*service_workspaces.WorkspaceService

	gitHubService *service_github.Service
	remoteService *service_remote.EnterpriseService
}

var _ service_workspaces.Service = (*Service)(nil)

func New(
	ossService *service_workspaces.WorkspaceService,
	gitHubService *service_github.Service,
	remoteService *service_remote.EnterpriseService,
) *Service {
	return &Service{
		WorkspaceService: ossService,
		gitHubService:    gitHubService,
		remoteService:    remoteService,
	}
}

func (s *Service) LandChange(ctx context.Context, ws *workspaces.Workspace, patchIDs []string, diffOpts ...vcs.DiffOption) (*changes.Change, error) {
	gitHubRepository, err := s.gitHubService.GetRepositoryByCodebaseID(ctx, ws.CodebaseID)
	switch {
	case err == nil, errors.Is(err, sql.ErrNoRows):
	default:
		return nil, fmt.Errorf("failed to get gitHubRepository: %w", err)
	}

	if gitHubRepository != nil && gitHubRepository.IntegrationEnabled && gitHubRepository.GitHubSourceOfTruth {
		return nil, fmt.Errorf("landing disallowed when a github integration exists for codebase (github is source of truth)")
	}

	change, err := s.WorkspaceService.LandChange(ctx, ws, patchIDs, diffOpts...)
	if err != nil {
		return nil, err
	}

	if gitHubRepository != nil && gitHubRepository.IntegrationEnabled && !gitHubRepository.GitHubSourceOfTruth {
		// TODO: move to a queue.
		if err := s.gitHubService.Push(ctx, gitHubRepository, change); err != nil {
			return nil, fmt.Errorf("failed to push to github: %w", err)
		}
		return change, nil
	}

	return change, nil
}

func (s *Service) Push(ctx context.Context, user *users.User, ws *workspaces.Workspace) error {
	// if codebase has github integration, push to github
	_, err := s.gitHubService.CreateOrUpdatePullRequest(ctx, user, ws)
	switch {
	case errors.Is(err, service_github.ErrIntegrationNotEnabled):
	// continue, check push to other provider
	case err != nil:
		return fmt.Errorf("failed to push to github: %w", err)
	default:
		return nil
	}

	if err := s.remoteService.Push(ctx, user, ws); err != nil {
		return fmt.Errorf("failed to push to remote: %w", err)
	}

	return nil
}

package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/change"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	"getsturdy.com/api/pkg/workspaces"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
)

type Service struct {
	*service_workspaces.WorkspaceService

	gitHubService *service_github.Service
}

func New(
	ossService *service_workspaces.WorkspaceService,
	gitHubService *service_github.Service,
) *Service {
	return &Service{
		WorkspaceService: ossService,
		gitHubService:    gitHubService,
	}
}

func (s *Service) LandChange(ctx context.Context, ws *workspaces.Workspace, patchIDs []string, diffOpts ...vcs.DiffOption) (*change.Change, error) {
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

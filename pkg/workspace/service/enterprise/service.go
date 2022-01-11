package enterprise

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mash/pkg/change"
	service_github "mash/pkg/github/service"
	"mash/pkg/workspace"
	service_workspace "mash/pkg/workspace/service"
	"mash/vcs"
)

type Service struct {
	*service_workspace.WorkspaceService

	gitHubService *service_github.Service
}

func New(
	ossService *service_workspace.WorkspaceService,
	gitHubService *service_github.Service,
) *Service {
	return &Service{
		WorkspaceService: ossService,
		gitHubService:    gitHubService,
	}
}

func (s *Service) LandChange(ctx context.Context, ws *workspace.Workspace, patchIDs []string, diffOpts ...vcs.DiffOption) (*change.Change, error) {
	gitHubRepository, err := s.gitHubService.GetRepositoryByCodebaseID(ctx, ws.CodebaseID)
	switch {
	case err == nil, errors.Is(err, sql.ErrNoRows):
	default:
		return nil, fmt.Errorf("failed to get gitHubRepository: %w", err)
	}

	integrationEnabled := gitHubRepository != nil && gitHubRepository.IntegrationEnabled
	isGitHubASourceOfTruth := integrationEnabled && gitHubRepository.GitHubSourceOfTruth
	if isGitHubASourceOfTruth {
		return nil, fmt.Errorf("landing disallowed when a github integration exists for codebase (github is source of truth)")
	}

	change, err := s.WorkspaceService.LandChange(ctx, ws, patchIDs, diffOpts...)
	if err != nil {
		return nil, err
	}

	if integrationEnabled && !isGitHubASourceOfTruth {
		return change, nil
	}

	// TODO: move to a queue.
	if err := s.gitHubService.Push(ctx, gitHubRepository, change); err != nil {
		return nil, fmt.Errorf("failed to push to github: %w", err)
	}

	return change, nil
}

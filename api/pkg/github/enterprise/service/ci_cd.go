package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"getsturdy.com/api/pkg/codebases"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	github_vcs "getsturdy.com/api/pkg/github/enterprise/vcs"
	"getsturdy.com/api/vcs"
)

func (svc *Service) CreateBuild(ctx context.Context, codebaseID codebases.ID, snapshotCommitSha, branchName string) error {
	gitHubRepository, err := svc.GetRepositoryByCodebaseID(ctx, codebaseID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("failed to get repository: %w", err)
	case !gitHubRepository.IntegrationEnabled:
		return nil
	}

	installation, err := svc.gitHubInstallationRepo.GetByInstallationID(gitHubRepository.InstallationID)
	if err != nil {
		return fmt.Errorf("failed to get github installation: %w", err)
	}

	logger := svc.logger.With(
		zap.Int64("github_installation_id", gitHubRepository.InstallationID),
		zap.Int64("github_repository_id", gitHubRepository.GitHubRepositoryID),
	)

	accessToken, err := github_client.GetAccessToken(
		ctx,
		logger,
		svc.gitHubAppConfig,
		installation,
		gitHubRepository.GitHubRepositoryID,
		svc.gitHubRepositoryRepo,
		svc.gitHubInstallationClientProvider,
	)
	if err != nil {
		return fmt.Errorf("failed to get github access token: %w", err)
	}

	// Push in a git executor context
	var userVisibleError string
	if err := svc.executorProvider.New().GitWrite(func(repo vcs.RepoGitWriter) error {
		localBranchName := fmt.Sprintf("push-%s", snapshotCommitSha)
		if err := repo.CreateNewBranchAt(localBranchName, snapshotCommitSha); err != nil {
			return fmt.Errorf("failed to create new branch: %w", err)
		}

		userVisibleError, err = github_vcs.PushToGitHubWithRefspec(repo, accessToken, fmt.Sprintf("+refs/heads/%s:refs/heads/%s", localBranchName, branchName))
		if err != nil {
			return err
		}

		if err := repo.DeleteBranch(localBranchName); err != nil {
			return fmt.Errorf("failed to delete branch: %w", err)
		}

		return nil
	}).ExecTrunk(codebaseID, "pushToGitHub"); err != nil {
		logger.Error("failed to push to github", zap.Error(err))
		// save that the push failed
		t := time.Now()
		gitHubRepository.LastPushAt = &t
		gitHubRepository.LastPushErrorMessage = &userVisibleError
		if err := svc.gitHubRepositoryRepo.Update(gitHubRepository); err != nil {
			logger.Error("failed to update status of github integration", zap.Error(err))
		}

		return fmt.Errorf("failed to push to github: %w", err)
	}

	return nil
}

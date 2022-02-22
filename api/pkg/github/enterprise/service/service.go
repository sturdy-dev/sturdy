package service

import (
	"context"
	"fmt"
	"time"

	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/change"
	service_change "getsturdy.com/api/pkg/change/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/github"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	config_github "getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	github_vcs "getsturdy.com/api/pkg/github/enterprise/vcs"
	"getsturdy.com/api/pkg/notification/sender"
	db_review "getsturdy.com/api/pkg/review/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	service_sync "getsturdy.com/api/pkg/sync/service"
	service_user "getsturdy.com/api/pkg/users/service"
	sender_workspace_activity "getsturdy.com/api/pkg/workspaces/activity/sender"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type ImporterQueue interface {
	Enqueue(ctx context.Context, codebaseID string, userID string) error
}

type ClonerQueue interface {
	Enqueue(context.Context, *github.CloneRepositoryEvent) error
}

type Service struct {
	logger *zap.Logger

	gitHubRepositoryRepo   db_github.GitHubRepositoryRepo
	gitHubInstallationRepo db_github.GitHubInstallationRepo
	gitHubUserRepo         db_github.GitHubUserRepo
	gitHubPullRequestRepo  db_github.GitHubPRRepo

	gitHubPullRequestImporterQueue *ImporterQueue
	gitHubCloneQueue               *ClonerQueue

	gitHubAppConfig                  *config_github.GitHubAppConfig
	gitHubInstallationClientProvider github_client.InstallationClientProvider
	gitHubPersonalClientProvider     github_client.PersonalClientProvider
	gitHubAppClientProvider          github_client.AppClientProvider

	workspaceWriter  db_workspaces.WorkspaceWriter
	workspaceReader  db_workspaces.WorkspaceReader
	codebaseUserRepo db_codebase.CodebaseUserRepository
	codebaseRepo     db_codebase.CodebaseRepository
	reviewRepo       db_review.ReviewRepository

	executorProvider executor.Provider

	snap               snapshotter.Snapshotter
	analyticsService   *service_analytics.Service
	notificationSender sender.NotificationSender
	eventsSender       events.EventSender
	activitySender     sender_workspace_activity.ActivitySender

	userService     service_user.Service
	syncService     *service_sync.Service
	commentsService *service_comments.Service
	changeService   *service_change.Service

	buildQueue *workers_ci.BuildQueue
}

func New(
	logger *zap.Logger,

	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubInstallationRepo db_github.GitHubInstallationRepo,
	gitHubUserRepo db_github.GitHubUserRepo,
	gitHubPullRequestRepo db_github.GitHubPRRepo,
	gitHubAppConfig *config_github.GitHubAppConfig,
	gitHubInstallationClientProvider github_client.InstallationClientProvider,
	gitHubPersonalClientProvider github_client.PersonalClientProvider,
	gitHubAppClientProvider github_client.AppClientProvider,

	gitHubPullRequestImporterQueue *ImporterQueue,
	gitHubCloneQueue *ClonerQueue,

	workspaceWriter db_workspaces.WorkspaceWriter,
	workspaceReader db_workspaces.WorkspaceReader,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	codebaseRepo db_codebase.CodebaseRepository,
	reviewRepo db_review.ReviewRepository,

	executorProvider executor.Provider,
	snap snapshotter.Snapshotter,
	analyticsService *service_analytics.Service,
	notificationSender sender.NotificationSender,
	eventsSender events.EventSender,
	activitySender sender_workspace_activity.ActivitySender,

	userService service_user.Service,
	syncService *service_sync.Service,
	commentsService *service_comments.Service,
	changeService *service_change.Service,

	buildQueue *workers_ci.BuildQueue,
) *Service {
	return &Service{
		logger: logger,

		gitHubRepositoryRepo:             gitHubRepositoryRepo,
		gitHubInstallationRepo:           gitHubInstallationRepo,
		gitHubUserRepo:                   gitHubUserRepo,
		gitHubPullRequestRepo:            gitHubPullRequestRepo,
		gitHubAppConfig:                  gitHubAppConfig,
		gitHubInstallationClientProvider: gitHubInstallationClientProvider,
		gitHubPersonalClientProvider:     gitHubPersonalClientProvider,
		gitHubAppClientProvider:          gitHubAppClientProvider,

		gitHubPullRequestImporterQueue: gitHubPullRequestImporterQueue,
		gitHubCloneQueue:               gitHubCloneQueue,

		workspaceWriter:  workspaceWriter,
		workspaceReader:  workspaceReader,
		codebaseUserRepo: codebaseUserRepo,
		codebaseRepo:     codebaseRepo,
		reviewRepo:       reviewRepo,

		executorProvider:   executorProvider,
		snap:               snap,
		analyticsService:   analyticsService,
		notificationSender: notificationSender,
		eventsSender:       eventsSender,
		activitySender:     activitySender,

		userService:     userService,
		syncService:     syncService,
		commentsService: commentsService,
		changeService:   changeService,

		buildQueue: buildQueue,
	}
}

func (s *Service) GetRepositoryByCodebaseID(_ context.Context, codebaseID string) (*github.GitHubRepository, error) {
	return s.gitHubRepositoryRepo.GetByCodebaseID(codebaseID)
}

func (s *Service) GetRepositoryByInstallationAndRepoID(_ context.Context, installationID, repositoryID int64) (*github.GitHubRepository, error) {
	return s.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repositoryID)
}

func (s *Service) Push(ctx context.Context, gitHubRepository *github.GitHubRepository, change *change.Change) error {
	installation, err := s.gitHubInstallationRepo.GetByInstallationID(gitHubRepository.InstallationID)
	if err != nil {
		return fmt.Errorf("failed to get github installation: %w", err)
	}

	logger := s.logger.With(
		zap.Int64("github_installation_id", gitHubRepository.InstallationID),
		zap.Int64("github_repository_id", gitHubRepository.GitHubRepositoryID),
	)

	accessToken, err := github_client.GetAccessToken(
		ctx,
		logger,
		s.gitHubAppConfig,
		installation,
		gitHubRepository.GitHubRepositoryID,
		s.gitHubRepositoryRepo,
		s.gitHubInstallationClientProvider,
	)
	if err != nil {
		return fmt.Errorf("failed to get github access token: %w", err)
	}

	t := time.Now()

	// GitHub Repository might have been modified at this point, refresh it
	gitHubRepository, err = s.gitHubRepositoryRepo.GetByID(gitHubRepository.ID)
	if err != nil {
		return fmt.Errorf("failed to re-get github repository: %w", err)
	}

	// Push in a git executor context
	var userVisibleError string
	if err := s.executorProvider.New().GitWrite(func(repo vcs.RepoGitWriter) error {
		userVisibleError, err = github_vcs.PushTrackedToGitHub(
			logger,
			repo,
			accessToken,
			gitHubRepository.TrackedBranch,
		)
		if err != nil {
			return err
		}
		return nil
	}).ExecTrunk(change.CodebaseID, "landChangePushTrackedToGitHub"); err != nil {
		logger.Error("failed to push to github (sturdy is source of truth)", zap.Error(err))
		// save that the push failed
		gitHubRepository.LastPushAt = &t
		gitHubRepository.LastPushErrorMessage = &userVisibleError
		if err := s.gitHubRepositoryRepo.Update(gitHubRepository); err != nil {
			logger.Error("failed to update status of github integration", zap.Error(err))
		}

		return fmt.Errorf("failed to push to github: %w", err)
	}

	// Mark as successfully pushed
	gitHubRepository.LastPushAt = &t
	gitHubRepository.LastPushErrorMessage = nil
	if err := s.gitHubRepositoryRepo.Update(gitHubRepository); err != nil {
		return fmt.Errorf("failed to update status of github integration: %w", err)
	}

	logger.Info("pushed to github")

	return nil
}

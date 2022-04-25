package service

import (
	"context"
	"fmt"
	"time"

	sender_workspace_activity "getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/changes"
	service_change "getsturdy.com/api/pkg/changes/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	config_github "getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	github_vcs "getsturdy.com/api/pkg/github/enterprise/vcs"
	"getsturdy.com/api/pkg/notification/sender"
	service_remote "getsturdy.com/api/pkg/remote/enterprise/service"
	db_review "getsturdy.com/api/pkg/review/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/users"
	service_user "getsturdy.com/api/pkg/users/service"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	gitHubRepositoryRepo   db_github.GitHubRepositoryRepository
	gitHubInstallationRepo db_github.GitHubInstallationRepository
	gitHubUserRepo         db_github.GitHubUserRepository
	gitHubPullRequestRepo  db_github.GitHubPRRepository

	gitHubPullRequestImporterQueue *ImporterQueue
	gitHubCloneQueue               *ClonerQueue

	gitHubAppConfig                  *config_github.GitHubAppConfig
	gitHubInstallationClientProvider github_client.InstallationClientProvider
	gitHubPersonalClientProvider     github_client.PersonalClientProvider
	gitHubAppClientProvider          github_client.AppClientProvider

	workspaceWriter  db_workspaces.WorkspaceWriter
	workspaceReader  db_workspaces.WorkspaceReader
	codebaseUserRepo db_codebases.CodebaseUserRepository
	codebaseRepo     db_codebases.CodebaseRepository
	reviewRepo       db_review.ReviewRepository

	executorProvider executor.Provider

	snap               *service_snapshots.Service
	analyticsService   *service_analytics.Service
	notificationSender sender.NotificationSender
	eventsSender       events.EventSender
	activitySender     sender_workspace_activity.ActivitySender
	eventsPublisher    *eventsv2.Publisher

	userService       service_user.Service
	syncService       *service_sync.Service
	commentsService   *service_comments.Service
	changeService     *service_change.Service
	remoteService     *service_remote.EnterpriseService
	workspacesService *service_workspaces.Service
	activityService   *service_activity.Service

	buildQueue *workers_ci.BuildQueue
}

func New(
	logger *zap.Logger,

	gitHubRepositoryRepo db_github.GitHubRepositoryRepository,
	gitHubInstallationRepo db_github.GitHubInstallationRepository,
	gitHubUserRepo db_github.GitHubUserRepository,
	gitHubPullRequestRepo db_github.GitHubPRRepository,
	gitHubAppConfig *config_github.GitHubAppConfig,
	gitHubInstallationClientProvider github_client.InstallationClientProvider,
	gitHubPersonalClientProvider github_client.PersonalClientProvider,
	gitHubAppClientProvider github_client.AppClientProvider,

	importerQueue *ImporterQueue,
	clonerQueue *ClonerQueue,

	workspaceWriter db_workspaces.WorkspaceWriter,
	workspaceReader db_workspaces.WorkspaceReader,
	codebaseUserRepo db_codebases.CodebaseUserRepository,
	codebaseRepo db_codebases.CodebaseRepository,
	reviewRepo db_review.ReviewRepository,

	executorProvider executor.Provider,
	snap *service_snapshots.Service,
	analyticsService *service_analytics.Service,
	notificationSender sender.NotificationSender,
	eventsSender events.EventSender,
	activitySender sender_workspace_activity.ActivitySender,
	eventsPublisher *eventsv2.Publisher,

	userService service_user.Service,
	syncService *service_sync.Service,
	commentsService *service_comments.Service,
	changeService *service_change.Service,
	remoteService *service_remote.EnterpriseService,
	workspacesService *service_workspaces.Service,
	activityService *service_activity.Service,

	buildQueue *workers_ci.BuildQueue,
) *Service {
	svc := &Service{
		logger: logger,

		gitHubRepositoryRepo:             gitHubRepositoryRepo,
		gitHubInstallationRepo:           gitHubInstallationRepo,
		gitHubUserRepo:                   gitHubUserRepo,
		gitHubPullRequestRepo:            gitHubPullRequestRepo,
		gitHubAppConfig:                  gitHubAppConfig,
		gitHubInstallationClientProvider: gitHubInstallationClientProvider,
		gitHubPersonalClientProvider:     gitHubPersonalClientProvider,
		gitHubAppClientProvider:          gitHubAppClientProvider,

		gitHubPullRequestImporterQueue: importerQueue,
		gitHubCloneQueue:               clonerQueue,

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
		eventsPublisher:    eventsPublisher,

		userService:       userService,
		syncService:       syncService,
		commentsService:   commentsService,
		changeService:     changeService,
		remoteService:     remoteService,
		workspacesService: workspacesService,
		activityService:   activityService,

		buildQueue: buildQueue,
	}
	clonerQueue.setService(svc)
	importerQueue.setService(svc)
	return svc
}

func (svc *Service) GetRepositoryByCodebaseID(_ context.Context, codebaseID codebases.ID) (*github.Repository, error) {
	return svc.gitHubRepositoryRepo.GetByCodebaseID(codebaseID)
}

func (svc *Service) GetRepositoryByInstallationAndRepoID(_ context.Context, installationID, repositoryID int64) (*github.Repository, error) {
	return svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repositoryID)
}

func (svc *Service) Push(ctx context.Context, gitHubRepository *github.Repository, change *changes.Change) error {
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

	t := time.Now()

	// GitHub Repository might have been modified at this point, refresh it
	gitHubRepository, err = svc.gitHubRepositoryRepo.GetByID(gitHubRepository.ID)
	if err != nil {
		return fmt.Errorf("failed to re-get github repository: %w", err)
	}

	// Push in a git executor context
	var userVisibleError string
	if err := svc.executorProvider.New().GitWrite(func(repo vcs.RepoGitWriter) error {
		userVisibleError, err = github_vcs.PushTrackedToGitHub(repo, accessToken, gitHubRepository.TrackedBranch)
		if err != nil {
			return err
		}
		return nil
	}).ExecTrunk(change.CodebaseID, "landChangePushTrackedToGitHub"); err != nil {
		logger.Error("failed to push to github (sturdy is source of truth)", zap.Error(err))
		// save that the push failed
		gitHubRepository.LastPushAt = &t
		gitHubRepository.LastPushErrorMessage = &userVisibleError
		if err := svc.gitHubRepositoryRepo.Update(gitHubRepository); err != nil {
			logger.Error("failed to update status of github integration", zap.Error(err))
		}

		return fmt.Errorf("failed to push to github: %w", err)
	}

	// Mark as successfully pushed
	gitHubRepository.LastPushAt = &t
	gitHubRepository.LastPushErrorMessage = nil
	if err := svc.gitHubRepositoryRepo.Update(gitHubRepository); err != nil {
		return fmt.Errorf("failed to update status of github integration: %w", err)
	}

	logger.Info("pushed to github")

	return nil
}

func (svc *Service) CheckPermissions(ctx context.Context) (bool, []string, []string, error) {
	provider, err := svc.gitHubAppClientProvider(svc.gitHubAppConfig)
	if err != nil {
		return false, nil, nil, fmt.Errorf("failed to get github app client provider: %w", err)
	}

	get, _, err := provider.Get(ctx, "")
	if err != nil {
		return false, nil, nil, fmt.Errorf("failed to get github app: %w", err)
	}

	var missingPermissions []string
	insertMissingPermission := func(key string) {
		missingPermissions = append(missingPermissions, key)
	}

	permissions := get.Permissions
	if permissions.GetContents() != "write" {
		insertMissingPermission("Content")
	}
	if permissions.GetMetadata() != "read" {
		insertMissingPermission("Metadata")
	}
	if permissions.GetPullRequests() != "write" {
		insertMissingPermission("Pull Request")
	}
	if permissions.GetStatuses() != "read" {
		insertMissingPermission("Status")
	}
	if permissions.GetWorkflows() != "write" {
		insertMissingPermission("Workflows")
	}

	eventsMap := make(map[string]string)
	eventsMap["pull_request"] = "Pull Request"
	eventsMap["pull_request_review"] = "Pull Request Review"
	eventsMap["push"] = "Push"
	eventsMap["status"] = "Status"
	eventsMap["workflow_job"] = "Workflow Job"

	eventsEnabled := make(map[string]bool)
	for _, e := range get.Events {
		eventsEnabled[e] = true
	}

	var missingEvents []string
	for key, element := range eventsMap {
		if !eventsEnabled[key] {
			missingEvents = append(missingEvents, element)
		}
	}

	hasErrors := len(missingPermissions) > 0 || len(missingEvents) > 0

	return !hasErrors, missingPermissions, missingEvents, nil
}

func (svc *Service) InheritShadowData(
	ctx context.Context,
	gitHubUser *github.User,
	realUser *users.User,
	shadowUser *users.User,
) error {
	if gitHubUser.UserID == realUser.ID {
		// noop
		return nil
	}

	// inherit workspaces
	workspaces, err := svc.workspaceReader.ListByUserID(ctx, gitHubUser.UserID)
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	for _, workspace := range workspaces {
		if err := svc.workspaceWriter.UpdateFields(ctx, workspace.ID, db_workspaces.SetUserID(realUser.ID)); err != nil {
			return fmt.Errorf("failed to update workspace: %w", err)
		}
	}

	// inherit github user
	gitHubUser.UserID = realUser.ID
	if err := svc.gitHubUserRepo.Update(gitHubUser); err != nil {
		return fmt.Errorf("failed to update github user: %w", err)
	}

	if err := svc.userService.Inherit(ctx, realUser.ID, shadowUser); err != nil {
		return fmt.Errorf("failed to deactivate shadow user: %w", err)
	}

	return nil
}

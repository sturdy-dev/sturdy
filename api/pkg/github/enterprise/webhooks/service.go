package webhooks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	gh "github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_change "getsturdy.com/api/pkg/change/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/github/enterprise/client"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	config_github "getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	vcs_github "getsturdy.com/api/pkg/github/enterprise/vcs"
	db_review "getsturdy.com/api/pkg/review/db"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/workspace/activity"
	sender_workspace_activity "getsturdy.com/api/pkg/workspace/activity/sender"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
)

type Service struct {
	logger *zap.Logger

	gitHubPullRequestRepo  db_github.GitHubPRRepo
	gitHubRepositoryRepo   db_github.GitHubRepositoryRepo
	gitHubInstallationRepo db_github.GitHubInstallationRepo

	workspaceWriter db_workspace.WorkspaceWriter
	workspaceReader db_workspace.WorkspaceReader
	reviewRepo      db_review.ReviewRepository

	gitHubAppConfig                  *config_github.GitHubAppConfig
	gitHubInstallationClientProvider github_client.InstallationClientProvider
	gitHubPersonalClientProvider     github_client.PersonalClientProvider
	gitHubAppClientProvider          github_client.AppClientProvider

	executorProvider executor.Provider

	eventsSender     events.EventSender
	analyticsService *service_analytics.Service
	activitySender   sender_workspace_activity.ActivitySender

	syncService      *service_sync.Service
	workspaceService service_workspace.Service
	commentsService  *service_comments.Service
	changeService    *service_change.Service

	buildQueue *workers_ci.BuildQueue
}

func New(
	logger *zap.Logger,

	gitHubPullRequestRepo db_github.GitHubPRRepo,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubInstallationRepo db_github.GitHubInstallationRepo,

	workspaceWriter db_workspace.WorkspaceWriter,
	workspaceReader db_workspace.WorkspaceReader,
	reviewRepo db_review.ReviewRepository,

	gitHubAppConfig *config_github.GitHubAppConfig,
	gitHubInstallationClientProvider github_client.InstallationClientProvider,
	gitHubPersonalClientProvider github_client.PersonalClientProvider,
	gitHubAppClientProvider github_client.AppClientProvider,

	executorProvider executor.Provider,

	eventsSender events.EventSender,
	analyticsService *service_analytics.Service,
	activitySender sender_workspace_activity.ActivitySender,

	syncService *service_sync.Service,
	workspaceService service_workspace.Service,
	commentsService *service_comments.Service,
	changeService *service_change.Service,

	buildQueue *workers_ci.BuildQueue,
) *Service {
	return &Service{
		logger: logger,

		gitHubPullRequestRepo:  gitHubPullRequestRepo,
		gitHubRepositoryRepo:   gitHubRepositoryRepo,
		gitHubInstallationRepo: gitHubInstallationRepo,

		workspaceWriter: workspaceWriter,
		workspaceReader: workspaceReader,
		reviewRepo:      reviewRepo,

		gitHubAppConfig:                  gitHubAppConfig,
		gitHubInstallationClientProvider: gitHubInstallationClientProvider,
		gitHubPersonalClientProvider:     gitHubPersonalClientProvider,
		gitHubAppClientProvider:          gitHubAppClientProvider,

		executorProvider: executorProvider,

		eventsSender:     eventsSender,
		analyticsService: analyticsService,
		activitySender:   activitySender,

		syncService:      syncService,
		workspaceService: workspaceService,
		commentsService:  commentsService,
		changeService:    changeService,

		buildQueue: buildQueue,
	}
}

func (svc *Service) HandlePullRequestEvent(event *gh.PullRequestEvent) error {
	ctx := context.Background()

	apiPR := event.GetPullRequest()
	pr, err := svc.gitHubPullRequestRepo.GetByGitHubID(apiPR.GetID())
	if errors.Is(err, sql.ErrNoRows) {
		return nil // noop
	} else if err != nil {
		return fmt.Errorf("failed to get github pull request from db: %w", err)
	}

	if apiPR.GetState() == "closed" || apiPR.GetState() == "open" {
		pr.Open = apiPR.GetState() == "open"
		t0 := time.Now()
		pr.UpdatedAt = &t0
		pr.ClosedAt = apiPR.ClosedAt
		pr.Merged = apiPR.GetMerged()
		pr.MergedAt = apiPR.MergedAt
		if err := svc.gitHubPullRequestRepo.Update(pr); err != nil {
			return fmt.Errorf("failed to update github pull request in db: %w", err)
		}
	}

	ws, err := svc.workspaceReader.Get(pr.WorkspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		svc.logger.Warn("handled a github pull request webhook for non-existing workspace", zap.String("workspace_id", pr.WorkspaceID), zap.String("github_pr_id", pr.ID), zap.String("github_pr_link", apiPR.GetHTMLURL()))
		return nil // noop
	} else if err != nil {
		return fmt.Errorf("failed to get workspace from db: %w", err)
	}

	// import / sync workspace
	if apiPR.GetState() == "closed" && apiPR.GetMerged() {

		accessToken, err := svc.accessToken(event.GetInstallation().GetID(), event.GetRepo().GetID())
		if err != nil {
			return fmt.Errorf("failed to get github access token: %w", err)
		}

		repo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(event.GetInstallation().GetID(), event.GetRepo().GetID())
		if err != nil {
			return fmt.Errorf("failed to get github repo from db: %w", err)
		}

		// pull from github if sturdy doesn't have the commits
		err = svc.pullFromGitHubIfCommitNotExists(ws.CodebaseID, []string{
			apiPR.GetMergeCommitSHA(),
			event.GetPullRequest().GetBase().GetSHA(),
		}, accessToken, repo.TrackedBranch)
		if err != nil {
			return fmt.Errorf("failed to pullFromGitHubIfCommitNotExists: %w", err)
		}

		ch, err := svc.changeService.CreateWithCommitAsParent(ctx, ws, apiPR.GetMergeCommitSHA(), event.GetPullRequest().GetBase().GetSHA())
		if err != nil {
			return fmt.Errorf("failed to create change: %w", err)
		}

		// unset the draft description
		ws.DraftDescription = ""
		if err := svc.workspaceWriter.Update(ctx, ws); err != nil {
			return fmt.Errorf("failed to update workspace: %w", err)
		}

		hasConflicts, err := svc.workspaceService.HasConflicts(ctx, ws)
		if err != nil {
			svc.logger.Error("failed to check for conflicts", zap.Error(err), zap.Any("workspace_id", ws.ID))
			// do not fail
		} else if err == nil {
			// sync workspace with head if possible
			if !hasConflicts && ws.ViewID != nil {
				if _, err := svc.syncService.OnTrunk(ctx, *ws.ViewID); err != nil {
					return fmt.Errorf("failed to sync workspace: %w", err)
				}
				if err := svc.eventsSender.Codebase(ws.CodebaseID, events.ViewUpdated, *ws.ViewID); err != nil {
					svc.logger.Error("failed to send workspace updated event", zap.Error(err))
					// do not fail
				}
			}
		}

		if err := svc.eventsSender.Workspace(ws.ID, events.WorkspaceUpdated, ws.ID); err != nil {
			svc.logger.Error("failed to send workspace updated event", zap.Error(err))
			// do not fail
		}

		if err := svc.workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(ws.CodebaseID); err != nil {
			return fmt.Errorf("failed to unset up to date with trunk for all in codebase: %w", err)
		}

		if err := svc.commentsService.MoveCommentsFromWorkspaceToChange(ctx, ws.ID, ch.ID); err != nil {
			return fmt.Errorf("failed to migrate comments: %w", err)
		}

		if err := svc.reviewRepo.DismissAllInWorkspace(ctx, ws.ID); err != nil {
			return fmt.Errorf("failed to dissmiss reviews: %w", err)
		}

		svc.analyticsService.Capture(ctx, "pull request merged",
			analytics.DistinctID(ws.UserID),
			analytics.CodebaseID(ws.CodebaseID),
			analytics.Property("workspace_id", ws.ID),
			analytics.Property("github", true),
		)

		// Create workspace activity
		if err := svc.activitySender.Codebase(ctx, ws.CodebaseID, ws.ID, ws.UserID, activity.WorkspaceActivityTypeCreatedChange, string(ch.ID)); err != nil {
			return fmt.Errorf("failed to create workspace activity: %w", err)
		}

		// Send events that the codebase has been updated
		if err := svc.eventsSender.Codebase(ws.CodebaseID, events.CodebaseUpdated, ws.CodebaseID); err != nil {
			svc.logger.Error("failed to send codebase event", zap.Error(err))
			// do not fail
		}

		if err := svc.buildQueue.EnqueueChange(ctx, ch); err != nil {
			svc.logger.Error("failed to enqueue change", zap.Error(err))
			// do not fail
		}
	}

	if err := svc.eventsSender.Codebase(ws.CodebaseID, events.GitHubPRUpdated, pr.WorkspaceID); err != nil {
		svc.logger.Error("failed to send codebase event", zap.String("workspace_id", pr.WorkspaceID), zap.String("github_pr_id", pr.ID), zap.String("github_pr_link", apiPR.GetHTMLURL()), zap.Error(err))
		// do not fail
	}

	return nil
}

func (svc *Service) HandlePushEvent(event *gh.PushEvent) error {
	installationID := event.GetInstallation().GetID()
	repoID := event.GetRepo().GetID()
	repo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repoID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil
	default:
		return fmt.Errorf("failed to get github repo from db: %w", err)
	}

	if event.GetRef() != fmt.Sprintf("refs/heads/%s", repo.TrackedBranch) {
		return nil
	}

	logger := svc.logger.With(zap.String("codebase_id", repo.CodebaseID),
		zap.String("repo_id", repo.ID),
		zap.String("repo_tracked_branch", repo.TrackedBranch))

	if !repo.GitHubSourceOfTruth || !repo.IntegrationEnabled {
		logger.Info("skipping github push event, the integration is disabled or github is not the source of truth")
		return nil
	}

	accessToken, err := svc.accessToken(installationID, repoID)
	if err != nil {
		if strings.Contains(err.Error(), "The permissions requested are not granted to this installation") {
			logger.Info("did not have permissions to get a github token")
			return nil
		}
		return fmt.Errorf("failed to get access token: %w", err)
	}

	if err := svc.executorProvider.New().
		GitWrite(vcs_github.FetchTrackedToSturdytrunk(accessToken, event.GetRef())).
		ExecTrunk(repo.CodebaseID, "githubPushEvent"); err != nil {
		return fmt.Errorf("failed to fetch changes from github: %w", err)
	}

	repo, err = svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repoID)
	if err != nil {
		return fmt.Errorf("failed to get github repo from db: %w", err)
	}

	t := time.Now()
	repo.SyncedAt = &t
	if err := svc.gitHubRepositoryRepo.Update(repo); err != nil {
		return fmt.Errorf("failed to update github repository in db: %w", err)
	}

	// Allow all workspaces to be rebased/synced on the latest head
	if err := svc.workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(repo.CodebaseID); err != nil {
		return fmt.Errorf("failed to unset up to date with trunk for all in codebase: %w", err)
	}

	return nil
}

func (svc *Service) pullFromGitHubIfCommitNotExists(codebaseID string, commitShas []string, accessToken, trackedBranchName string) error {
	shouldPull := false

	if err := svc.executorProvider.New().
		GitRead(func(repo vcs.RepoGitReader) error {
			for _, sha := range commitShas {
				if _, err := repo.GetCommitDetails(sha); err != nil {
					shouldPull = true
				}
			}
			return nil
		}).
		ExecTrunk(codebaseID, "pullFromGitHubIfCommitNotExists.Check"); err != nil {
		return fmt.Errorf("failed to fetch changes from github: %w", err)
	}

	if !shouldPull {
		return nil
	}

	if err := svc.executorProvider.New().
		GitWrite(vcs_github.FetchTrackedToSturdytrunk(accessToken, "refs/heads/"+trackedBranchName)).
		ExecTrunk(codebaseID, "pullFromGitHubIfCommitNotExists.Pull"); err != nil {
		return fmt.Errorf("failed to fetch changes from github: %w", err)
	}

	return nil
}

func (svc *Service) accessToken(installationID, repositoryID int64) (string, error) {
	repo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repositoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get github repo from db: %w", err)
	}

	installation, err := svc.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return "", fmt.Errorf("could not get installation: %w", err)
	}

	accessToken, err := client.GetAccessToken(context.Background(), svc.logger, svc.gitHubAppConfig, installation, repo.GitHubRepositoryID, svc.gitHubRepositoryRepo, svc.gitHubInstallationClientProvider)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	return accessToken, nil
}

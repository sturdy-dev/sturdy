package push

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	"getsturdy.com/api/pkg/change/decorate"
	"getsturdy.com/api/pkg/change/message"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/github/enterprise/db"
	vcs_github "getsturdy.com/api/pkg/github/enterprise/vcs"
	db_review "getsturdy.com/api/pkg/review/db"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/workspace"
	"getsturdy.com/api/pkg/workspace/activity"
	sender_workspace_activity "getsturdy.com/api/pkg/workspace/activity/sender"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func HandlePushEvent(
	ctx context.Context,
	logger *zap.Logger,
	event *gh.PushEvent,
	githubRepositoryRepo db.GitHubRepositoryRepo,
	githubInstallationRepo db.GitHubInstallationRepo,
	workspaceWriter db_workspace.WorkspaceWriter,
	workspaceReadRepo db_workspace.WorkspaceReader,
	workspaceService service_workspace.Service,
	syncService *service_sync.Service,
	gitHubPRRepo db.GitHubPRRepo,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,
	executorProvider executor.Provider,
	gitHubAppConfig *config.GitHubAppConfig,
	githubClientProvider client.InstallationClientProvider,
	eventsSender events.EventSender,
	analyticsClient analytics.Client,
	reviewRepo db_review.ReviewRepository,
	activitySender sender_workspace_activity.ActivitySender,
	commentsService *service_comments.Service,
	buildQueue *workers_ci.BuildQueue,
) error {
	installationID := event.GetInstallation().GetID()
	repoID := event.GetRepo().GetID()
	repo, err := githubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repoID)
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

	logger = logger.With(zap.String("codebase_id", repo.CodebaseID),
		zap.String("repo_id", repo.ID),
		zap.String("repo_tracked_branch", repo.TrackedBranch))

	if !repo.GitHubSourceOfTruth || !repo.IntegrationEnabled {
		logger.Info("skipping github push event, the integration is disabled or github is not the source of truth")
		return nil
	}

	installation, err := githubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return fmt.Errorf("could not get installation: %w", err)
	}

	accessToken, err := client.GetAccessToken(ctx, logger, gitHubAppConfig, installation, repo.GitHubRepositoryID, githubRepositoryRepo, githubClientProvider)
	if err != nil {
		if strings.Contains(err.Error(), "The permissions requested are not granted to this installation") {
			logger.Info("did not have permissions to get a github token")
			return nil
		}
		return fmt.Errorf("failed to get access token: %w", err)
	}

	var maybeImportedChanges []*vcs.LogEntry
	if err := executorProvider.New().
		GitWrite(vcs_github.FetchTrackedToSturdytrunk(accessToken, event.GetRef())).
		GitWrite(func(repo vcs.RepoGitWriter) error {
			logger.Info("listing imported changes")

			// ListImportedChanges lists the latest 50 commits
			// Normally a merge only has 1 new commit, but octopus merges makes it possible to merge multiple PRs at the same time
			// using 50 commits as a upper limit for now
			changes, err := vcs_github.ListImportedChanges(repo)
			if err != nil {
				return fmt.Errorf("failed to list imported changes: %w", err)
			}
			maybeImportedChanges = changes
			return nil
		}).
		ExecTrunk(repo.CodebaseID, "githubPushEvent"); err != nil {
		return fmt.Errorf("failed to fetch changes from github: %w", err)
	}

	repo, err = githubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repoID)
	if err != nil {
		return fmt.Errorf("failed to get github repo from db: %w", err)
	}

	// Decorate / import new commits
	for _, maybeNewChange := range maybeImportedChanges {
		_, err := changeCommitRepo.GetByCommitID(maybeNewChange.ID, repo.CodebaseID)
		switch {
		case err == nil:
			continue
		case errors.Is(err, sql.ErrNoRows):
			// Create both a change and a change commit
			ch := change.Change{
				ID:         change.ID(uuid.New().String()),
				CodebaseID: repo.CodebaseID,

				GitCreatedAt:    &maybeNewChange.Time,
				GitCreatorEmail: &maybeNewChange.Email,
				GitCreatorName:  &maybeNewChange.Name,
			}

			meta := decorate.ParseCommitMessage(maybeNewChange.RawCommitMessage)

			var ws *workspace.Workspace

			// If the imported commit has a workspace_id set (this is the case when using GitHub Pull Requests), use the workspaces draftDescription as the description
			if len(meta.WorkspaceID) > 0 {
				ws, err = workspaceReadRepo.Get(meta.WorkspaceID)
				if err == nil && ws.CodebaseID == repo.CodebaseID {
					cleanCommitMessage := message.CommitMessage(ws.DraftDescription)
					cleanCommitMessageTitle := strings.Split(cleanCommitMessage, "\n")[0]
					ch.UpdatedDescription = ws.DraftDescription
					ch.Title = &cleanCommitMessageTitle
					ch.UserID = &ws.UserID

					// Set CreatedAt to the time the change was imported to Sturdy
					t := time.Now()
					ch.CreatedAt = &t
				}
			}

			if err := changeRepo.Insert(ch); err != nil {
				return fmt.Errorf("failed to create change: %w", err)
			}

			chCommit := change.ChangeCommit{
				ChangeID:   ch.ID,
				CommitID:   maybeNewChange.ID,
				CodebaseID: repo.CodebaseID,
				Trunk:      true,
			}
			if err := changeCommitRepo.Insert(chCommit); err != nil {
				return fmt.Errorf("failed to create change_commit: %w", err)
			}

			// Reset the workspace draftDescription
			if ws != nil {
				ws.DraftDescription = ""
				if err := workspaceWriter.Update(ws); err != nil {
					return fmt.Errorf("failed to update workspace: %w", err)
				}

				hasConflicts, err := workspaceService.HasConflicts(ctx, ws)
				if err != nil {
					logger.Error("failed to check for conflicts", zap.Error(err), zap.Any("workspace_id", ws.ID))
					// do not fail
				}

				// sync workspace with head if possible
				if !hasConflicts && ws.ViewID != nil {
					if _, err := syncService.OnTrunk(*ws.ViewID); err != nil {
						return fmt.Errorf("failed to sync workspace: %w", err)
					}
					if err := eventsSender.Codebase(ws.CodebaseID, events.ViewUpdated, *ws.ViewID); err != nil {
						logger.Error("failed to send workspace updated event", zap.Error(err))
						// do not fail
					}
				}

				if err := eventsSender.Workspace(ws.ID, events.WorkspaceUpdated, ws.ID); err != nil {
					logger.Error("failed to send workspace updated event", zap.Error(err))
					// do not fail
				}

				if err := workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(repo.CodebaseID); err != nil {
					return fmt.Errorf("failed to unset up to date with trunk for all in codebase: %w", err)
				}

				if err := commentsService.MoveCommentsFromWorkspaceToChange(ctx, meta.WorkspaceID, ch.ID); err != nil {
					return fmt.Errorf("failed to migrate comments: %w", err)
				}

				if err := reviewRepo.DismissAllInWorkspace(ctx, meta.WorkspaceID); err != nil {
					return fmt.Errorf("failed to dissmiss reviews: %w", err)
				}

				if err := analyticsClient.Enqueue(analytics.Capture{
					Event:      "pull request merged",
					DistinctId: ws.UserID,
					Properties: analytics.NewProperties().
						Set("github", true).
						Set("codebase_id", ws.CodebaseID).
						Set("workspace_id", ws.ID),
				}); err != nil {
					logger.Error("analytics failed", zap.Error(err))
				}

				// Create workspace activity
				if err := activitySender.Codebase(ctx, ws.CodebaseID, ws.ID, ws.UserID, activity.WorkspaceActivityTypeCreatedChange, string(ch.ID)); err != nil {
					return fmt.Errorf("failed to create workspace activity: %w", err)
				}

				// Send events that the codebase has been updated
				if err := eventsSender.Codebase(ws.CodebaseID, events.CodebaseUpdated, ws.CodebaseID); err != nil {
					logger.Error("failed to send codebase event", zap.Error(err))
					// do not fail
				}

				if err := buildQueue.EnqueueChange(ctx, &ch); err != nil {
					logger.Error("failed to enqueue change", zap.Error(err))
					// do not fail
				}
			}

			logger.Info("decorate commit imported from github", zap.Any("change_id", ch.ID), zap.Any("change_commit_id", chCommit.CommitID), zap.Any("meta", meta))

		default:
			return fmt.Errorf("failed to get change commit from db: %w", err)
		}
	}

	t := time.Now()
	repo.SyncedAt = &t
	if err := githubRepositoryRepo.Update(repo); err != nil {
		return fmt.Errorf("failed to update github repository in db: %w", err)
	}

	// If applicable, send an event notifying the PR update
	refTokens := strings.Split(event.GetRef(), "refs/heads/")
	if len(refTokens) == 2 {
		prs, err := gitHubPRRepo.ListByHeadAndRepositoryID(refTokens[1], event.GetRepo().GetID())
		if err != nil {
			return fmt.Errorf("failed to get pr for ref: %w", err)
		}
		if len(prs) > 0 {
			if err := eventsSender.Codebase(repo.CodebaseID, events.GitHubPRUpdated, prs[0].WorkspaceID); err != nil {
				logger.Error("failed to send event", zap.Error(err))
				// do not fail
			}
		}
	}

	// Allow all workspaces to be rebased/synced on the latest head
	if err := workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(repo.CodebaseID); err != nil {
		return fmt.Errorf("failed to unset up to date with trunk for all in codebase: %w", err)
	}

	return nil
}

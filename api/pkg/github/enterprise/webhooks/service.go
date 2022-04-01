package webhooks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/activity"
	sender_workspace_activity "getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_codebases "getsturdy.com/api/pkg/codebases/service"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/api"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	config_github "getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	vcs_github "getsturdy.com/api/pkg/github/enterprise/vcs"
	db_review "getsturdy.com/api/pkg/review/db"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/users"
	service_users "getsturdy.com/api/pkg/users/service"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
)

type Service struct {
	logger *zap.Logger

	gitHubPullRequestRepo  db_github.GitHubPRRepository
	gitHubRepositoryRepo   db_github.GitHubRepositoryRepository
	gitHubInstallationRepo db_github.GitHubInstallationRepository
	gitHubUserRepo         db_github.GitHubUserRepository

	workspaceWriter db_workspaces.WorkspaceWriter
	workspaceReader db_workspaces.WorkspaceReader
	reviewRepo      db_review.ReviewRepository
	codebaseRepo    db_codebases.CodebaseRepository
	viewRepo        db_view.Repository

	gitHubAppConfig                  *config_github.GitHubAppConfig
	gitHubInstallationClientProvider github_client.InstallationClientProvider
	gitHubPersonalClientProvider     github_client.PersonalClientProvider
	gitHubAppClientProvider          github_client.AppClientProvider

	executorProvider executor.Provider

	eventsSender     events.EventSender
	eventsSenderV2   *eventsv2.Publisher
	analyticsService *service_analytics.Service
	activitySender   sender_workspace_activity.ActivitySender

	syncService      *service_sync.Service
	workspaceService service_workspace.Service
	codebaseService  *service_codebases.Service
	commentsService  *service_comments.Service
	activityService  *service_activity.Service
	changeService    *service_change.Service
	statusService    *service_statuses.Service
	githubService    *service_github.Service
	usersService     service_users.Service

	buildQueue *workers_ci.BuildQueue
}

func New(
	logger *zap.Logger,

	gitHubPullRequestRepo db_github.GitHubPRRepository,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepository,
	gitHubInstallationRepo db_github.GitHubInstallationRepository,
	gitHubUserRepo db_github.GitHubUserRepository,

	workspaceWriter db_workspaces.WorkspaceWriter,
	workspaceReader db_workspaces.WorkspaceReader,
	reviewRepo db_review.ReviewRepository,
	codebaseRepo db_codebases.CodebaseRepository,
	viewRepo db_view.Repository,

	gitHubAppConfig *config_github.GitHubAppConfig,
	gitHubInstallationClientProvider github_client.InstallationClientProvider,
	gitHubPersonalClientProvider github_client.PersonalClientProvider,
	gitHubAppClientProvider github_client.AppClientProvider,

	executorProvider executor.Provider,

	eventsSender events.EventSender,
	eventsSenderV2 *eventsv2.Publisher,
	analyticsService *service_analytics.Service,
	activitySender sender_workspace_activity.ActivitySender,

	syncService *service_sync.Service,
	workspaceService service_workspace.Service,
	codebaseService *service_codebases.Service,
	commentsService *service_comments.Service,
	activityService *service_activity.Service,
	changeService *service_change.Service,
	statusService *service_statuses.Service,
	githubService *service_github.Service,
	usersService service_users.Service,

	buildQueue *workers_ci.BuildQueue,
) *Service {
	return &Service{
		logger: logger.Named("github_webhooks"),

		gitHubPullRequestRepo:  gitHubPullRequestRepo,
		gitHubRepositoryRepo:   gitHubRepositoryRepo,
		gitHubInstallationRepo: gitHubInstallationRepo,
		gitHubUserRepo:         gitHubUserRepo,

		workspaceWriter: workspaceWriter,
		workspaceReader: workspaceReader,
		reviewRepo:      reviewRepo,
		codebaseRepo:    codebaseRepo,
		viewRepo:        viewRepo,

		gitHubAppConfig:                  gitHubAppConfig,
		gitHubInstallationClientProvider: gitHubInstallationClientProvider,
		gitHubPersonalClientProvider:     gitHubPersonalClientProvider,
		gitHubAppClientProvider:          gitHubAppClientProvider,

		executorProvider: executorProvider,

		eventsSender:     eventsSender,
		eventsSenderV2:   eventsSenderV2,
		analyticsService: analyticsService,
		activitySender:   activitySender,

		syncService:      syncService,
		workspaceService: workspaceService,
		codebaseService:  codebaseService,
		commentsService:  commentsService,
		activityService:  activityService,
		changeService:    changeService,
		statusService:    statusService,
		githubService:    githubService,
		usersService:     usersService,

		buildQueue: buildQueue,
	}
}

func getPRStatus(apiPR *api.PullRequest) github.PullRequestState {
	switch apiPR.GetState() {
	case "open":
		return github.PullRequestStateOpen
	case "closed":
		if apiPR.GetMerged() {
			return github.PullRequestStateMerged
		}
		return github.PullRequestStateClosed
	default:
		return github.PullRequestStateUnknown
	}
}

func (svc *Service) getPullRequestAuthor(ctx context.Context, repo *github.Repository, event *PullRequestEvent) (*users.User, error) {
	if ghUser, err := svc.gitHubUserRepo.GetByUsername(event.GetPullRequest().GetUser().GetLogin()); errors.Is(err, sql.ErrNoRows) {
		// user with this username doesn't exist yet, will create shadow
	} else if err != nil {
		return nil, fmt.Errorf("failed to get github user: %w", err)
	} else {
		// user exists, return it
		return svc.usersService.GetByID(ctx, ghUser.UserID)
	}

	// make up email from the user's login, similar to what github does
	// see https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-github-user-account/managing-email-preferences/setting-your-commit-email-address
	email := fmt.Sprintf("%d+%s@users.noreply.github.com.com", event.GetPullRequest().GetUser().GetID(), event.GetPullRequest().GetUser().GetLogin())
	name := event.GetPullRequest().GetUser().GetLogin()
	referer := service_users.GitHubPullRequestReferer(event.GetRepo().GetID(), event.GetPullRequest().GetID())
	user, err := svc.usersService.CreateShadow(ctx, email, referer, &name)
	if err != nil {
		return nil, fmt.Errorf("failed to create shadow user: %w", err)
	}

	// save shadow user <-> github connection
	// this will be used to connect shadow user's data with real user once the real use signs up
	if err := svc.gitHubUserRepo.Create(github.User{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Username:  event.GetPullRequest().GetUser().GetLogin(),
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("failed to create github user: %w", err)
	}

	return user, nil
}

func (svc *Service) importNewPullRequest(
	ctx context.Context,
	logger *zap.Logger,
	repo *github.Repository,
	event *PullRequestEvent,
) error {
	user, err := svc.getPullRequestAuthor(ctx, repo, event)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if _, err := svc.codebaseService.AddUser(ctx, repo.CodebaseID, user); err != nil {
		return fmt.Errorf("failed to add user to codebase: %w", err)
	}

	installation, err := svc.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return fmt.Errorf("could not get installation: %w", err)
	}

	accessToken, err := svc.accessToken(ctx, repo.InstallationID, repo.GitHubRepositoryID)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	if err := svc.githubService.ImportPullRequest(user.ID, event.GetPullRequest(), repo, installation, accessToken); errors.Is(err, service_github.ErrAlreadyImported) {
		return nil
	} else if err != nil {
		return err
	} else {
		return nil
	}
}

func (svc *Service) HandlePullRequestEvent(ctx context.Context, event *PullRequestEvent) error {
	logger := svc.logger.With(
		zap.Int64("pr_id", event.GetPullRequest().GetID()),
		zap.String("pr_state", event.GetPullRequest().GetState()),
		zap.Bool("pr_merged", event.GetPullRequest().GetMerged()),
		zap.Int64("installation_id", event.GetInstallation().GetID()),
		zap.Int64("repository_id", event.GetRepo().GetID()),
		zap.String("repository_name", event.GetRepo().GetFullName()),
	)
	start := time.Now()
	defer logger.Info("handle pull request event", zap.Duration("duration", time.Since(start)))

	repo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(event.GetInstallation().GetID(), event.GetRepo().GetID())
	if errors.Is(err, sql.ErrNoRows) {
		logger.Info("repository not found", zap.Error(err))
		return nil // noop
	} else if err != nil {
		logger.Error("failed to get GitHub repository", zap.Error(err))
		return fmt.Errorf("could not get installation: %w", err)
	}

	logger = logger.With(zap.Stringer("codebase_id", repo.CodebaseID), zap.String("gh_repo_id", repo.ID))

	if pr, err := svc.gitHubPullRequestRepo.GetByGitHubIDAndCodebaseID(event.GetPullRequest().GetID(), repo.CodebaseID); errors.Is(err, sql.ErrNoRows) {
		logger.Info("pull request not found, importing")
		if err := svc.importNewPullRequest(ctx, logger, repo, event); err != nil {
			return fmt.Errorf("failed to import new pull request: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to get github pull request from db: %w", err)
	} else {
		logger.Info("pull request found, updating")
		return svc.updateExistingPullRequest(ctx, logger.With(zap.String("pr_id", pr.ID)), repo, event, pr)
	}
}

func (svc *Service) updateExistingPullRequest(
	ctx context.Context,
	logger *zap.Logger,
	repo *github.Repository,
	event *PullRequestEvent,
	pr *github.PullRequest,
) error {
	now := time.Now()
	pr.UpdatedAt = &now
	pr.ClosedAt = event.GetPullRequest().ClosedAt
	pr.MergedAt = event.GetPullRequest().MergedAt
	pr.State = getPRStatus(event.GetPullRequest())
	if err := svc.gitHubPullRequestRepo.Update(ctx, pr); err != nil {
		return fmt.Errorf("failed to update github pull request in db: %w", err)
	}

	defer func() {
		// make sure we send pr updated event after this function returns in any case
		if err := svc.eventsSenderV2.GitHubPRUpdated(ctx, eventsv2.Codebase(pr.CodebaseID), pr); err != nil {
			logger.Error("failed to send codebase event", zap.Error(err))
			// do not fail
		}
	}()

	if pr.State != github.PullRequestStateMerged {
		// no need to sync if PR is not merged
		return nil
	}

	// import / sync workpsace

	ws, err := svc.workspaceReader.Get(pr.WorkspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		svc.logger.Warn("handled a github pull request webhook for non-existing workspace",
			zap.String("workspace_id", pr.WorkspaceID),
			zap.String("github_pr_id", pr.ID),
			zap.String("github_pr_link", event.GetPullRequest().GetHTMLURL()),
		)
		return nil // noop
	} else if err != nil {
		return fmt.Errorf("failed to get workspace from db: %w", err)
	}

	accessToken, err := svc.accessToken(ctx, event.GetInstallation().GetID(), event.GetRepo().GetID())
	if err != nil {
		return fmt.Errorf("failed to get github access token: %w", err)
	}

	// pull from github if sturdy doesn't have the commits
	if err := svc.pullFromGitHubIfCommitNotExists(pr.CodebaseID, []string{
		event.GetPullRequest().GetMergeCommitSHA(),
		event.GetPullRequest().GetBase().GetSHA(),
	}, accessToken, repo.TrackedBranch); err != nil {
		return fmt.Errorf("failed to pullFromGitHubIfCommitNotExists: %w", err)
	}

	ch, err := svc.changeService.CreateWithCommitAsParent(ctx, ws, event.GetPullRequest().GetMergeCommitSHA(), event.GetPullRequest().GetBase().GetSHA())
	if err != nil {
		return fmt.Errorf("failed to create change: %w", err)
	}

	svc.analyticsService.Capture(ctx, "pull request merged",
		analytics.DistinctID(ws.UserID.String()),
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
		analytics.Property("github", true),
	)

	// send change to ci
	if err := svc.buildQueue.EnqueueChange(ctx, ch); err != nil {
		svc.logger.Error("failed to enqueue change", zap.Error(err))
		// do not fail
	}

	// all workspaces that were up to date with trunk are not up to data anymore
	if err := svc.workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(ws.CodebaseID); err != nil {
		return fmt.Errorf("failed to unset up to date with trunk for all in codebase: %w", err)
	}

	// Create workspace activity that it has created a change
	if err := svc.activitySender.Codebase(ctx, ws.CodebaseID, ws.ID, ws.UserID, activity.TypeCreatedChange, string(ch.ID)); err != nil {
		return fmt.Errorf("failed to create workspace activity: %w", err)
	}

	// copy all workspace activities to change activities
	if err := svc.activityService.SetChange(ctx, ws.ID, ch.ID); err != nil {
		return fmt.Errorf("failed to set change: %w", err)
	}

	if err := svc.commentsService.MoveCommentsFromWorkspaceToChange(ctx, ws.ID, ch.ID); err != nil {
		return fmt.Errorf("failed to migrate comments: %w", err)
	}

	// optimisticly archive workspace. if everything is ok with it, it will stay archived
	// if not, we'll unarchive it a later in this function
	if err := svc.workspaceService.ArchiveWithChange(ctx, ws, ch); err != nil {
		return fmt.Errorf("failed to archive workspace: %w", err)
	}

	// ws might have been modified (reload)
	ws, err = svc.workspaceReader.Get(pr.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to reload workpsace: %w", err)
	}

	hasConflicts, err := svc.workspaceService.HasConflicts(ctx, ws)
	if err != nil {
		return fmt.Errorf("failed to check for conflicts: %w", err)
	}

	if hasConflicts {
		// if workspace has conflicts, that probably means that it was changed between the pr was opened and merged
		// in that case, unarchive workspace and make the user fix it
		return svc.workspaceService.Unarchive(ctx, ws)
	}

	// no conflits, so we can sync the workspace with the trunk
	if _, err := svc.syncService.OnTrunk(ctx, ws); err != nil {
		return fmt.Errorf("failed to sync workspace: %w", err)
	}

	diffs, _, err := svc.workspaceService.Diffs(ctx, ws.ID)
	if err != nil {
		return fmt.Errorf("failed to get diffs: %w", err)
	}

	if len(diffs) != 0 {
		// there are some diffs left after the sync. That means that the workspace was changed between the pr was opened
		// and merged. in that case, unarchive workspace and make the user fix it
		return svc.workspaceService.Unarchive(ctx, ws)
	}

	// workspace is synced with with trunk, has no diffs - very good!
	return nil
}

func (svc *Service) HandlePushEvent(ctx context.Context, event *PushEvent) error {
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

	logger := svc.logger.With(zap.Stringer("codebase_id", repo.CodebaseID),
		zap.String("repo_id", repo.ID),
		zap.String("repo_tracked_branch", repo.TrackedBranch))

	if !repo.GitHubSourceOfTruth || !repo.IntegrationEnabled {
		logger.Info("skipping github push event, the integration is disabled or github is not the source of truth")
		return nil
	}

	accessToken, err := svc.accessToken(ctx, installationID, repoID)
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

	// Unset codebase head change cache
	if err := svc.changeService.UnsetHeadChangeCache(repo.CodebaseID); err != nil {
		return fmt.Errorf("failed to unset head change cache: %w", err)
	}

	// Allow all workspaces to be rebased/synced on the latest head
	if err := svc.workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(repo.CodebaseID); err != nil {
		return fmt.Errorf("failed to unset up to date with trunk for all in codebase: %w", err)
	}

	return nil
}

func (svc *Service) pullFromGitHubIfCommitNotExists(codebaseID codebases.ID, commitShas []string, accessToken, trackedBranchName string) error {
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

func (svc *Service) accessToken(ctx context.Context, installationID, repositoryID int64) (string, error) {
	repo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, repositoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get github repo from db: %w", err)
	}

	installation, err := svc.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return "", fmt.Errorf("could not get installation: %w", err)
	}

	accessToken, err := github_client.GetAccessToken(ctx, svc.logger, svc.gitHubAppConfig, installation, repo.GitHubRepositoryID, svc.gitHubRepositoryRepo, svc.gitHubInstallationClientProvider)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	return accessToken, nil
}

func (svc *Service) HandleInstallationEvent(ctx context.Context, event *InstallationEvent) error {
	svc.analyticsService.IdentifyGitHubInstallation(ctx,
		event.GetInstallation().GetID(),
		event.GetInstallation().GetAccount().GetLogin(),
		event.GetInstallation().GetAccount().GetEmail(),
	)

	svc.analyticsService.Capture(ctx, fmt.Sprintf("github installation %s", event.GetAction()),
		analytics.DistinctID(fmt.Sprintf("%d", event.GetInstallation().GetID())),
	)

	if event.GetAction() == "created" ||
		event.GetAction() == "deleted" ||
		event.GetAction() == "new_permissions_accepted" {

		t := time.Now()
		var uninstalledAt *time.Time
		if event.GetAction() == "deleted" {
			uninstalledAt = &t
		}

		// Check if it's already installed (stale data or whatever)
		existing, err := svc.gitHubInstallationRepo.GetByInstallationID(event.Installation.GetID())
		if errors.Is(err, sql.ErrNoRows) {
			// Save new installation
			err := svc.gitHubInstallationRepo.Create(github.Installation{
				ID:                     uuid.NewString(),
				InstallationID:         event.Installation.GetID(),
				Owner:                  event.Installation.GetAccount().GetLogin(), // "sturdy-dev" or "zegl", etc.
				CreatedAt:              t,
				UninstalledAt:          uninstalledAt,
				HasWorkflowsPermission: true,
			})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// Update existing entry
			existing.UninstalledAt = uninstalledAt

			if event.GetAction() == "new_permissions_accepted" {
				existing.HasWorkflowsPermission = true
			}

			err = svc.gitHubInstallationRepo.Update(existing)
			if err != nil {
				return err
			}
		}

		// Save all repositories
		if event.GetAction() == "created" {
			for _, r := range event.Repositories {

				if err := svc.handleInstalledRepository(event.GetInstallation().GetID(), r); err != nil {
					return err
				}
			}
		}
		// TODO: Handle deleted - update installed repositories?

		return nil
	}

	// Unhandled actions:
	// suspend, unsuspend

	return nil
}

func (svc *Service) HandleInstallationRepositoriesEvent(ctx context.Context, event *InstallationRepositoriesEvent) error {
	_, err := svc.gitHubInstallationRepo.GetByInstallationID(event.GetInstallation().GetID())
	// If the original InstallationEvent webhook was missed (otherwise user has to remove and re-add app)
	if errors.Is(err, sql.ErrNoRows) {
		installation := github.Installation{
			ID:                     uuid.NewString(),
			InstallationID:         event.Installation.GetID(),
			Owner:                  event.Installation.Account.GetLogin(), // "sturdy-dev" or "zegl", etc.
			CreatedAt:              time.Now(),
			HasWorkflowsPermission: true,
		}

		if err := svc.gitHubInstallationRepo.Create(installation); err != nil {
			return err
		}
	}

	if event.GetRepositorySelection() == "selected" || event.GetRepositorySelection() == "all" {
		// Add new repos
		for _, r := range event.RepositoriesAdded {
			if err := svc.handleInstalledRepository(event.GetInstallation().GetID(), r); err != nil {
				return err
			}
		}

		// Mark uninstalled repos as uninstalled
		for _, r := range event.RepositoriesRemoved {
			installedRepo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(event.GetInstallation().GetID(), r.GetID())
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else if err != nil {
				svc.logger.Error(
					"failed to mark as uninstalled",
					zap.Error(err),
					zap.Int64("repository_id", r.GetID()),
					zap.Int64("installation_id", event.GetInstallation().GetID()),
				)
				continue
			}

			t := time.Now()
			installedRepo.UninstalledAt = &t

			if err := svc.gitHubRepositoryRepo.Update(installedRepo); err != nil {
				svc.logger.Error(
					"failed to mark as uninstalled",
					zap.Error(err),
					zap.Int64("repository_id", r.GetID()),
					zap.Int64("installation_id", event.GetInstallation().GetID()),
				)
				continue
			}
		}
	}

	return nil
}

func (svc *Service) handleInstalledRepository(installationID int64, ghRepo *api.Repository) error {
	ctx := context.Background()

	// CreateWithCommitAsParent a non-ready codebase (add the initiating user), and put the event on a queue
	logger := svc.logger.With(zap.String("repo_name", ghRepo.GetName()), zap.Int64("installation_id", installationID))
	logger.Info("handleInstalledRepository")

	existingGitHubRepo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, ghRepo.GetID())
	// If GitHub repo already exists (previously installed and then uninstalled) un-archive it
	if err == nil {
		logger.Info("handleInstalledRepository repository already exists", zap.Any("existing_github_repo", existingGitHubRepo))

		// un-archive the codebase if archived
		cb, err := svc.codebaseRepo.GetAllowArchived(existingGitHubRepo.CodebaseID)
		if err != nil {
			logger.Error("failed to get codebase", zap.Error(err))
			return err
		}

		svc.analyticsService.Capture(ctx, "installed repository", analytics.DistinctID(fmt.Sprintf("%d", installationID)),
			analytics.CodebaseID(cb.ID),
			analytics.Property("repo_name", ghRepo.GetName()),
		)

		if cb.ArchivedAt != nil {
			cb.ArchivedAt = nil
			if err := svc.codebaseRepo.Update(cb); err != nil {
				logger.Error("failed to un-archive codebase", zap.Error(err))
				return err
			}
		}

		if existingGitHubRepo.UninstalledAt != nil {
			existingGitHubRepo.UninstalledAt = nil
			err := svc.gitHubRepositoryRepo.Update(existingGitHubRepo)
			if err != nil {
				logger.Error("failed to un-archive github repository entry", zap.Error(err))
				return err
			}
		}

		return nil
	}

	return nil
}

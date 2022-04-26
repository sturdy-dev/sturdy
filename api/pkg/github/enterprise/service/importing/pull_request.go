package importing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	gh "github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"

	"getsturdy.com/api/pkg/activity"
	sender_workspace_activity "getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/changes/message"
	service_change "getsturdy.com/api/pkg/changes/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebases"
	service_comments "getsturdy.com/api/pkg/comments/service"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/api"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	vcs_github "getsturdy.com/api/pkg/github/enterprise/vcs"
	"getsturdy.com/api/pkg/workspaces"
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

	workspaceWriter db_workspaces.WorkspaceWriter

	executorProvider executor.Provider

	analyticsService *service_analytics.Service

	activitySender  sender_workspace_activity.ActivitySender
	eventsPublisher *eventsv2.Publisher

	commentsService *service_comments.Service
	changeService   *service_change.Service

	workspacesService *service_workspaces.Service
	activityService   *service_activity.Service

	buildQueue *workers_ci.BuildQueue

	gitHubService *service_github.Service
}

func New(
	logger *zap.Logger,

	gitHubRepositoryRepo db_github.GitHubRepositoryRepository,
	gitHubInstallationRepo db_github.GitHubInstallationRepository,
	gitHubUserRepo db_github.GitHubUserRepository,
	gitHubPullRequestRepo db_github.GitHubPRRepository,
	workspaceWriter db_workspaces.WorkspaceWriter,

	executorProvider executor.Provider,
	analyticsService *service_analytics.Service,
	activitySender sender_workspace_activity.ActivitySender,
	eventsPublisher *eventsv2.Publisher,

	commentsService *service_comments.Service,
	changeService *service_change.Service,
	workspacesService *service_workspaces.Service,
	activityService *service_activity.Service,

	buildQueue *workers_ci.BuildQueue,
	gitHubService *service_github.Service,
) *Service {
	svc := &Service{
		logger: logger,

		gitHubRepositoryRepo:   gitHubRepositoryRepo,
		gitHubInstallationRepo: gitHubInstallationRepo,
		gitHubUserRepo:         gitHubUserRepo,
		gitHubPullRequestRepo:  gitHubPullRequestRepo,
		workspaceWriter:        workspaceWriter,

		executorProvider: executorProvider,
		analyticsService: analyticsService,
		activitySender:   activitySender,
		eventsPublisher:  eventsPublisher,

		commentsService:   commentsService,
		changeService:     changeService,
		workspacesService: workspacesService,
		activityService:   activityService,

		buildQueue:    buildQueue,
		gitHubService: gitHubService,
	}
	return svc
}

func (svc *Service) MergePullRequest(ctx context.Context, ws *workspaces.Workspace) error {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return err
	}

	existingGitHubUser, err := svc.gitHubUserRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get github user: %w", err)
	}

	if existingGitHubUser.AccessToken == nil {
		return fmt.Errorf("no github access token found for user %s", userID)
	}

	prs, err := svc.gitHubPullRequestRepo.ListOpenedByWorkspace(ws.ID)
	if err != nil {
		return fmt.Errorf("pull request not found: %w", err)
	}
	if len(prs) != 1 {
		return fmt.Errorf("unexpected number of open pull requests found")
	}

	pr := prs[0]

	ghRepo, err := svc.gitHubRepositoryRepo.GetByCodebaseID(pr.CodebaseID)
	if err != nil {
		return fmt.Errorf("gh repo not found: %w", err)
	}

	ghInstallation, err := svc.gitHubInstallationRepo.GetByInstallationID(ghRepo.InstallationID)
	if err != nil {
		return fmt.Errorf("gh installation not found: %w", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *existingGitHubUser.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	userAuthClient := gh.NewClient(tc)

	mergeOpts := &gh.PullRequestOptions{
		CommitTitle: fmt.Sprintf("Merge pull request #%d - %s", pr.GitHubPRNumber, ws.NameOrFallback()),
		// TODO: Do we want to set this to rebase?
		// MergeMethod: "rebase",
	}

	commitMessage := message.CommitMessage(ws.DraftDescription) + "\n\nMerged via Sturdy"

	// update PR state to merging
	previousState := pr.State
	if err := svc.setPRState(ctx, pr, github.PullRequestStateMerging); err != nil {
		return fmt.Errorf("failed to update pull request state: %w", err)
	}

	// check if pr is already merged, without any api error checking
	// actual error will come when trying to merge pr, if needed
	if apiPR, _, err := userAuthClient.PullRequests.Get(ctx, ghInstallation.Owner, ghRepo.Name, pr.GitHubPRNumber); err == nil && apiPR.GetMerged() {
		return svc.UpdatePullRequestFromGitHub(ctx, ghRepo, pr, api.ConvertPullRequest(apiPR), *existingGitHubUser.AccessToken)
	}

	//nolint:contextcheck
	res, resp, err := userAuthClient.PullRequests.Merge(ctx, ghInstallation.Owner, ghRepo.Name, pr.GitHubPRNumber, commitMessage, mergeOpts)
	if err != nil {
		// rollback github pr state
		if err := svc.setPRState(ctx, pr, previousState); err != nil {
			return fmt.Errorf("failed to update pull request state: %w", err)
		}

		var errorResponse *gh.ErrorResponse
		if resp.StatusCode == http.StatusMethodNotAllowed && errors.As(err, &errorResponse) {
			// 405 not allowed
			// This happens if the repo is configured with branch protection rules (require approvals, tests to pass, etc).
			// Proxy GitHub's error message to the end user.
			//
			// Examples:
			// * "failed to merge pr: PUT https://api.github.com/repos/zegl/empty-11/pulls/4/merge: 405 At least 1 approving review is required by reviewers with write access. []"
			return service_github.GitHubUserError{Msg: errorResponse.Message}
		}

		return fmt.Errorf("failed to merge pr: %w", err)
	}

	if !res.GetMerged() {
		// rollback github pr state
		if err := svc.setPRState(ctx, pr, previousState); err != nil {
			return fmt.Errorf("failed to update pull request state: %w", err)
		}
		return fmt.Errorf("pull request was not merged")
	}

	apiPR, _, err := userAuthClient.PullRequests.Get(ctx, ghInstallation.Owner, ghRepo.Name, pr.GitHubPRNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	return svc.UpdatePullRequestFromGitHub(ctx, ghRepo, pr, api.ConvertPullRequest(apiPR), *existingGitHubUser.AccessToken)
}

func (svc *Service) UpdatePullRequestFromGitHub(
	ctx context.Context,
	repo *github.Repository,
	pr *github.PullRequest,
	gitHubPR *api.PullRequest,
	accessToken string,
) error {
	ws, err := svc.workspacesService.GetByID(ctx, pr.WorkspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		svc.logger.Warn("handled a github pull request webhook for non-existing workspace",
			zap.String("workspace_id", pr.WorkspaceID),
			zap.String("github_pr_id", pr.ID),
			zap.String("github_pr_link", gitHubPR.GetHTMLURL()),
		)
		return nil // noop
	} else if err != nil {
		return fmt.Errorf("failed to get workspace from db: %w", err)
	}

	// make context user context to connect all events to the same user
	ctx = auth.NewUserContext(ctx, ws.UserID)

	if shoudUpdateDescription := pr.Importing; shoudUpdateDescription {
		newDescription, err := service_github.DescriptionFromPullRequest(gitHubPR)
		if err != nil {
			return fmt.Errorf("failed to build description: %w", err)
		}
		ws.DraftDescription = newDescription
		if err := svc.workspaceWriter.UpdateFields(ctx, ws.ID, db_workspaces.SetDraftDescription(newDescription)); err != nil {
			return fmt.Errorf("failed to update workspace: %w", err)
		}
	}

	newState := service_github.GetPRState(gitHubPR)

	if shouldUnarchive := newState == github.PullRequestStateOpen && pr.Importing; shouldUnarchive {
		// if pr is open and not importing, unarchive it
		if err := svc.workspacesService.Unarchive(ctx, ws); err != nil {
			return fmt.Errorf("failed to unarchive workspace: %w", err)
		}
	}

	if shouldArchive := newState == github.PullRequestStateClosed && pr.Importing; shouldArchive {
		// if pr is closed and importing, archive it
		if err := svc.gitHubService.UpdatePRFromGitHub(ctx, pr, gitHubPR); err != nil {
			return fmt.Errorf("failed to update pull request state: %w", err)
		}
		return svc.workspacesService.Archive(ctx, ws)
	} else if shouldUpdateState := newState != github.PullRequestStateMerged; shouldUpdateState {
		// if the PR is not merged, fetch the latest PR data from GitHub
		if err := svc.gitHubService.UpdatePullRequest(ctx, pr, gitHubPR, accessToken, ws); err != nil {
			return fmt.Errorf("failed to update pull request: %w", err)
		}
		if err := svc.gitHubService.UpdatePRFromGitHub(ctx, pr, gitHubPR); err != nil {
			return fmt.Errorf("failed to update pull request state: %w", err)
		}
		return nil
	} else {
		// merge pr
		if err := svc.mergePullRequest(ctx, repo, pr, gitHubPR, accessToken); err != nil {
			return fmt.Errorf("failed to merge pull request: %w", err)
		}
		if err := svc.gitHubService.UpdatePRFromGitHub(ctx, pr, gitHubPR); err != nil {
			return fmt.Errorf("failed to update pull request state: %w", err)
		}
		return nil
	}
}

func (svc *Service) mergePullRequest(
	ctx context.Context,
	repo *github.Repository,
	pr *github.PullRequest,
	gitHubPR *api.PullRequest,
	accessToken string,
) error {
	if pr.State == github.PullRequestStateMerging {
		// another goroutine is already merging this PR
		return nil
	}

	ws, err := svc.workspacesService.GetByID(ctx, pr.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	if alreadyMerged := pr.State == github.PullRequestStateMerged && ws.ChangeID != nil; alreadyMerged {
		// this is a noop, pr is already merged
		return nil
	}

	// pull from github if sturdy doesn't have the commits
	if err := svc.pullFromGitHubIfCommitNotExists(pr.CodebaseID, []string{
		gitHubPR.GetMergeCommitSHA(),
		gitHubPR.GetBase().GetSHA(),
	}, accessToken, repo.TrackedBranch); err != nil {
		return fmt.Errorf("failed to pullFromGitHubIfCommitNotExists: %w", err)
	}

	ch, err := svc.changeService.CreateWithCommitAsParent(ctx, ws, gitHubPR.GetMergeCommitSHA(), gitHubPR.GetBase().GetSHA())
	if errors.Is(err, service_change.ErrAlreadyExists) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create change: %w", err)
	}

	svc.analyticsService.Capture(ctx, "pull request merged",
		analytics.DistinctID(ws.UserID.String()),
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
		analytics.Property("github", true),
		analytics.Property("importing", pr.Importing),
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

	// Archive the workspace
	if err := svc.workspacesService.ArchiveWithChange(ctx, ws, ch); err != nil {
		return fmt.Errorf("failed to archive workspace: %w", err)
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

func (svc *Service) setPRState(ctx context.Context, pr *github.PullRequest, state github.PullRequestState) error {
	pr.State = state
	if err := svc.gitHubPullRequestRepo.Update(ctx, pr); err != nil {
		return fmt.Errorf("failed to update pull request: %w", err)
	}

	if err := svc.eventsPublisher.GitHubPRUpdated(ctx, eventsv2.Codebase(pr.CodebaseID), pr); err != nil {
		svc.logger.Error("failed to send codebase event",
			zap.String("workspace_id", pr.WorkspaceID),
			zap.String("github_pr_id", pr.ID),
			zap.Error(err),
		)
		// do not fail
	}
	return nil
}

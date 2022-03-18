package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/changes/message"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/vcs"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/workspaces"
)

type GitHubUserError struct {
	msg string
}

func (g GitHubUserError) Error() string {
	return g.msg
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
		&oauth2.Token{AccessToken: existingGitHubUser.AccessToken},
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
	pr.State = github.PullRequestStateMerging
	if err := svc.gitHubPullRequestRepo.Update(ctx, pr); err != nil {
		return fmt.Errorf("failed to update pull request: %w", err)
	}
	if err := svc.eventsPublisher.GitHubPRUpdated(ctx, events.Codebase(pr.CodebaseID), pr); err != nil {
		svc.logger.Error("failed to send codebase event",
			zap.String("workspace_id", pr.WorkspaceID),
			zap.String("github_pr_id", pr.ID),
			zap.Error(err),
		)
		// do not fail
	}

	// check if pr is already merged, without any api error checking
	// actual error will come when trying to merge pr, if needed
	if apiPR, _, err := userAuthClient.PullRequests.Get(ctx, ghInstallation.Owner, ghRepo.Name, pr.GitHubPRNumber); err == nil && apiPR.GetMerged() {
		// pr is already merged, continue waiting for the webhook
		//
		// TODO: do not wait for the webhook, we have all the information to update the pr status accordingly
		return nil
	}

	//nolint:contextcheck
	res, resp, err := userAuthClient.PullRequests.Merge(ctx, ghInstallation.Owner, ghRepo.Name, pr.GitHubPRNumber, commitMessage, mergeOpts)
	if err != nil {
		// rollback github pr state
		pr.State = previousState
		if err := svc.gitHubPullRequestRepo.Update(ctx, pr); err != nil {
			return fmt.Errorf("failed to update pull request: %w", err)
		}
		if err := svc.eventsPublisher.GitHubPRUpdated(ctx, events.Codebase(pr.CodebaseID), pr); err != nil {
			svc.logger.Error("failed to send codebase event",
				zap.String("workspace_id", pr.WorkspaceID),
				zap.String("github_pr_id", pr.ID),
				zap.Error(err),
			)
			// do not fail
		}

		var errorResponse *gh.ErrorResponse
		if resp.StatusCode == http.StatusMethodNotAllowed && errors.As(err, &errorResponse) {
			// 405 not allowed
			// This happens if the repo is configured with branch protection rules (require approvals, tests to pass, etc).
			// Proxy GitHub's error message to the end user.
			//
			// Examples:
			// * "failed to merge pr: PUT https://api.github.com/repos/zegl/empty-11/pulls/4/merge: 405 At least 1 approving review is required by reviewers with write access. []"
			return GitHubUserError{errorResponse.Message}
		}

		svc.logger.Error("unable to merge github pull request", zap.Error(err))

		return fmt.Errorf("failed to merge pr: %w", err)
	}

	if !res.GetMerged() {
		// rollback github pr state
		pr.State = previousState
		if err := svc.gitHubPullRequestRepo.Update(ctx, pr); err != nil {
			return fmt.Errorf("failed to update pull request: %w", err)
		}
		if err := svc.eventsPublisher.GitHubPRUpdated(ctx, events.Codebase(pr.CodebaseID), pr); err != nil {
			svc.logger.Error("failed to send codebase event",
				zap.String("workspace_id", pr.WorkspaceID),
				zap.String("github_pr_id", pr.ID),
				zap.Error(err),
			)
			// do not fail
		}
		return fmt.Errorf("pull request was not merged")
	}

	// This endpoint is not marking the PR as merged
	// That happens when GitHub is sending the webhook to us
	// Doing so, allows us to handle PR merges both from our and github's UI in the same way
	//
	// TODO: do not wait for the webhook, we already know everything we need
	return nil
}

var ErrIntegrationNotEnabled = errors.New("github integration is not enabled")

func (svc *Service) CreateOrUpdatePullRequest(ctx context.Context, ws *workspaces.Workspace, patchIDs []string) (*github.PullRequest, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	ghRepo, err := svc.gitHubRepositoryRepo.GetByCodebaseID(ws.CodebaseID)
	if err != nil {
		return nil, err
	}

	// Pull Requests can only be made if the integration is enabled and GitHub is considered to be the source of truth
	if !ghRepo.IntegrationEnabled || !ghRepo.GitHubSourceOfTruth {
		return nil, ErrIntegrationNotEnabled
	}

	ghInstallation, err := svc.gitHubInstallationRepo.GetByInstallationID(ghRepo.InstallationID)
	if err != nil {
		return nil, err
	}

	ghUser, err := svc.gitHubUserRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := svc.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cb, err := svc.codebaseRepo.Get(ws.CodebaseID)
	if err != nil {
		return nil, err
	}

	logger := svc.logger.With(
		zap.String("codebase_id", cb.ID),
		zap.Int64("github_installation_id", ghInstallation.InstallationID),
		zap.String("workspace_id", ws.ID),
		zap.Stringer("user_id", userID),
	)

	prs, err := svc.gitHubPullRequestRepo.ListOpenedByWorkspace(ws.ID)
	if err != nil {
		return nil, err
	}

	prBranch := "sturdy-pr-" + ws.ID
	remoteBranchName := prBranch

	// PRs that have been imported to Sturdy have user defined branch names, push update to that branch
	if len(prs) == 1 && prs[0].Head != "" {
		remoteBranchName = prs[0].Head
	}

	gitCommitMessage := message.CommitMessage(ws.DraftDescription)

	prSHA, err := svc.remoteService.PrepareBranchForPush(ctx, prBranch, ws, gitCommitMessage, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare branch: %w", err)
	}

	accessToken, err := client.GetAccessToken(ctx, logger, svc.gitHubAppConfig, ghInstallation, ghRepo.GitHubRepositoryID, svc.gitHubRepositoryRepo, svc.gitHubInstallationClientProvider)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	// GitHub Repository might have been modified at this point, refresh it
	ghRepo, err = svc.gitHubRepositoryRepo.GetByID(ghRepo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to re-load ghRepo: %w", err)
	}

	// GitHub automatically promotes the first branch pushed to a repository to be the HEAD branch
	// If the repository is _empty_ there is a risk that the branch pushed for the PR is the first branch pushed to GH
	// If this is the case, first push the sturdytrunk to be the new "master"/"main".
	// This is done _without_ force, to not screw anything up if we're in the wrong.
	if err := vcs.HaveTrackedBranch(svc.executorProvider, ws.CodebaseID, ghRepo.TrackedBranch); err != nil {
		logger.Info("pushing sturdytrunk to github")
		userVisibleError, pushTrunkErr := vcs.PushBranchToGithubSafely(logger, svc.executorProvider, ws.CodebaseID, "sturdytrunk", ghRepo.TrackedBranch, accessToken)
		if pushTrunkErr != nil {
			logger.Error("failed to push trunk to github (github is source of truth)", zap.Error(pushTrunkErr))

			// save that the push failed
			ghRepo.LastPushAt = &t
			ghRepo.LastPushErrorMessage = &userVisibleError
			if err := svc.gitHubRepositoryRepo.Update(ghRepo); err != nil {
				logger.Error("failed to update status of github integration", zap.Error(err))
			}

			return nil, gqlerrors.Error(pushTrunkErr, "pushFailure", userVisibleError)
		}
	} else {
		logger.Info("github have a default branch, not pushing sturdytrunk")
	}

	userVisibleError, pushErr := vcs.PushBranchToGithubWithForce(logger, svc.executorProvider, ws.CodebaseID, prBranch, remoteBranchName, ghUser.AccessToken)
	if pushErr != nil {
		logger.Error("failed to push to github (github is source of truth)", zap.Error(pushErr))

		// save that the push failed
		ghRepo.LastPushAt = &t
		ghRepo.LastPushErrorMessage = &userVisibleError
		if err := svc.gitHubRepositoryRepo.Update(ghRepo); err != nil {
			logger.Error("failed to update status of github integration", zap.Error(err))
		}

		return nil, gqlerrors.Error(pushErr, "pushFailure", userVisibleError)
	}

	// Mark as successfully pushed
	ghRepo.LastPushAt = &t
	ghRepo.LastPushErrorMessage = nil
	if err := svc.gitHubRepositoryRepo.Update(ghRepo); err != nil {
		logger.Error("failed to update status of github integration", zap.Error(err))
	}

	pullRequestDescription := prDescription(user.Name, ghUser.Username, cb, ws)

	// GitHub Client to be used on behalf of this user
	// TODO: Fallback to make these requests with the tokenClient if the users Access Token is invalid? (or they don't have one?)
	personalClient, err := svc.gitHubPersonalClientProvider(ghUser.AccessToken)
	if err != nil {
		return nil, err
	}

	// GitHub client to be used on behalf of the app
	tokenClient, _, err := svc.gitHubInstallationClientProvider(
		svc.gitHubAppConfig,
		ghRepo.InstallationID,
	)
	if err != nil {
		return nil, err
	}

	pullRequestTitle := ws.NameOrFallback()

	if len(prs) == 0 {
		// Create Pull Request using the personal client
		apiPR, _, err := personalClient.PullRequests.Create(ctx, ghInstallation.Owner, ghRepo.Name, &gh.NewPullRequest{
			Title: &pullRequestTitle,
			Head:  &prBranch,
			Base:  &ghRepo.TrackedBranch,
			Body:  pullRequestDescription,
		})
		if err != nil {
			return nil, gqlerrors.Error(err, "createPullRequestFailure", "Failed to create a GitHub Pull Request")
		}
		pr := github.PullRequest{
			ID:                 uuid.NewString(),
			WorkspaceID:        ws.ID,
			GitHubID:           apiPR.GetID(),
			GitHubRepositoryID: ghRepo.GitHubRepositoryID,
			CreatedBy:          userID,
			GitHubPRNumber:     apiPR.GetNumber(),
			Head:               prBranch,
			HeadSHA:            &prSHA,
			CodebaseID:         ghRepo.CodebaseID,
			Base:               ghRepo.TrackedBranch,
			State:              github.PullRequestStateOpen,
			CreatedAt:          time.Now(),
		}
		if err := svc.gitHubPullRequestRepo.Create(pr); err != nil {
			return nil, err
		}

		svc.analyticsService.Capture(ctx, "created pull request",
			analytics.CodebaseID(ws.CodebaseID),
			analytics.Property("github", true),
		)

		return &pr, nil
	}
	if len(prs) > 1 {
		logger.Error("more than one opened pull requests for a workspace - this is an erroneous state", zap.Error(err))
	}

	currentPR := prs[0]
	apiPR, _, err := tokenClient.PullRequests.Get(ctx, ghInstallation.Owner, ghRepo.Name, currentPR.GitHubPRNumber)
	if err != nil {
		return nil, gqlerrors.Error(err, "getPullRequestFailure", fmt.Sprintf("Failed to get Pull Request #%d from GitHub", currentPR.GitHubPRNumber))
	}
	apiPR.Title = &pullRequestTitle
	apiPR.Body = pullRequestDescription
	// Update the Pull Request on behalf of the user
	_, _, err = personalClient.PullRequests.Edit(ctx, ghInstallation.Owner, ghRepo.Name, currentPR.GitHubPRNumber, apiPR)
	if err != nil {
		return nil, gqlerrors.Error(err, "updatePullRequestFailure", fmt.Sprintf("Failed to update Pull Request #%d on GitHub", currentPR.GitHubPRNumber))
	}

	t0 := time.Now()
	currentPR.UpdatedAt = &t0
	currentPR.HeadSHA = &prSHA
	// Only updated_at time saved?
	if err := svc.gitHubPullRequestRepo.Update(ctx, currentPR); err != nil {
		return nil, err
	}
	svc.analyticsService.Capture(ctx, "updated pull request",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("github", true),
	)

	return currentPR, nil
}

// GitHub support (some) HTML the Pull Request descriptions, so we don't need to clean that up here.
func prDescription(userName, userGitHubLogin string, cb *codebases.Codebase, ws *workspaces.Workspace) *string {
	var builder strings.Builder
	builder.WriteString(ws.DraftDescription)
	builder.WriteString("\n\n---\n\n")

	workspaceUrl := fmt.Sprintf("https://getsturdy.com/%s/%s", cb.GenerateSlug(), ws.ID)
	builder.WriteString(fmt.Sprintf("This PR was created by %s (%s) on [Sturdy](%s).\n\n", userName, userGitHubLogin, workspaceUrl))
	builder.WriteString("Update this PR by making changes through Sturdy.\n")

	res := builder.String()
	return &res
}

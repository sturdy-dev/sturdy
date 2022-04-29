package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/changes/message"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/api"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	github_vcs "getsturdy.com/api/pkg/github/enterprise/vcs"
	"getsturdy.com/api/pkg/snapshots"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/users"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (svc *Service) ImportOpenPullRequestsByUser(ctx context.Context, codebaseID codebases.ID, userID users.ID) error {
	repo, err := svc.gitHubRepositoryRepo.GetByCodebaseID(codebaseID)
	if err != nil {
		return fmt.Errorf("failed to get github repo: %w", err)
	}

	installation, err := svc.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return fmt.Errorf("failed to get github installation: %w", err)
	}

	gitHubClient, _, err := svc.gitHubInstallationClientProvider(svc.gitHubAppConfig, installation.InstallationID)
	if err != nil {
		return fmt.Errorf("failed to get github api client: %w", err)
	}

	ghUser, err := svc.gitHubUserRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get github user: %w", err)
	}

	accessToken, err := github_client.GetAccessToken(ctx, svc.logger, svc.gitHubAppConfig, installation, repo.GitHubRepositoryID, svc.gitHubRepositoryRepo, svc.gitHubInstallationClientProvider)
	if err != nil {
		return fmt.Errorf("failed to get github access token: %w", err)
	}

	pullRequests, _, err := gitHubClient.PullRequests.List(ctx, installation.Owner, repo.Name, &gh.PullRequestListOptions{
		State: "open",
		ListOptions: gh.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get github pull requests: %w", err)
	}

	for _, pr := range pullRequests {
		// only import PRs by the requested user
		if pr.GetUser().GetLogin() != ghUser.Username {
			continue
		}

		svc.logger.Info("importing pull request", zap.Stringer("codebase_id", codebaseID), zap.Int("pr_number", pr.GetNumber()))

		err := svc.ImportPullRequest(userID, api.ConvertPullRequest(pr), repo, installation, accessToken)
		switch {
		case errors.Is(err, ErrAlreadyImported):
			continue
		case err != nil:
			return fmt.Errorf("failed import pull request: %w", err)
		}
	}

	return nil
}

var ErrAlreadyImported = errors.New("pull request has already been imported")

func (svc *Service) UpdatePullRequest(ctx context.Context, pullRequest *github.PullRequest, apiPullRequest *api.PullRequest, accessToken string, workspace *workspaces.Workspace) error {
	if !pullRequest.Importing {
		return nil
	}

	err := svc.fetchAndSnapshotPullRequest(
		apiPullRequest,
		accessToken,
		workspace,
	)
	if err != nil {
		return fmt.Errorf("failed to update workspace snapshot from pull reqeust: %w", err)
	}

	return nil
}

type PullRequestTitleDescriptioner interface {
	GetNumber() int
	GetTitle() string
	GetBody() string
}

func DescriptionFromPullRequest(pr PullRequestTitleDescriptioner) (string, error) {
	pullRequestName := pr.GetTitle()
	if pullRequestName == "" {
		pullRequestName = fmt.Sprintf("PR %d", pr.GetNumber())
	}
	pullRequestDescription, err := message.MarkdownToHtml(pr.GetBody())
	if err != nil {
		return "", fmt.Errorf("failed to render github body: %w", err)
	}
	return "<p>" + pullRequestName + "</p>" + pullRequestDescription, nil
}

func (svc *Service) fetchAndSnapshotPullRequest(
	gitHubPR *api.PullRequest,
	accessToken string,
	workspace *workspaces.Workspace,
) error {
	importBranchName := fmt.Sprintf("import-pull-request-%d-%s", gitHubPR.GetNumber(), uuid.NewString())
	refspec := fmt.Sprintf("+refs/pull/%d/head:refs/heads/%s", gitHubPR.GetNumber(), importBranchName)

	var commonAncestor string

	// Fetch to trunk
	if err := svc.executorProvider.New().
		GitWrite(github_vcs.FetchBranchWithRefspec(accessToken, refspec)).
		GitWrite(func(repo vcs.RepoGitWriter) error {
			if err := repo.CreateNewBranchAt(importBranchName, gitHubPR.GetHead().GetSHA()); err != nil {
				return fmt.Errorf("failed to create import branch")
			}

			head, err := repo.HeadCommit()
			if err != nil {
				return fmt.Errorf("could not get head: %w", err)
			}

			commonAncestor, err = repo.CommonAncestor(head.Id().String(), gitHubPR.GetHead().GetSHA())
			if err != nil {
				return fmt.Errorf("could not find common ancestor: %w", err)
			}

			svc.logger.Info("importing pull request",
				zap.String("commonAncestor", commonAncestor),
				zap.String("sturdyTrunkHead", head.Id().String()),
				zap.String("pullRequestHead", gitHubPR.GetHead().GetSHA()),
				zap.String("workspaceId", workspace.ID),
				zap.Stringer("codebaseId", workspace.CodebaseID),
			)

			// Create the workspace branch
			if err := repo.CreateNewBranchAt(workspace.ID, commonAncestor); err != nil {
				return fmt.Errorf("failed to create workspace branch")
			}

			return nil
		}).ExecTrunk(workspace.CodebaseID, "gitHubImportBranchFetch"); err != nil {
		return fmt.Errorf("failed to fetch pull to trunk: %w", err)
	}

	// make a snapshot
	//
	// step1:
	//   reset to the head of the branch. that will make all changes from pr NOT staged (not index)
	//   just what a usual sturdy user would have
	// step2:
	//   make a snapshot
	if err := svc.executorProvider.New().
		Write(vcs_view.CheckoutBranch(importBranchName)).
		Write(func(repo vcs.RepoWriter) error {
			if err := repo.ResetMixed(commonAncestor); err != nil {
				return fmt.Errorf("failed to reset temporary view to common ancestor: %w", err)
			}

			if _, err := svc.snap.Snapshot(workspace.CodebaseID, workspace.ID,
				snapshots.ActionImported,
				service_snapshots.WithMarkAsLatestInWorkspace(),
				service_snapshots.WithOnView(*repo.ViewID()),
				service_snapshots.WithOnRepo(repo), // Re-use repo context
			); err != nil {
				return fmt.Errorf("failed to create snapshot: %w", err)
			}

			return nil
		}).ExecTemporaryView(workspace.CodebaseID, "gitHubImportBranch"); err != nil {
		return fmt.Errorf("failed to create workspace from pr: %w", err)
	}

	if err := svc.executorProvider.New().Write(func(repo vcs.RepoWriter) error {
		if err := repo.DeleteBranch(importBranchName); err != nil {
			return fmt.Errorf("failed to delete importBranchName: %w", err)
		}
		return nil
	}).ExecTrunk(workspace.CodebaseID, "gitHubImportBranchCleanup"); err != nil {
		return fmt.Errorf("failed to cleanup import branch: %w", err)
	}

	return nil
}

func (svc *Service) ImportPullRequest(
	userID users.ID,
	gitHubPR *api.PullRequest,
	ghRepo *github.Repository,
	ghInstallation *github.Installation,
	accessToken string,
) error {
	codebaseID := ghRepo.CodebaseID
	// check that this pull request hasn't been imported before
	if _, err := svc.gitHubPullRequestRepo.GetByGitHubIDAndCodebaseID(gitHubPR.GetID(), codebaseID); err == nil {
		return ErrAlreadyImported
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check if pr is imported: %w", err)
	}

	draftDescription, err := DescriptionFromPullRequest(gitHubPR)
	if err != nil {
		return fmt.Errorf("failed to create description: %w", err)
	}

	workspaceID := uuid.NewString()

	t := time.Now()
	// Create the workspace
	ws := workspaces.Workspace{
		ID:               workspaceID,
		CodebaseID:       codebaseID,
		UserID:           userID,
		DraftDescription: draftDescription,
		CreatedAt:        &t,
	}
	if err := svc.workspaceWriter.Create(ws); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// fetch pr branch and create a snapshot
	if err := svc.fetchAndSnapshotPullRequest(gitHubPR, accessToken, &ws); err != nil {

		// import failed, archive the workspace that we created
		if err := svc.workspaceWriter.UpdateFields(context.Background(), ws.ID, db_workspaces.SetArchivedAt(&t)); err != nil {
			return fmt.Errorf("failed to archive workspace after failed import: %w", err)
		}

		return fmt.Errorf("failed to import pull request: %w", err)
	}

	// GitHub apps can only push to the repository where the app is installed, and not to it's forks.
	// If the PR is created from a branch in the "head" repo, Sturdy will _update_ the existing PR.
	// If the PR is created from a fork, Sturdy will create a new PR instead.

	isFork := gitHubPR.GetHead().GetUser().GetLogin() != ghInstallation.Owner

	// Create pull request object, to enable updates to existing PRs
	sturdyPR := github.PullRequest{
		ID:                 uuid.NewString(),
		WorkspaceID:        workspaceID,
		GitHubID:           gitHubPR.GetID(),
		GitHubRepositoryID: ghRepo.GitHubRepositoryID,
		CreatedBy:          userID,
		GitHubPRNumber:     gitHubPR.GetNumber(),
		Head:               gitHubPR.GetHead().GetRef(),
		HeadSHA:            gitHubPR.GetHead().SHA,
		CodebaseID:         codebaseID,
		Base:               gitHubPR.GetBase().GetRef(),
		State:              github.PullRequestStateOpen,
		CreatedAt:          gitHubPR.GetCreatedAt(),
		UpdatedAt:          nil,
		ClosedAt:           nil,
		MergedAt:           nil,
		Importing:          true,
		Fork:               isFork,
	}

	if err := svc.gitHubPullRequestRepo.Create(sturdyPR); err != nil {

		// import failed, archive the workspace that we created
		if err := svc.workspacesService.Archive(context.Background(), &ws); err != nil {
			return fmt.Errorf("failed to archive workspace after failed import: %w", err)
		}

		return fmt.Errorf("failed to save pull request record: %w", err)
	}

	svc.logger.Info("imported pull request",
		zap.String("id", sturdyPR.ID),
		zap.String("workspace_id", sturdyPR.WorkspaceID),
		zap.Bool("fork", sturdyPR.Fork))

	return nil
}

func (svc *Service) EnqueueGitHubPullRequestImport(ctx context.Context, codebaseID codebases.ID, userID users.ID) error {
	if err := (*svc.gitHubPullRequestImporterQueue).Enqueue(ctx, codebaseID, userID); err != nil {
		return err
	}
	return nil
}

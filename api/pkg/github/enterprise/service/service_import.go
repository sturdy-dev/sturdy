package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/changes/message"
	"getsturdy.com/api/pkg/github"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	github_vcs "getsturdy.com/api/pkg/github/enterprise/vcs"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/users"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (svc *Service) ImportOpenPullRequestsByUser(ctx context.Context, codebaseID string, userID users.ID) error {
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

		svc.logger.Info("importing pull request", zap.String("codebase_id", codebaseID), zap.Int("pr_number", pr.GetNumber()))

		err := svc.importPullRequest(codebaseID, userID, pr, repo, installation, accessToken)
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

func (svc *Service) importPullRequest(codebaseID string, userID users.ID, gitHubPR *gh.PullRequest, ghRepo *github.Repository, ghInstallation *github.Installation, accessToken string) error {
	// check that this pull request hasn't been imported before
	if _, err := svc.gitHubPullRequestRepo.GetByGitHubIDAndCodebaseID(gitHubPR.GetID(), codebaseID); err == nil {
		return ErrAlreadyImported
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check if pr is imported: %w", err)
	}

	workspaceID := uuid.NewString()

	pullRequestName := gitHubPR.GetTitle()
	if pullRequestName == "" {
		pullRequestName = fmt.Sprintf("PR %d", gitHubPR.GetNumber())
	}

	pullRequestDescription, err := message.MarkdownToHtml(gitHubPR.GetBody())
	if err != nil {
		return fmt.Errorf("failed to render github body: %w", err)
	}

	importBranchName := fmt.Sprintf("import-pull-request-%d-%s", gitHubPR.GetNumber(), uuid.NewString())
	importedTemporaryBranchName := fmt.Sprintf("imported-tmp-pull-request-%d-%s", gitHubPR.GetNumber(), uuid.NewString())
	refspec := fmt.Sprintf("+refs/pull/%d/head:refs/heads/%s", gitHubPR.GetNumber(), importBranchName)

	// this is done because the repository does not have a "sturdytrunk" yet
	var trunkHeadCommitID string

	// Fetch to trunk
	if err := svc.executorProvider.New().
		GitWrite(github_vcs.FetchBranchWithRefspec(accessToken, refspec)).
		GitWrite(func(repo vcs.RepoGitWriter) error {
			// Create the workspace branch
			if err := repo.CreateNewBranchOnHEAD(workspaceID); err != nil {
				return fmt.Errorf("failed to create workspace branch")
			}

			// Get trunk head
			headCommit, err := repo.HeadCommit()
			if err != nil {
				return fmt.Errorf("failed to get head: %w", err)
			}

			trunkHeadCommitID = headCommit.Id().String()

			return nil
		}).ExecTrunk(codebaseID, "gitHubImportBranchFetch"); err != nil {
		return fmt.Errorf("failed to fetch pull to trunk: %w", err)
	}

	if err := svc.executorProvider.New().
		Write(vcs_view.CheckoutBranch(importBranchName)).
		Write(func(repo vcs.RepoWriter) error {
			if err := repo.ResetMixed(trunkHeadCommitID); err != nil {
				return fmt.Errorf("failed to reset temporary view to trunk: %w", err)
			}
			if _, err := repo.AddAndCommit(fmt.Sprintf("Import from GitHub Pull Request %d", gitHubPR.GetNumber())); err != nil {
				return fmt.Errorf("failed to commit for snapshot: %w", err)
			}
			if err := repo.CreateNewBranchOnHEAD(importedTemporaryBranchName); err != nil {
				return fmt.Errorf("failed to create branch for snapshot: %w", err)
			}
			if err := repo.Push(svc.logger, importedTemporaryBranchName); err != nil {
				return fmt.Errorf("failed to push branch for snapshot: %w", err)
			}
			return nil
		}).ExecTemporaryView(codebaseID, "gitHubImportBranch"); err != nil {
		return fmt.Errorf("failed to create workspace from pr: %w", err)
	}

	t := time.Now()
	// Create the workspace
	ws := workspaces.Workspace{
		ID:               workspaceID,
		CodebaseID:       codebaseID,
		UserID:           userID,
		Name:             &pullRequestName,
		DraftDescription: pullRequestDescription,
		CreatedAt:        &t,
	}
	err = svc.workspaceWriter.Create(ws)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Create a snapshot
	makeSnapshotExec := svc.executorProvider.New().FileReadGitWrite(func(repo vcs.RepoReaderGitWriter) error {
		wsHead, err := repo.BranchCommitID(importedTemporaryBranchName)
		if err != nil {
			return fmt.Errorf("failed to get head of imported branch: %w", err)
		}

		svc.logger.Info("HEAD IS", zap.String("wsHead", wsHead), zap.String("branch", importedTemporaryBranchName))

		// Create a snapshot
		if _, err := svc.snap.Snapshot(codebaseID, workspaceID,
			snapshots.ActionSyncCompleted,
			snapshotter.WithOnTemporaryView(),
			snapshotter.WithMarkAsLatestInWorkspace(),
			snapshotter.WithOnExistingCommit(wsHead),
			snapshotter.WithOnRepo(repo), // Re-use repo context
		); err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}

		if err := repo.DeleteBranch(importBranchName); err != nil {
			return fmt.Errorf("failed to delete importBranchName: %w", err)
		}

		if err := repo.DeleteBranch(importedTemporaryBranchName); err != nil {
			return fmt.Errorf("failed to delete importedTemporaryBranchName: %w", err)
		}

		return nil
	})
	if err := makeSnapshotExec.ExecTrunk(codebaseID, "gitHubImportBranchFetch"); err != nil {
		return fmt.Errorf("failed to fetch pull to trunk: %w", err)
	}

	// GitHub apps can only push to the repository where the app is installed, and not to it's forks.
	// If the PR is created from a branch in the "head" repo, Sturdy will _update_ the existing PR.
	// If the PR is created from a fork, Sturdy will create a new PR instead.

	if gitHubPR.GetHead().GetUser().GetLogin() == ghInstallation.Owner {
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
			Open:               true,
			Merged:             false,
			CreatedAt:          gitHubPR.GetCreatedAt(),
			UpdatedAt:          nil,
			ClosedAt:           nil,
			MergedAt:           nil,
		}

		if err := svc.gitHubPullRequestRepo.Create(sturdyPR); err != nil {
			return fmt.Errorf("failed to save pull request record: %w", err)
		}
	}

	return nil
}

func (svc *Service) EnqueueGitHubPullRequestImport(ctx context.Context, codebaseID string, userID users.ID) error {
	if err := (*svc.gitHubPullRequestImporterQueue).Enqueue(ctx, codebaseID, userID); err != nil {
		return err
	}
	return nil
}

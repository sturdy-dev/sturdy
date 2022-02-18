package enterprise

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/auth"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	github_vcs "getsturdy.com/api/pkg/github/enterprise/vcs"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/workspace"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

func (r *codebaseGitHubIntegrationRootResolver) CreateWorkspaceFromGitHubBranch(ctx context.Context, args resolvers.CreateWorkspaceFromGitHubBranchArgs) (resolvers.WorkspaceResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	codebaseID := string(args.Input.CodebaseID)
	viewID := "github-import-" + uuid.NewString()
	workspaceID := uuid.NewString()

	repo, err := r.gitHubRepositoryRepo.GetByCodebaseID(codebaseID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, repo); err != nil {
		return nil, gqlerrors.Error(err)
	}

	installation, err := r.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	accessToken, err := github_client.GetAccessToken(ctx, r.logger, r.gitHubAppConfig, installation, repo.GitHubRepositoryID, r.gitHubRepositoryRepo, r.gitHubClientProvider)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/import-branch-%s", args.Input.BranchName, args.Input.BranchName)
	if err := r.gitExecutorProvider.New().
		GitWrite(github_vcs.FetchBranchWithRefspec(accessToken, refspec)).
		ExecTrunk(codebaseID, "fetchGithubBranch"); err != nil {
		return nil, gqlerrors.Error(err)
	}
	if err := r.gitExecutorProvider.
		New().
		AllowRebasingState(). // allowed because the view does not exist yet
		Schedule(func(repoProvider provider.RepoProvider) error {
			clone, err := vcs.CloneRepo(
				repoProvider.TrunkPath(codebaseID),
				repoProvider.ViewPath(codebaseID, viewID),
			)
			if err != nil {
				return fmt.Errorf("failed to clone branch")
			}

			r.logger.Info("fetched branch")

			importBranchName := fmt.Sprintf("import-branch-%s", args.Input.BranchName)
			if err = clone.FetchBranch(importBranchName); err != nil {
				return fmt.Errorf("failed to fetch branch: %w", err)
			}

			if err := clone.CreateBranchTrackingUpstream(importBranchName); err != nil {
				return fmt.Errorf("failed to create branch tracking upstream: %w", err)
			}

			if err := clone.CheckoutBranchWithForce(importBranchName); err != nil {
				return fmt.Errorf("failed to checkout branch: %w", err)
			}

			trunkCommit, err := clone.BranchCommitID("sturdytrunk")
			if err != nil {
				return fmt.Errorf("failed to get trunk commit: %w", err)
			}

			if err := clone.ResetMixed(trunkCommit); err != nil {
				return fmt.Errorf("failed to reset mixed: %w", err)
			}

			if err := clone.CreateNewBranchOnHEAD(workspaceID); err != nil {
				return fmt.Errorf("failed to create new branch: %w", err)
			}

			if err := clone.Push(r.logger, workspaceID); err != nil {
				return fmt.Errorf("failed to push: %w", err)
			}

			t := time.Now()
			// Create the workspace
			ws := workspace.Workspace{
				ID:         workspaceID,
				CodebaseID: codebaseID,
				UserID:     userID,
				Name:       &args.Input.BranchName,
				CreatedAt:  &t,
			}

			if err := r.workspaceWriter.Create(ws); err != nil {
				return fmt.Errorf("failed to create workspace: %w", err)
			}

			// Create a snapshot
			snapshot, err := r.snapshotter.Snapshot(codebaseID, workspaceID, snapshots.ActionSyncCompleted, snapshotter.WithOnView(viewID))
			if err != nil {
				return fmt.Errorf("failed to create snapshot: %w", err)
			}

			// mark as latest snapshot (will not be done automatically by the snapshotter in this scenario)
			ws.LatestSnapshotID = &snapshot.ID
			if err = r.workspaceWriter.Update(ctx, &ws); err != nil {
				return fmt.Errorf("failed to update workspace: %w", err)
			}
			return nil
		}).ExecView(codebaseID, viewID, "createWorkspaceFromGitHubBranch"); err != nil {
		return nil, gqlerrors.Error(err)
	}

	r.logger.Info("successfully imported workspace from github branch", zap.String("workspace_id", workspaceID), zap.String("branch_name", args.Input.BranchName))

	return (*r.workspaceRootResolver).Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(workspaceID)})
}

func (r *codebaseGitHubIntegrationRootResolver) ImportGitHubPullRequests(ctx context.Context, args resolvers.ImportGitHubPullRequestsInputArgs) (resolvers.CodebaseResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	codebaseID := string(args.Input.CodebaseID)

	cb, err := r.codebaseService.GetByID(ctx, codebaseID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.gitHubService.ImportOpenPullRequestsByUser(ctx, codebaseID, userID); err != nil {
		return nil, gqlerrors.Error(err)
	}

	id := graphql.ID(codebaseID)
	return (*r.codebaseRootResolver).Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
}

func (r *codebaseGitHubIntegrationRootResolver) RefreshGitHubCodebases(ctx context.Context) ([]resolvers.CodebaseResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	// add user to codebases
	if err := r.gitHubService.AddUserToCodebases(ctx, userID); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// find repos that are installed but not cloned, and add them to the clone queue
	if err := r.gitHubService.CloneMissingRepositories(ctx, userID); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return (*r.codebaseRootResolver).Codebases(ctx)
}

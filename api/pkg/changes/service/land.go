package service

import (
	"fmt"

	"github.com/google/uuid"
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"

	vcs_changes "getsturdy.com/api/pkg/changes/vcs"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/vcs"
)

func (s *Service) CreateAndLandFromView(
	viewRepo vcs.RepoWriter,
	codebaseID codebases.ID,
	workspaceID string,
	message string,
	signature git.Signature,
	diffOpts ...vcs.DiffOption,
) (commitID string, pushFunc func(vcs.RepoGitWriter) error, retErr error) {
	viewID := viewRepo.ViewID()
	if viewID == nil {
		return "", nil, fmt.Errorf("can not create on a non view")
	}

	snapshot, err := s.snap.Snapshot(codebaseID, workspaceID, snapshots.ActionPreChangeLand, service_snapshots.WithOnRepo(viewRepo), service_snapshots.WithOnView(*viewID))
	if err != nil {
		return "", nil, fmt.Errorf("failed to snapshot: %w", err)
	}

	defer func() {
		if retErr == nil {
			return
		}

		// if something went wrong, restore the view to how it was
		s.logger.Error("create and land failed, trying to restore", zap.Error(retErr))

		if err := viewRepo.CheckoutBranchWithForce(workspaceID); err != nil {
			s.logger.Error("failed to checkout workspace branch", zap.Error(err))
		}

		if err := s.snap.Restore(snapshot, viewRepo); err != nil {
			s.logger.Error("failed to restore from snapshot", zap.Error(err))
			return
		}

		s.logger.Info("successfully restored view after failed landing")
	}()

	createdCommitID, err := vcs_changes.CreateChangeFromPatchesOnRepo(s.logger, viewRepo, codebaseID, nil, message, signature, diffOpts...)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create the new change: %w", err)
	}

	if err = fastLand(viewRepo, createdCommitID); err != nil {
		return "", nil, fmt.Errorf("landing failed: %w", err)
	}

	// move the workspace branch to be the same as the new sturdytrunk
	if err := viewRepo.MoveBranch(workspaceID, "sturdytrunk"); err != nil {
		return "", nil, fmt.Errorf("failed to move workspace to new trunk: %w", err)
	}

	if err := viewRepo.CheckoutBranchWithForce(workspaceID); err != nil {
		return "", nil, fmt.Errorf("failed to checkout workspace branch: %w", err)
	}

	// LFS Pull
	if err := viewRepo.LargeFilesPull(); err != nil {
		// Log and continue (repo can have LFS files from outside of Sturdy)
		s.logger.Warn("failed to pull large files", zap.Error(err))
	}

	newBranchCommit, err := viewRepo.BranchCommitID(workspaceID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get commit after land: %w", err)
	}

	// will be executed once the new state has been recorded in the databases
	resPushFunc := func(viewRepo vcs.RepoGitWriter) error {
		if err := viewRepo.Push(s.logger, "sturdytrunk"); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}

		if err := viewRepo.ForcePush(s.logger, workspaceID); err != nil {
			return fmt.Errorf("failed to push updated workspace branch: %w", err)
		}

		return nil
	}

	return newBranchCommit, resPushFunc, nil
}

func fastLand(viewRepo vcs.RepoWriter, commitID string) (err error) {
	if err = viewRepo.FetchBranch("sturdytrunk"); err != nil {
		return fmt.Errorf("failed to fetch before fastland: %w", err)
	}

	if err := syncSingleCommitOnBranch(viewRepo, commitID, "origin", "sturdytrunk"); err != nil {
		return fmt.Errorf("failed to land: %w", err)
	}

	return nil
}

func syncSingleCommitOnBranch(repo vcs.RepoWriter, syncCommit, remoteName, branchName string) error {
	err := repo.FetchBranch(branchName)
	if err != nil {
		return fmt.Errorf("fetch origin failed: %w", err)
	}

	onto, err := repo.RemoteBranchCommit(remoteName, branchName)
	if err != nil {
		return err
	}
	defer onto.Free()

	syncingBranchName := fmt.Sprintf("syncing-%s", uuid.NewString())

	if err := repo.CreateAndCheckoutBranchAtCommit(onto.Id().String(), syncingBranchName); err != nil {
		return fmt.Errorf("create and checkout branch failed: %w", err)
	}

	_, conflicted, _, err := repo.CherryPickOnto(syncCommit, onto.Id().String())
	if conflicted {
		return fmt.Errorf("could not sync, had conflicts")
	}
	if err != nil {
		return fmt.Errorf("cherry pick failed: %w", err)
	}

	if err := repo.MoveBranchToHEAD(branchName); err != nil {
		return err
	}

	return nil
}

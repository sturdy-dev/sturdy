package service

import (
	"fmt"
	change_vcs "mash/pkg/change/vcs"
	"mash/pkg/snapshots"
	"mash/pkg/snapshots/snapshotter"
	"mash/pkg/sync"
	"mash/pkg/sync/vcs"
	"mash/pkg/unidiff"
	db_view "mash/pkg/view/db"
	db_workspace "mash/pkg/workspace/db"
	ws_meta "mash/pkg/workspace/meta"
	vcsvcs "mash/vcs"
	"mash/vcs/executor"
	"mash/vcs/provider"
	"time"

	git "github.com/libgit2/git2go/v33"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger           *zap.Logger
	executorProvider executor.Provider
	viewRepo         db_view.Repository
	workspaceReader  db_workspace.WorkspaceReader
	workspaceWriter  db_workspace.WorkspaceWriter
	snap             snapshotter.Snapshotter
}

func New(
	logger *zap.Logger,
	executorProvider executor.Provider,
	viewRepo db_view.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	workspaceWriter db_workspace.WorkspaceWriter,
	snap snapshotter.Snapshotter,
) *Service {
	return &Service{
		logger:           logger.Named("syncService"),
		executorProvider: executorProvider,
		viewRepo:         viewRepo,
		workspaceReader:  workspaceReader,
		workspaceWriter:  workspaceWriter,
		snap:             snap,
	}
}

// OnTrunk starts a sync of viewID on top of the current sturdytrunk
// If the work in progress changes on the view conflicts with trunk, a conflicting sync.RebaseStatusResponse is returned
// which has to be resolved by the user (see Resolve).
//
// The current work in progress will be added to a commit, that is rebased on top of the trunk.
// After the syncing is done, the commit is "git reset --mixed HEAD^1"-ed, to restore it to the WIP.
func (s *Service) OnTrunk(viewID string) (*sync.RebaseStatusResponse, error) {
	syncID := uuid.NewString()

	view, err := s.viewRepo.Get(viewID)
	if err != nil {
		return nil, err
	}

	branchName := "sync-" + syncID

	var rebaseStatusResponse *sync.RebaseStatusResponse

	startSyncFunc := func(repoProvider provider.RepoProvider) error {
		repo, err := repoProvider.ViewRepo(view.CodebaseID, view.ID)
		if err != nil {
			return err
		}

		if err := repo.FetchBranch("sturdytrunk"); err != nil {
			return err
		}

		trunkHeadCommit, err := repo.RemoteBranchCommit("origin", "sturdytrunk")
		if err != nil {
			return err
		}

		if err := repo.CreateNewBranchOnHEAD(branchName + "_withunsaved"); err != nil {
			return fmt.Errorf("failed to create new branch during Syncer start: %w", err)
		}

		if err := repo.CheckoutBranchSafely(branchName + "_withunsaved"); err != nil {
			return fmt.Errorf("failed to safely checkout new branch during Syncer start: %w", err)
		}

		treeID, err := change_vcs.CreateChangesTreeFromPatches(s.logger, repo, view.CodebaseID, nil)
		if err != nil {
			return fmt.Errorf("failed to create tree from patches during sync: %w", err)
		}

		// no changes, early return
		if treeID == nil {
			if err := repo.MoveBranchToCommit(branchName, trunkHeadCommit.Id().String()); err != nil {
				return fmt.Errorf("failed to move branch to commit in early return: %w", err)
			}
			if err := repo.CheckoutBranchWithForce(branchName); err != nil {
				return fmt.Errorf("failed to checkout branch in early return: %w", err)
			}
			if err := s.complete(repo, view.CodebaseID, view.WorkspaceID, view.ID, nil, nil); err != nil {
				return fmt.Errorf("failed to complete in early return: %w", err)
			}
			rebaseStatusResponse = &sync.RebaseStatusResponse{HaveConflicts: false}
			return nil
		}

		sig := git.Signature{
			Name:  "Sturdy",
			Email: "support@getsturdy.com",
			When:  time.Now(),
		}

		unsavedCommitID, err := repo.CommitIndexTree(treeID, vcs.UnsavedCommitMessage, sig)
		if err != nil {
			return fmt.Errorf("failed to create commit with unsave changes: %w", err)
		}

		err = repo.CreateAndCheckoutBranchAtCommit(trunkHeadCommit.Id().String(), branchName)
		if err != nil {
			return fmt.Errorf("create and checkout branch failed: %w", err)
		}

		// Apply our unsaved changes
		rb, rebasedCommits, err := repo.InitRebaseRaw(
			unsavedCommitID,
			trunkHeadCommit.Id().String(),
		)

		if err != nil {
			return err
		}

		rebaseStatus, err := rb.Status()
		if err != nil {
			return err
		}

		// We have conflicts, require resolution from user
		if rebaseStatus == vcsvcs.RebaseHaveConflicts {
			// Restore large files
			if err := repo.LargeFilesPull(); err != nil {
				// don't fail
				s.logger.Error("failed to restore large files", zap.Error(err))
			}

			rebaseStatusResponse, err = Status(s.logger, rb)
			if err != nil {
				return fmt.Errorf("failed to get conflict status: %w", err)
			}

			return nil
		}

		// No conflicts

		if err := repo.MoveBranchToHEAD(branchName); err != nil {
			return fmt.Errorf("branch to head failed: %w", err)
		}

		if err := s.complete(repo, view.CodebaseID, view.WorkspaceID, view.ID, &unsavedCommitID, rebasedCommits); err != nil {
			return err
		}

		rebaseStatusResponse = &sync.RebaseStatusResponse{HaveConflicts: false}
		return nil
	}

	err = s.executorProvider.New().
		AssertBranchName(view.WorkspaceID).
		AllowRebasingState().
		Schedule(startSyncFunc).
		ExecView(view.CodebaseID, view.ID, "syncOnTrunk2")
	if err != nil {
		return nil, err
	}

	if rebaseStatusResponse == nil {
		return nil, fmt.Errorf("no rebase status found")
	}

	return rebaseStatusResponse, nil
}

// complete is called by OnTrunk (if there where no conflicts) and Resolve (when all conflicts have been resolved)
func (svc *Service) complete(repo vcsvcs.RepoWriter, codebaseID, workspaceID, viewID string, unsavedCommitID *string, rebasedCommits []vcsvcs.RebasedCommit) error {
	if err := repo.MoveBranchToHEAD(workspaceID); err != nil {
		return fmt.Errorf("failed to move workspace to head: %w", err)
	}

	if err := repo.CheckoutBranchWithForce(workspaceID); err != nil {
		return fmt.Errorf("failed to checkout workspace with force: %w", err)
	}

	// Restore large files
	if err := repo.LargeFilesPull(); err != nil {
		// don't fail
		svc.logger.Error("failed to restore large files", zap.Error(err))
	}

	if unsavedCommitID != nil && shouldResetHead(*unsavedCommitID, rebasedCommits) {
		head, err := repo.HeadCommit()
		if err != nil {
			return fmt.Errorf("could not get head commit: %w", err)
		}
		parents, err := repo.GetCommitParents(head.Id().String())
		if err != nil {
			return fmt.Errorf("failed to get parents: %w", err)
		}

		if len(parents) != 1 {
			return fmt.Errorf("unexpected number of parents: %d", len(parents))
		}

		err = repo.ResetMixed(parents[0])
		if err != nil {
			return fmt.Errorf("failed to reset: %w", err)
		}
	}

	if err := repo.ForcePush(svc.logger, workspaceID); err != nil {
		return fmt.Errorf("failed to push result: %w", err)
	}

	// Make a snapshot (right away)
	// The "conflict" status is calculated based on the latest snapshot of a workspace
	// Create a snapshot right away to re-calculate the conflicting status
	if _, err := svc.snap.Snapshot(codebaseID, workspaceID, snapshots.ActionSyncCompleted,
		snapshotter.WithOnView(viewID),
		snapshotter.WithOnRepo(repo),
	); err != nil {
		svc.logger.Error("failed to snapshot", zap.Error(err))
		// Don't fail
	}

	// Update workspace
	if err := ws_meta.Updated(svc.workspaceReader, svc.workspaceWriter, workspaceID); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	/*err = postHogClient.Enqueue(posthog.Capture{
		Event:      "completed sync",
		DistinctId: syncer.UserID,
		Properties: posthog.NewProperties().
			Set("codebase_id", syncer.CodebaseID).
			Set("workspace_id", syncer.WorkspaceID),
	})
	if err != nil {
		logger.Error("posthog failed", zap.Error(err))
	}
	*/

	return nil
}

// shouldResetHead returns true if the commit with the unsaved was committed
// This is _not_ the case if the changes in WIP where identical to a commit that has been applied on the trunk.
func shouldResetHead(unsavedCommitOld string, rebasedCommits []vcsvcs.RebasedCommit) bool {
	for _, c := range rebasedCommits {
		if c.OldCommitID == unsavedCommitOld {
			// The 'unsaved changes' commit that we would want to reset is not noop
			if c.Noop {
				return false
			}
			return true
		}
	}
	return false
}

func Status(logger *zap.Logger, rebasing *vcsvcs.SturdyRebase) (*sync.RebaseStatusResponse, error) {
	conflictingFiles, err := rebasing.ConflictingFiles()
	if err != nil {
		return nil, err
	}

	var cf []sync.ConflictingFile
	for _, p := range conflictingFiles {
		patches, err := rebasing.ConflictDiff(p)
		if err != nil {
			return nil, fmt.Errorf("failed to get conflict patches for %s: %w", p, err)
		}

		workspaceDiff, err := unidiff.NewUnidiff(unidiff.NewStringsPatchReader([]string{patches.WorkspacePatch}), logger).DecorateSingle()
		if err != nil {
			return nil, fmt.Errorf("failed to decorate workspace diff: %w", err)
		}

		trunkDiff, err := unidiff.NewUnidiff(unidiff.NewStringsPatchReader([]string{patches.TrunkPatch}), logger).DecorateSingle()
		if err != nil {
			return nil, fmt.Errorf("failed to decorate workspace diff: %w", err)
		}

		cf = append(cf, sync.ConflictingFile{
			Path:          p,
			WorkspaceDiff: workspaceDiff,
			TrunkDiff:     trunkDiff,
		})
	}

	return &sync.RebaseStatusResponse{
		IsRebasing:       true,
		HaveConflicts:    true,
		ConflictingFiles: cf,
		ProgressCurrent:  1, // TODO: Remove from API
		ProgressTotal:    1, // TODO: Remove from API
	}, nil
}

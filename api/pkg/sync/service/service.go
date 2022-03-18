package service

import (
	"context"
	"fmt"
	"time"

	change_vcs "getsturdy.com/api/pkg/changes/vcs"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/sync"
	"getsturdy.com/api/pkg/sync/vcs"
	"getsturdy.com/api/pkg/unidiff"
	db_view "getsturdy.com/api/pkg/view/db"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	ws_meta "getsturdy.com/api/pkg/workspaces/meta"
	vcsvcs "getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

type Service struct {
	logger           *zap.Logger
	executorProvider executor.Provider
	viewRepo         db_view.Repository
	workspaceReader  db_workspaces.WorkspaceReader
	workspaceWriter  db_workspaces.WorkspaceWriter
	snap             snapshotter.Snapshotter

	eventsPublisher *events.Publisher
}

func New(
	logger *zap.Logger,
	executorProvider executor.Provider,
	viewRepo db_view.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	workspaceWriter db_workspaces.WorkspaceWriter,
	snap snapshotter.Snapshotter,
	eventsPublisher *events.Publisher,
) *Service {
	return &Service{
		logger:           logger.Named("syncService"),
		executorProvider: executorProvider,
		viewRepo:         viewRepo,
		workspaceReader:  workspaceReader,
		workspaceWriter:  workspaceWriter,
		snap:             snap,
		eventsPublisher:  eventsPublisher,
	}
}

// OnTrunk starts a sync of the workspace on top of the current sturdytrunk
// If the work in progress changes on the workspace conflicts with trunk, a conflicting sync.RebaseStatusResponse is returned
// which has to be resolved by the user (see Resolve).
//
// The current work in progress will be added to a commit, that is rebased on top of the trunk.
// After the syncing is done, the commit is "git reset --mixed HEAD^1"-ed, to restore it to the WIP.
func (svc *Service) OnTrunk(ctx context.Context, ws *workspaces.Workspace) (*sync.RebaseStatusResponse, error) {
	syncID := uuid.NewString()

	branchName := fmt.Sprintf("sync-%s", syncID)

	var rebaseStatusResponse *sync.RebaseStatusResponse

	rebaseFunc := func(repo vcsvcs.RepoWriter) error {
		// Already rebasing, exit
		if repo.IsRebasing() {
			rb, err := repo.OpenRebase()
			if err != nil {
				return fmt.Errorf("failed to open previous rebase: %w", err)
			}
			rebaseStatusResponse, err = Status(svc.logger, rb)
			if err != nil {
				return fmt.Errorf("failed to get conflict status: %w", err)
			}
			return nil
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

		if err := svc.logFiles(ws.ID, "before", repo); err != nil {
			return fmt.Errorf("failed to log changed files before sync: %w", err)
		}

		treeID, err := change_vcs.CreateChangesTreeFromPatches(svc.logger, repo, ws.CodebaseID, nil)
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
			if err := svc.complete(ctx, repo, ws.CodebaseID, ws.ID, *repo.ViewID(), nil, nil); err != nil {
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

		if err := repo.CreateAndCheckoutBranchAtCommit(trunkHeadCommit.Id().String(), branchName); err != nil {
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
				svc.logger.Error("failed to restore large files", zap.Error(err))
			}

			rebaseStatusResponse, err = Status(svc.logger, rb)
			if err != nil {
				return fmt.Errorf("failed to get conflict status: %w", err)
			}

			return nil
		}

		// No conflicts

		if err := repo.MoveBranchToHEAD(branchName); err != nil {
			return fmt.Errorf("branch to head failed: %w", err)
		}

		if err := svc.complete(ctx, repo, ws.CodebaseID, ws.ID, *repo.ViewID(), &unsavedCommitID, rebasedCommits); err != nil {
			return err
		}

		rebaseStatusResponse = &sync.RebaseStatusResponse{HaveConflicts: false}
		return nil
	}

	if ws.ViewID != nil {
		if err := svc.executorProvider.New().
			AssertBranchName(ws.ID).
			AllowRebasingState(). // allowed to get the state of existing conflicts
			Write(rebaseFunc).
			ExecView(ws.CodebaseID, *ws.ViewID, "syncOnTrunk"); err != nil {
			return nil, err
		}
		vw, err := svc.viewRepo.Get(*ws.ViewID)
		if err != nil {
			return nil, fmt.Errorf("failed to get view: %w", err)
		}

		if err := svc.eventsPublisher.ViewUpdated(ctx, events.Codebase(vw.CodebaseID), vw); err != nil {
			svc.logger.Error("failed to send workspace updated event", zap.Error(err))
			// do not fail
		}
	} else {
		if err := svc.executorProvider.New().
			Write(vcs_view.CheckoutBranch(ws.ID)).
			Write(rebaseFunc).
			ExecTemporaryView(ws.CodebaseID, "syncOnTrunk"); err != nil {
			return nil, err
		}
	}

	if rebaseStatusResponse == nil {
		return nil, fmt.Errorf("no rebase status found")
	}

	return rebaseStatusResponse, nil
}

// complete is called by OnTrunk (if there where no conflicts) and Resolve (when all conflicts have been resolved)
func (svc *Service) complete(ctx context.Context, repo vcsvcs.RepoWriter, codebaseID codebases.ID, workspaceID, viewID string, unsavedCommitID *string, rebasedCommits []vcsvcs.RebasedCommit) error {
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

	if err := svc.logFiles(workspaceID, "complete", repo); err != nil {
		return fmt.Errorf("failed to log changed files before sync: %w", err)
	}

	// Make a snapshot (right away)
	// The "conflict" status is calculated based on the latest snapshot of a workspace
	// Create a snapshot right away to re-calculate the conflicting status
	if _, err := svc.snap.Snapshot(codebaseID, workspaceID, snapshots.ActionSyncCompleted,
		snapshotter.WithOnView(viewID),
		snapshotter.WithOnRepo(repo),
		snapshotter.WithMarkAsLatestInWorkspace(),
	); err != nil {
		svc.logger.Error("failed to snapshot", zap.Error(err))
		// Don't fail
	}

	// Update workspace
	if err := ws_meta.Updated(ctx, svc.workspaceReader, svc.workspaceWriter, workspaceID); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

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
	}, nil
}

func (svc *Service) logFiles(workspaceID, state string, repo vcsvcs.RepoGitReader) error {
	preSyncDiff, err := repo.CurrentDiffNoIndex()
	if err != nil {
		return fmt.Errorf("could not get diffs: %w", err)
	}
	err = preSyncDiff.ForEach(func(delta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		svc.logger.Info("debug-sync-diff", zap.String("state", state), zap.String("workspace_id", workspaceID), zap.String("new_file_path", delta.NewFile.Path), zap.String("old_file_path", delta.OldFile.Path))
		return nil, nil
	}, git.DiffDetailFiles)
	if err != nil {
		return fmt.Errorf("could log diffs: %w", err)
	}
	return nil
}

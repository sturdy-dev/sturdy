package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/sync"
	vcsvcs "getsturdy.com/api/vcs"
)

// Resolve resolves the conflicts in viewID with the resolutions in resolves
//
// For each conflicting file in the index, resolves contains the the file path and if the resolution should be
// * use the version from trunk
// * use the version from the workspace
// * use the current version of the file on disk (called "custom")
func (svc *Service) Resolve(ctx context.Context, viewID string, resolves []vcsvcs.SturdyRebaseResolve) (*sync.RebaseStatusResponse, error) {
	view, err := svc.viewRepo.Get(viewID)
	if err != nil {
		return nil, err
	}

	var rebaseStatusResponse *sync.RebaseStatusResponse

	resolveSyncFunc := func(repo vcsvcs.RepoWriter) error {
		rb, err := repo.OpenRebase()
		if err != nil {
			return err
		}

		if err := rb.ResolveFiles(resolves); err != nil {
			return err
		}

		conflicts, rebasedCommits, err := rb.Continue()
		if err != nil {
			return err
		}
		if conflicts {
			return fmt.Errorf("unexpected conflict after conflict resolution")
		}
		if len(rebasedCommits) != 1 {
			return fmt.Errorf("unexpected number of rebased commits")
		}

		// No conflicts

		if err := svc.complete(ctx, repo, view.CodebaseID, view.WorkspaceID, view.ID, &rebasedCommits[0].OldCommitID, rebasedCommits); err != nil {
			return err
		}

		rebaseStatusResponse = &sync.RebaseStatusResponse{HaveConflicts: false}
		return nil
	}

	err = svc.executorProvider.New().
		AllowRebasingState(). // allowed to get the state of existing conflicts
		Write(resolveSyncFunc).
		ExecView(view.CodebaseID, view.ID, "syncResolve2")
	if err != nil {
		return nil, err
	}

	if rebaseStatusResponse == nil {
		return nil, fmt.Errorf("no rebase status found")
	}

	return rebaseStatusResponse, nil
}

package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/changes/message"
	service_change "getsturdy.com/api/pkg/changes/service"
	vcs_change "getsturdy.com/api/pkg/changes/vcs"
	"getsturdy.com/api/pkg/remote"
	db_remote "getsturdy.com/api/pkg/remote/enterprise/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
)

type Service struct {
	repo             db_remote.Repository
	executorProvider executor.Provider
	logger           *zap.Logger
	workspaceReader  db_workspaces.WorkspaceReader
	snap             snapshotter.Snapshotter
	changeService    *service_change.Service
}

func New(
	repo db_remote.Repository,
	executorProvider executor.Provider,
	logger *zap.Logger,
	workspaceReader db_workspaces.WorkspaceReader,
	snap snapshotter.Snapshotter,
	changeService *service_change.Service,
) *Service {
	return &Service{
		repo:             repo,
		executorProvider: executorProvider,
		logger:           logger,
		workspaceReader:  workspaceReader,
		snap:             snap,
		changeService:    changeService,
	}
}

func (svc *Service) Get(ctx context.Context, codebaseID string) (*remote.Remote, error) {
	rep, err := svc.repo.GetByCodebaseID(ctx, codebaseID)
	if err != nil {
		return nil, err
	}
	return rep, nil
}

func (svc *Service) SetRemote(ctx context.Context, codebaseID, name, url, username, password, trackedBranch string) (*remote.Remote, error) {
	// update existing if exists
	rep, err := svc.repo.GetByCodebaseID(ctx, codebaseID)
	switch {
	case err == nil:
		// update
		rep.Name = name
		rep.URL = url
		rep.BasicAuthUsername = username
		rep.BasicAuthPassword = password
		rep.TrackedBranch = trackedBranch
		if err := svc.repo.Update(ctx, rep); err != nil {
			return nil, fmt.Errorf("failed to update remote: %w", err)
		}
		return rep, nil
	case errors.Is(err, sql.ErrNoRows):
		// create
		r := remote.Remote{
			ID:                uuid.NewString(),
			CodebaseID:        codebaseID,
			Name:              name,
			URL:               url,
			BasicAuthUsername: username,
			BasicAuthPassword: password,
			TrackedBranch:     trackedBranch,
		}
		if err := svc.repo.Create(ctx, r); err != nil {
			return nil, fmt.Errorf("failed to add remote: %w", err)
		}
		return &r, nil
	default:
		return nil, fmt.Errorf("failed to set remote: %w", err)
	}
}

func (svc *Service) Push(ctx context.Context, user *users.User, ws *workspaces.Workspace) error {
	rem, err := svc.repo.GetByCodebaseID(ctx, ws.CodebaseID)
	if err != nil {
		return fmt.Errorf("could not get remote: %w", err)
	}

	localBranchName := "sturdy-" + ws.ID
	gitCommitMessage := message.CommitMessage(ws.DraftDescription)

	_, err = svc.PrepareBranchForPush(ctx, localBranchName, ws, gitCommitMessage, user.Name, user.Email)
	if err != nil {
		return err
	}

	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/sturdy-%s", localBranchName, ws.ID)

	push := func(repo vcs.RepoGitWriter) error {
		_, err := repo.PushRemoteUrlWithRefspec(
			svc.logger,
			rem.URL,
			newCredentialsCallback(rem.BasicAuthPassword, rem.BasicAuthPassword),
			[]string{refspec},
		)
		if err != nil {
			return fmt.Errorf("failed to pull: %w", err)
		}
		return nil
	}

	if err := svc.executorProvider.New().GitWrite(push).ExecTrunk(ws.CodebaseID, "pushRemote"); err != nil {
		return fmt.Errorf("failed to push from trunk: %w", err)
	}

	return nil
}

func (svc *Service) Pull(ctx context.Context, codebaseID string) error {
	rem, err := svc.repo.GetByCodebaseID(ctx, codebaseID)
	if err != nil {
		return fmt.Errorf("could not get remote: %w", err)
	}

	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/sturdytrunk", rem.TrackedBranch)

	pull := func(repo vcs.RepoGitWriter) error {
		err := repo.FetchUrlRemoteWithCreds(
			rem.URL,
			newCredentialsCallback(rem.BasicAuthPassword, rem.BasicAuthPassword),
			[]string{refspec},
		)
		if err != nil {
			return fmt.Errorf("failed to pull: %w", err)
		}
		return nil
	}

	if err := svc.executorProvider.New().GitWrite(pull).ExecTrunk(codebaseID, "pullRemote"); err != nil {
		return fmt.Errorf("failed to pull: %w", err)
	}

	if err := svc.changeService.UnsetHeadChangeCache(codebaseID); err != nil {
		return fmt.Errorf("failed to unset head: %w", err)
	}

	return nil
}

func newCredentialsCallback(username, password string) git.CredentialsCallback {
	return func(url string, usernameFromUrl string, allowedTypes git.CredentialType) (*git.Credential, error) {
		cred, _ := git.NewCredentialUserpassPlaintext(username, password)
		return cred, nil
	}
}

func (svc *Service) PrepareBranchForPush(ctx context.Context, prBranchName string, ws *workspaces.Workspace, commitMessage, userName, userEmail string) (commitSha string, err error) {
	if ws.ViewID == nil && ws.LatestSnapshotID != nil {
		commitSha, err = svc.prepareBranchForPullRequestFromSnapshot(ctx, prBranchName, ws, commitMessage, userName, userEmail)
		if err != nil {
			return "", fmt.Errorf("failed to prepare branch from snapshot: %w", err)
		}
		return
	} else if ws.ViewID != nil {
		commitSha, err = svc.prepareBranchForPullRequestWithView(prBranchName, ws, commitMessage, userName, userEmail)
		if err != nil {
			return "", fmt.Errorf("failed to prepare branch from snapshot: %w", err)
		}
		return
	} else {
		return "", errors.New("workspace does not have either view nor snapshot")
	}

}

func (svc *Service) prepareBranchForPullRequestFromSnapshot(ctx context.Context, prBranchName string, ws *workspaces.Workspace, commitMessage, userName, userEmail string) (string, error) {
	signature := git.Signature{
		Name:  userName,
		Email: userEmail,
		When:  time.Now(),
	}

	snapshot, err := svc.snap.GetByID(ctx, *ws.LatestSnapshotID)
	if err != nil {
		return "", fmt.Errorf("failed to get snapshot: %w", err)
	}

	var resSha string

	exec := svc.executorProvider.New().GitWrite(func(r vcs.RepoGitWriter) error {
		sha, err := r.CreateNewCommitBasedOnCommit(prBranchName, snapshot.CommitID, signature, commitMessage)
		if err != nil {
			return err
		}

		resSha = sha
		return nil
	})

	if err := exec.ExecTrunk(ws.CodebaseID, "prepareBranchForPullRequestFromSnapshot"); err != nil {
		return "", fmt.Errorf("failed to create pr branch from snapshot")
	}

	return resSha, nil
}

func (svc *Service) prepareBranchForPullRequestWithView(prBranchName string, ws *workspaces.Workspace, commitMessage, userName, userEmail string) (string, error) {
	signature := git.Signature{
		Name:  userName,
		Email: userEmail,
		When:  time.Now(),
	}

	var resSha string

	exec := svc.executorProvider.New().FileReadGitWrite(func(r vcs.RepoReaderGitWriter) error {
		treeID, err := vcs_change.CreateChangesTreeFromPatches(svc.logger, r, ws.CodebaseID, nil)
		if err != nil {
			return err
		}

		// No changes where added
		if treeID == nil {
			return fmt.Errorf("no changes to add")
		}

		if err := r.CreateNewBranchOnHEAD(prBranchName); err != nil {
			return fmt.Errorf("failed to create pr branch: %w", err)
		}

		sha, err := r.CommitIndexTreeWithReference(treeID, commitMessage, signature, "refs/heads/"+prBranchName)
		if err != nil {
			return fmt.Errorf("failed save change: %w", err)
		}

		if err := r.ForcePush(svc.logger, prBranchName); err != nil {
			return fmt.Errorf("failed to push to sturdytrunk: %w", err)
		}

		resSha = sha
		return nil
	})

	if err := exec.ExecView(ws.CodebaseID, *ws.ViewID, "prepareBranchForPullRequestWithView"); err != nil {
		return "", fmt.Errorf("failed to create pr branch from view: %w", err)
	}

	return resSha, nil
}
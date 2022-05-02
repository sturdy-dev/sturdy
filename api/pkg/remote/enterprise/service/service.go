package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/uuid"
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
	ssh2 "golang.org/x/crypto/ssh"

	"getsturdy.com/api/pkg/analytics"
	analytics_service "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/changes/message"
	service_change "getsturdy.com/api/pkg/changes/service"
	vcs_change "getsturdy.com/api/pkg/changes/vcs"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/crypto"
	db_crypto "getsturdy.com/api/pkg/crypto/db"
	"getsturdy.com/api/pkg/remote"
	db_remote "getsturdy.com/api/pkg/remote/enterprise/db"
	"getsturdy.com/api/pkg/remote/service"
	service_snapshotter "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
)

type EnterpriseService struct {
	repo              db_remote.Repository
	executorProvider  executor.Provider
	logger            *zap.Logger
	workspaceReader   db_workspaces.WorkspaceReader
	workspaceWriter   db_workspaces.WorkspaceWriter
	snap              *service_snapshotter.Service
	changeService     *service_change.Service
	analyticsService  *analytics_service.Service
	keyPairRepository db_crypto.KeyPairRepository
}

var _ service.Service = (*EnterpriseService)(nil)

func New(
	repo db_remote.Repository,
	executorProvider executor.Provider,
	logger *zap.Logger,
	workspaceReader db_workspaces.WorkspaceReader,
	workspaceWriter db_workspaces.WorkspaceWriter,
	snap *service_snapshotter.Service,
	changeService *service_change.Service,
	analyticsService *analytics_service.Service,
	keyPairRepository db_crypto.KeyPairRepository,
) *EnterpriseService {
	return &EnterpriseService{
		repo:              repo,
		executorProvider:  executorProvider,
		logger:            logger,
		workspaceReader:   workspaceReader,
		workspaceWriter:   workspaceWriter,
		snap:              snap,
		changeService:     changeService,
		analyticsService:  analyticsService,
		keyPairRepository: keyPairRepository,
	}
}

func (svc *EnterpriseService) Get(ctx context.Context, codebaseID codebases.ID) (*remote.Remote, error) {
	rep, err := svc.repo.GetByCodebaseID(ctx, codebaseID)
	if err != nil {
		return nil, err
	}
	return rep, nil
}

func (svc *EnterpriseService) GetWithFixedURL(ctx context.Context, codebaseID codebases.ID) (*remote.Remote, error) {
	rep, err := svc.repo.GetByCodebaseID(ctx, codebaseID)
	if err != nil {
		return nil, err
	}

	if rep.KeyPairID != nil {
		rep.URL = rewriteSshUrl(rep.URL)
	}

	return rep, nil
}

type SetRemoteInput struct {
	Name              string
	URL               string
	TrackedBranch     string
	BasicAuthUsername *string
	BasicAuthPassword *string
	KeyPairID         *crypto.KeyPairID
	BrowserLinkRepo   string
	BrowserLinkBranch string
	Enabled           bool
}

func (svc *EnterpriseService) SetRemote(ctx context.Context, codebaseID codebases.ID, input *SetRemoteInput) (*remote.Remote, error) {
	hasBasic := input.BasicAuthUsername != nil && input.BasicAuthPassword != nil
	hasKeyPair := input.KeyPairID != nil

	if hasBasic && hasKeyPair {
		return nil, fmt.Errorf("basic auth and keypair auth are mutually exclusive")
	}
	if !hasBasic && !hasKeyPair {
		return nil, fmt.Errorf("no auth method set")
	}

	// make sure that only the relevant fields are set
	if hasBasic {
		input.KeyPairID = nil
	} else if hasKeyPair {
		input.BasicAuthUsername = nil
		input.BasicAuthPassword = nil
	} else {
		return nil, fmt.Errorf("unexpected auth configuration")
	}

	// update existing if exists
	rep, err := svc.repo.GetByCodebaseID(ctx, codebaseID)
	switch {
	case err == nil:
		// update
		rep.Name = input.Name
		rep.URL = input.URL
		rep.TrackedBranch = input.TrackedBranch
		rep.BasicAuthUsername = input.BasicAuthUsername
		rep.BasicAuthPassword = input.BasicAuthPassword
		rep.KeyPairID = input.KeyPairID
		rep.BrowserLinkRepo = input.BrowserLinkRepo
		rep.BrowserLinkBranch = input.BrowserLinkBranch
		rep.Enabled = input.Enabled
		if err := svc.repo.Update(ctx, rep); err != nil {
			return nil, fmt.Errorf("failed to update remote: %w", err)
		}

		svc.analyticsService.Capture(ctx, "updated remote integration", analytics.CodebaseID(codebaseID), analytics.Property("remote_name", rep.Name))

		return rep, nil
	case errors.Is(err, sql.ErrNoRows):
		// create
		r := remote.Remote{
			ID:                uuid.NewString(),
			CodebaseID:        codebaseID,
			Name:              input.Name,
			URL:               input.URL,
			TrackedBranch:     input.TrackedBranch,
			BasicAuthUsername: input.BasicAuthUsername,
			BasicAuthPassword: input.BasicAuthPassword,
			KeyPairID:         input.KeyPairID,
			BrowserLinkRepo:   input.BrowserLinkRepo,
			BrowserLinkBranch: input.BrowserLinkBranch,
			Enabled:           input.Enabled,
		}

		if err := svc.repo.Create(ctx, r); err != nil {
			return nil, fmt.Errorf("failed to add remote: %w", err)
		}

		svc.analyticsService.Capture(ctx, "created remote integration", analytics.CodebaseID(codebaseID), analytics.Property("remote_name", r.Name))

		return &r, nil
	default:
		return nil, fmt.Errorf("failed to set remote: %w", err)
	}
}

var ErrRemoteDisabled = errors.New("this remote is disabled")

func (svc *EnterpriseService) Push(ctx context.Context, user *users.User, ws *workspaces.Workspace) error {
	rem, err := svc.GetWithFixedURL(ctx, ws.CodebaseID)
	if err != nil {
		return fmt.Errorf("could not get remote: %w", err)
	}
	if !rem.Enabled {
		return ErrRemoteDisabled
	}

	localBranchName := "sturdy-" + ws.ID
	gitCommitMessage := message.CommitMessage(ws.DraftDescription)

	_, err = svc.PrepareBranchForPush(ctx, localBranchName, ws, gitCommitMessage, user.Name, user.Email)
	if err != nil {
		return err
	}

	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/sturdy-%s", localBranchName, ws.ID)

	creds, err := svc.newCredentialsCallback(ctx, rem)
	if err != nil {
		return fmt.Errorf("could not get creds: %w", err)
	}

	push := func(repo vcs.RepoGitWriter) error {
		_, err := repo.PushRemoteUrlWithRefspec(rem.URL, creds, []config.RefSpec{config.RefSpec(refspec)})
		switch {
		case errors.Is(err, gogit.NoErrAlreadyUpToDate):
			return nil
		case err != nil:
			return fmt.Errorf("failed to push: %w", err)
		default:
			return nil
		}
	}

	if err := svc.executorProvider.New().GitWrite(push).ExecTrunk(ws.CodebaseID, "pushRemote"); err != nil {
		return fmt.Errorf("failed to push workspace to remote: %w", err)
	}

	svc.analyticsService.CaptureUser(ctx, user.ID, "pushed workspace to remote", analytics.CodebaseID(ws.CodebaseID), analytics.Property("workspace_id", ws.ID))

	return nil
}

func (svc *EnterpriseService) PushTrunk(ctx context.Context, codebaseID codebases.ID) error {
	rem, err := svc.GetWithFixedURL(ctx, codebaseID)
	if err != nil {
		return fmt.Errorf("could not get remote: %w", err)
	}
	if !rem.Enabled {
		return ErrRemoteDisabled
	}

	refspec := fmt.Sprintf("refs/heads/sturdytrunk:refs/heads/%s", rem.TrackedBranch)

	creds, err := svc.newCredentialsCallback(ctx, rem)
	if err != nil {
		return fmt.Errorf("could not get creds: %w", err)
	}

	push := func(repo vcs.RepoGitWriter) error {
		_, err := repo.PushRemoteUrlWithRefspec(rem.URL, creds, []config.RefSpec{config.RefSpec(refspec)})
		switch {
		case errors.Is(err, gogit.NoErrAlreadyUpToDate):
			return nil
		case err != nil:
			return fmt.Errorf("failed to push: %w", err)
		default:
			return nil
		}
	}

	if err := svc.executorProvider.New().GitWrite(push).ExecTrunk(codebaseID, "pushTrunkRemote"); err != nil {
		return fmt.Errorf("failed to push trunk to remote: %w", err)
	}

	svc.analyticsService.Capture(ctx, "pushed trunk to remote", analytics.CodebaseID(codebaseID))

	return nil
}

func (svc *EnterpriseService) Pull(ctx context.Context, codebaseID codebases.ID) error {
	rem, err := svc.GetWithFixedURL(ctx, codebaseID)
	if err != nil {
		return fmt.Errorf("could not get remote: %w", err)
	}
	if !rem.Enabled {
		return ErrRemoteDisabled
	}

	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/sturdytrunk", rem.TrackedBranch)

	creds, err := svc.newCredentialsCallback(ctx, rem)
	if err != nil {
		return fmt.Errorf("could not get creds: %w", err)
	}

	pull := func(repo vcs.RepoGitWriter) error {
		err := repo.FetchUrlRemoteWithCreds(rem.URL, creds, []config.RefSpec{config.RefSpec(refspec)})
		switch {
		case errors.Is(err, gogit.NoErrAlreadyUpToDate):
			return nil
		case err != nil:
			return fmt.Errorf("failed to pull: %w", err)
		default:
			return nil
		}
	}

	if err := svc.executorProvider.New().GitWrite(pull).ExecTrunk(codebaseID, "pullRemote"); err != nil {
		return fmt.Errorf("failed to pull: %w", err)
	}

	svc.analyticsService.Capture(ctx, "pulled trunk from remote", analytics.CodebaseID(codebaseID))

	if err := svc.changeService.UnsetHeadChangeCache(codebaseID); err != nil {
		return fmt.Errorf("failed to unset head: %w", err)
	}

	// Allow all workspaces to be rebased/synced on the latest head
	if err := svc.workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(codebaseID); err != nil {
		return fmt.Errorf("failed to unset up to date with trunk for all in codebase: %w", err)
	}

	return nil
}

func (svc *EnterpriseService) newCredentialsCallback(ctx context.Context, rem *remote.Remote) (cb transport.AuthMethod, err error) {
	if rem.KeyPairID != nil {
		kp, kpErr := svc.keyPairRepository.Get(ctx, *rem.KeyPairID)
		if kpErr != nil {
			return nil, fmt.Errorf("could not get kp: %w", kpErr)
		}

		am, err := ssh.NewPublicKeys("git", []byte(kp.PrivateKey), "")
		if err != nil {
			return nil, err
		}

		am.HostKeyCallback = ssh2.InsecureIgnoreHostKey()

		return am, nil
	}

	if rem.BasicAuthUsername != nil && rem.BasicAuthPassword != nil {
		return &http.BasicAuth{
			Username: *rem.BasicAuthUsername,
			Password: *rem.BasicAuthPassword,
		}, nil
	}

	return nil, errors.New("no auth method found")
}

func (svc *EnterpriseService) PrepareBranchForPush(ctx context.Context, prBranchName string, ws *workspaces.Workspace, commitMessage, userName, userEmail string) (commitSha string, err error) {
	if ws.ViewID == nil && ws.LatestSnapshotID != nil {
		commitSha, err = svc.prepareBranchForPullRequestFromSnapshot(ctx, prBranchName, ws, commitMessage, userName, userEmail)
		if err != nil {
			return "", fmt.Errorf("failed to prepare branch from snapshot: %w", err)
		}
		return
	} else if ws.ViewID != nil {
		commitSha, err = svc.prepareBranchForPullRequestWithView(ctx, prBranchName, ws, commitMessage, userName, userEmail)
		if err != nil {
			return "", fmt.Errorf("failed to prepare branch from snapshot: %w", err)
		}
		return
	} else {
		return "", errors.New("workspace does not have either view nor snapshot")
	}
}

func (svc *EnterpriseService) prepareBranchForPullRequestFromSnapshot(ctx context.Context, prBranchName string, ws *workspaces.Workspace, commitMessage, userName, userEmail string) (string, error) {
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
		sha, err := r.CreateNewCommitBasedOnCommit(prBranchName, snapshot.CommitSHA, signature, commitMessage)
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

func (svc *EnterpriseService) prepareBranchForPullRequestWithView(ctx context.Context, prBranchName string, ws *workspaces.Workspace, commitMessage, userName, userEmail string) (string, error) {
	signature := git.Signature{
		Name:  userName,
		Email: userEmail,
		When:  time.Now(),
	}

	var resSha string

	exec := svc.executorProvider.New().FileReadGitWrite(func(r vcs.RepoReaderGitWriter) error {
		treeID, err := vcs_change.CreateChangesTreeFromPatches(ctx, svc.logger, r, ws.CodebaseID, nil)
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

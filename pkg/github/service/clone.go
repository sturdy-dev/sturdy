package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"mash/pkg/codebase"
	"mash/pkg/github"
	"mash/pkg/shortid"
	"strings"
	"time"

	"mash/pkg/codebase/vcs"
	ghappclient "mash/pkg/github/client"
	"mash/pkg/view/events"
	"mash/vcs/provider"

	"go.uber.org/zap"
)

// Clone processes events to initiate the initial cloning of a codebase
//
// When a webhook event is received from GitHub that the app has been installed on a repository
// 1) The webhook handler creates a Codebase{isReady: false}, a partial GitHubRepository, and potentially a CodebaseUser
//    for the user that initiated the installation (if one can be identified)
// 2) An event is sent to this queue
// 3) This worker, clones the repository, populates the GitHubRepository, discovers all users that should have access,
//    and last sets isReady to true.
func (svc *Service) Clone(
	codebaseID string,
	installationID int64,
	gitHubRepositoryID int64,
	senderUserID string,
) error {
	installation, err := svc.gitHubInstallationRepo.GetByInstallationID(installationID)
	if err != nil {
		return fmt.Errorf("could not get github installation: %w", err)
	}

	logger := svc.logger.With(zap.String("codebase_id", codebaseID))

	ctx := context.Background()

	accessToken, err := ghappclient.GetFirstAccessToken(ctx, svc.gitHubAppConfig, installation, gitHubRepositoryID, svc.gitHubClientProvider)
	if err != nil {
		logger.Error("temporary log: could not get github access token", zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "The permissions requested are not granted to this installation") {
			// We don't have permissions to clone, nothing to do about it
			return nil
		}

		return fmt.Errorf("could not get github access token: %w", err)
	}

	tokenClient, _, err := svc.gitHubClientProvider(svc.gitHubAppConfig, installationID)
	if err != nil {
		return fmt.Errorf("could not get github client token: %w", err)
	}

	gitHubRepoDetails, _, err := tokenClient.Repositories.GetByID(ctx, gitHubRepositoryID)
	if err != nil {
		return fmt.Errorf("could not get repo details: %w", err)
	}

	logger = logger.With(
		zap.String("repo_owner", gitHubRepoDetails.GetOwner().GetLogin()),
		zap.String("repo_name", gitHubRepoDetails.GetName()),
	)

	// Populate the GitHubRepository
	sturdyGitHubRepo, err := svc.gitHubRepositoryRepo.GetByCodebaseID(codebaseID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not get from gitHubRepositoryRepo: %w", err)
	}

	sturdyGitHubRepo.TrackedBranch = gitHubRepoDetails.GetDefaultBranch()
	sturdyGitHubRepo.InstallationAccessToken = accessToken.Token
	sturdyGitHubRepo.InstallationAccessTokenExpiresAt = accessToken.ExpiresAt
	if err = svc.gitHubRepositoryRepo.Update(sturdyGitHubRepo); err != nil {
		return fmt.Errorf("the githubrepository could not be updated: %w", err)
	}

	logger.Info("cloning github repository")

	if err := svc.executorProvider.New().AllowRebasingState().Schedule(func(repoProvider provider.RepoProvider) error {
		return vcs.CloneFromGithub(repoProvider, codebaseID, gitHubRepoDetails.GetOwner().GetLogin(), gitHubRepoDetails.GetName(), *accessToken.Token)
	}).ExecTrunk(codebaseID, "clone github repository"); err != nil {
		return fmt.Errorf("cloning failed: %w", err)
	}

	// Import pull requests by the sender
	if senderUserID != "" {
		// enqueue import pull requests for this user
		if err := svc.EnqueueGitHubPullRequestImport(ctx, codebaseID, senderUserID); err != nil {
			return fmt.Errorf("failed to add to pr importer queue: %w", err)
		}
	}

	cb, err := svc.codebaseRepo.Get(codebaseID)
	if err != nil {
		return fmt.Errorf("could not get codebase: %w", err)
	}

	cb.IsReady = true
	if err := svc.codebaseRepo.Update(cb); err != nil {
		return fmt.Errorf("failed to mark codebase as ready: %w", err)
	}

	// Grant access for other users (the sender already has access)
	if err := svc.GrantCollaboratorsAccess(ctx, codebaseID, strOrNilIfEmpty(senderUserID)); err != nil {
		return fmt.Errorf("failed to GrantCollaboratorsAccess: %w", err)
	}

	// Send events
	svc.eventsSender.Codebase(cb.ID, events.CodebaseUpdated, cb.ID)

	logger.Info("successfully cloned repository, and marked it as ready!")

	return nil
}

func strOrNilIfEmpty(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func (svc *Service) CloneMissingRepositories(ctx context.Context, userID string) error {
	existingGitHubUser, err := svc.gitHubUserRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get github user: %w", err)
	}

	bgCtx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: existingGitHubUser.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	userAuthClient := gh.NewClient(tc)

	installations, err := svc.listAllUserInstallations(bgCtx, userAuthClient)
	if err != nil {
		return fmt.Errorf("failed to lookup installations for user: %w", err)
	}

	currentGitHubUser, _, err := userAuthClient.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get github user: %w", err)
	}

	for _, installation := range installations {
		// if the installation is missing, create it!
		_, err := svc.gitHubInstallationRepo.GetByInstallationID(installation.GetID())
		if errors.Is(err, sql.ErrNoRows) {
			if err := svc.gitHubInstallationRepo.Create(github.GitHubInstallation{
				ID:                     uuid.NewString(),
				InstallationID:         installation.GetID(),
				Owner:                  installation.GetAccount().GetLogin(),
				CreatedAt:              time.Now(),
				HasWorkflowsPermission: true, // This is an assumption
			}); err != nil {
				return fmt.Errorf("failed to re-create installation entry: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to get existing installation: %w", err)
		}

		repoIDs, err := svc.userAccessibleRepoIDs(ctx, userAuthClient, installation.GetID())
		if err != nil {
			return fmt.Errorf("failed to list repo ids: %w", err)
		}

		for _, repoID := range repoIDs {
			_, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installation.GetID(), repoID)
			if err == nil {
				continue
			}

			ghRepo, _, err := userAuthClient.Repositories.GetByID(ctx, repoID)
			if err != nil {
				return fmt.Errorf("could not get github repo: %w", err)
			}

			if err := svc.CreateNonReadyCodebaseAndClone(ctx, ghRepo, installation.GetID(), currentGitHubUser); err != nil {
				return fmt.Errorf("could enqueue clone: %w", err)
			}
		}
	}
	return nil
}

func (svc *Service) CreateNonReadyCodebaseAndClone(ctx context.Context, ghRepo *gh.Repository, installationID int64, sender *gh.User) error {
	svc.logger.Info("handleInstalledRepository setting up new non-ready codebase", zap.Int64("installation_ID", installationID), zap.String("gh_repo_name", ghRepo.GetName()))

	nonReadyCodebase := codebase.Codebase{
		ID:              uuid.NewString(),
		Name:            ghRepo.GetName(),
		ShortCodebaseID: codebase.ShortCodebaseID(shortid.New()),
		Description:     ghRepo.GetDescription(),
		IsReady:         false,
	}
	if err := svc.codebaseRepo.Create(nonReadyCodebase); err != nil {
		return fmt.Errorf("failed to create non-ready codebase: %w", err)
	}

	sturdyGitHubRepo := github.GitHubRepository{
		ID:                 uuid.NewString(),
		InstallationID:     installationID,
		Name:               ghRepo.GetName(),
		GitHubRepositoryID: ghRepo.GetID(),
		CreatedAt:          time.Now(),
		CodebaseID:         nonReadyCodebase.ID,
	}

	if err := svc.gitHubRepositoryRepo.Create(sturdyGitHubRepo); err != nil {
		return fmt.Errorf("failed to save new repo installation: %w", err)
	}

	// Grant access to the initiator right away
	// Access for other users will be added by the worker
	var senderUserID string
	if sender != nil {
		if gitHubUser, err := svc.gitHubUserRepo.GetByUsername(sender.GetLogin()); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get github user: %w", err)
		} else if err == nil {
			senderUserID = gitHubUser.UserID
			if err := svc.AddUser(nonReadyCodebase.ID, gitHubUser, &sturdyGitHubRepo); err != nil {
				return fmt.Errorf("failed to add sender to repo: %w", err)
			}
		}
	}

	// Put to queue!
	if err := svc.gitHubCloneQueue.Enqueue(ctx, &github.CloneRepositoryEvent{
		CodebaseID:         nonReadyCodebase.ID,
		InstallationID:     installationID,
		GitHubRepositoryID: ghRepo.GetID(),
		SenderUserID:       senderUserID,
	}); err != nil {
		return fmt.Errorf("failed to send EnqueueGitHubClone: %w", err)
	}

	return nil
}

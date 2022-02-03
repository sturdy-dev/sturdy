package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/codebase/vcs"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/github"
	ghappclient "getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/shortid"
	"getsturdy.com/api/vcs/provider"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
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

	if err := svc.executorProvider.New().
		AllowRebasingState(). // allowed because the repo does not exist yet
		Schedule(func(repoProvider provider.RepoProvider) error {
			return vcs.CloneFromGithub(logger, repoProvider, codebaseID, gitHubRepoDetails, *accessToken.Token)
		}).ExecTrunk(codebaseID, "clone github repository"); err != nil {
		return fmt.Errorf("cloning failed: %w", err)
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
	repos, err := svc.ListAllAccessibleRepositoriesFromGitHub(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to list repos from github: %w", err)
	}

	existingGitHubUser, err := svc.gitHubUserRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get github user: %w", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: existingGitHubUser.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	userAuthClient := gh.NewClient(tc)

	currentGitHubUser, _, err := userAuthClient.Users.Get(ctx, "")

	for _, repo := range repos {
		// if the installation is missing, create it!
		_, err := svc.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
		if errors.Is(err, sql.ErrNoRows) {
			if err := svc.gitHubInstallationRepo.Create(github.GitHubInstallation{
				ID:                     uuid.NewString(),
				InstallationID:         repo.InstallationID,
				Owner:                  repo.Owner,
				CreatedAt:              time.Now(),
				HasWorkflowsPermission: true, // This is an assumption
			}); err != nil {
				return fmt.Errorf("failed to re-create installation entry: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to get existing installation: %w", err)
		}

		_, err = svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(repo.InstallationID, repo.RepositoryID)
		if err == nil {
			continue
		}

		ghRepo, _, err := userAuthClient.Repositories.GetByID(ctx, repo.RepositoryID)
		if err != nil {
			return fmt.Errorf("could not get github repo: %w", err)
		}

		if _, err := svc.CreateNonReadyCodebaseAndClone(ctx, ghRepo, repo.InstallationID, currentGitHubUser, nil, nil); err != nil {
			return fmt.Errorf("could enqueue clone: %w", err)
		}
	}

	return nil
}

type GitHubRepo struct {
	InstallationID int64
	RepositoryID   int64

	Owner string
	Name  string
}

func (svc *Service) ListAllAccessibleRepositoriesFromGitHub(ctx context.Context, userID string) ([]GitHubRepo, error) {
	existingGitHubUser, err := svc.gitHubUserRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get github user: %w", err)
	}

	bgCtx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: existingGitHubUser.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	userAuthClient := gh.NewClient(tc)

	installations, err := svc.listAllUserInstallations(bgCtx, userAuthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup installations for user: %w", err)
	}

	var res []GitHubRepo

	for _, installation := range installations {
		repos, err := svc.userAccessibleRepoIDs(ctx, userAuthClient, installation.GetID())
		if err != nil {
			return nil, fmt.Errorf("failed to list repo ids: %w", err)
		}

		for _, repo := range repos {
			res = append(res, GitHubRepo{
				InstallationID: installation.GetID(),
				Owner:          installation.GetAccount().GetLogin(),

				RepositoryID: repo.id,
				Name:         repo.name,
			})
		}
	}

	return res, nil
}

func (svc *Service) CreateNonReadyCodebaseAndCloneByIDs(ctx context.Context, installationID, repositoryID int64, userID, organizationID string) (*codebase.Codebase, error) {
	client, _, err := ghappclient.NewClient(svc.gitHubAppConfig, installationID)
	if err != nil {
		return nil, fmt.Errorf("could not get github client: %w", err)
	}

	repo, _, err := client.Repositories.GetByID(ctx, repositoryID)
	if err != nil {
		return nil, fmt.Errorf("could not get repo details from github: %w", err)
	}

	return svc.CreateNonReadyCodebaseAndClone(ctx, repo, installationID, nil, &userID, &organizationID)
}

func (svc *Service) CreateNonReadyCodebaseAndClone(ctx context.Context, ghRepo *gh.Repository, installationID int64, sender *gh.User, addUserID *string, organizationID *string) (*codebase.Codebase, error) {
	svc.logger.Info("handleInstalledRepository setting up new non-ready codebase", zap.Int64("installation_ID", installationID), zap.String("gh_repo_name", ghRepo.GetName()))

	// Create the installation if it does not exist
	_, existingInstallationErr := svc.gitHubInstallationRepo.GetByInstallationID(installationID)
	switch {
	case errors.Is(existingInstallationErr, sql.ErrNoRows):
		_, jwtClient, err := ghappclient.NewClient(svc.gitHubAppConfig, installationID)
		if err != nil {
			return nil, fmt.Errorf("could not get github client: %w", err)
		}
		installation, _, err := jwtClient.Apps.GetInstallation(ctx, installationID)
		if err != nil {
			return nil, fmt.Errorf("could not get installation metadata from github: %w", err)
		}
		svc.gitHubInstallationRepo.Create(github.GitHubInstallation{
			ID:                     uuid.NewString(),
			InstallationID:         installationID,
			Owner:                  installation.GetAccount().GetLogin(),
			CreatedAt:              time.Now(),
			HasWorkflowsPermission: true, // this is a guess
		})
	case existingInstallationErr != nil:
		return nil, fmt.Errorf("could not get installation from repo: %w", existingInstallationErr)
	}

	nonReadyCodebase := codebase.Codebase{
		ID:              uuid.NewString(),
		Name:            ghRepo.GetName(),
		ShortCodebaseID: codebase.ShortCodebaseID(shortid.New()),
		Description:     ghRepo.GetDescription(),
		IsReady:         false,
		OrganizationID:  organizationID, // Optional (for now)
	}
	if err := svc.codebaseRepo.Create(nonReadyCodebase); err != nil {
		return nil, fmt.Errorf("failed to create non-ready codebase: %w", err)
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
		return nil, fmt.Errorf("failed to save new repo installation: %w", err)
	}

	// Grant access to the initiator right away
	// Access for other users will be added by the worker
	var senderUserID string
	if sender != nil {
		if gitHubUser, err := svc.gitHubUserRepo.GetByUsername(sender.GetLogin()); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get github user: %w", err)
		} else if err == nil {
			senderUserID = gitHubUser.UserID
			if err := svc.AddUser(nonReadyCodebase.ID, gitHubUser.UserID); err != nil {
				return nil, fmt.Errorf("failed to add sender to repo: %w", err)
			}
		}
	} else if addUserID != nil {
		senderUserID = *addUserID
		if err := svc.AddUser(nonReadyCodebase.ID, *addUserID); err != nil {
			return nil, fmt.Errorf("failed to add user to repo: %w", err)
		}
	}

	// Put to queue!
	if err := (*svc.gitHubCloneQueue).Enqueue(ctx, &github.CloneRepositoryEvent{
		CodebaseID:         nonReadyCodebase.ID,
		InstallationID:     installationID,
		GitHubRepositoryID: ghRepo.GetID(),
		SenderUserID:       senderUserID,
	}); err != nil {
		return nil, fmt.Errorf("failed to send EnqueueGitHubClone: %w", err)
	}

	return &nonReadyCodebase, nil
}

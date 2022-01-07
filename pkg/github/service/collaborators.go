package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"mash/pkg/codebase"
	sturdy_github "mash/pkg/github"
	"mash/pkg/github/client"
	"mash/pkg/notification"
	"mash/pkg/view/events"

	"github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

func (svc *Service) GrantCollaboratorsAccess(ctx context.Context, codebaseID string, authAsUserID *string) error {
	var didInviteAny bool

	gitHubRepo, err := svc.gitHubRepositoryRepo.GetByCodebaseID(codebaseID)
	if err != nil {
		return fmt.Errorf("failed to get github repo by codebase: %w", err)
	}

	installation, err := svc.gitHubInstallationRepo.GetByInstallationID(gitHubRepo.InstallationID)
	if err != nil {
		return fmt.Errorf("failed to get installation: %w", err)
	}

	repoClient, err := svc.authAsUserOrFallbackAsApp(authAsUserID, installation.InstallationID)
	if err != nil {
		return fmt.Errorf("failed to get a github client: %w", err)
	}

	collaborators, err := listAllCollaborators(ctx, repoClient, installation.Owner, gitHubRepo.Name)
	if err != nil {
		return fmt.Errorf("failed to list collaborators: %w", err)
	}

	for _, collaborator := range collaborators {
		logger := svc.logger.With(
			zap.String("codebase_id", collaborator.GetLogin()),
			zap.String("github_login", collaborator.GetLogin()),
		)

		logger.Info("setting up collaborator")

		gitHubUser, err := svc.gitHubUserRepo.GetByUsername(collaborator.GetLogin())
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logger.Error("failed to get github user from db", zap.Error(err))
			}
			continue
		}

		logger = logger.With(zap.String("user_id", gitHubUser.UserID))

		// If the gitHubUser was created within the last hour, this is a new github connection.
		// Only send notifications for old connections that have a new github repo imported.
		createdWithinTheLastHour := gitHubUser.CreatedAt.Add(time.Hour).After(time.Now())
		if !createdWithinTheLastHour {
			if err := svc.notificationSender.User(ctx, gitHubUser.UserID, gitHubRepo.CodebaseID, notification.GitHubRepositoryImported, gitHubRepo.ID); err != nil {
				logger.Error("failed to send github repo imported notification", zap.Error(err))
			} else {
				logger.Info("sent notification about the imported codebase")
			}
		} else {
			logger.Info("github user is to new, skipping sending notification")
		}

		_, err = svc.codebaseUserRepo.GetByUserAndCodebase(gitHubUser.UserID, codebaseID)
		switch {
		case err == nil:
			// The user is already a member (and is likely the user that installed the repo)
			logger.Info("github user is already a member of the codebase")

			// enqueue import pull requests for this user
			if err := svc.EnqueueGitHubPullRequestImport(ctx, codebaseID, gitHubUser.UserID); err != nil {
				logger.Error("failed to add to pr importer queue", zap.Error(err))
			}
		case errors.Is(err, sql.ErrNoRows):

			logger.Info("granting access to repository based on GitHub credentials")

			t0 := time.Now()
			if err := svc.codebaseUserRepo.Create(codebase.CodebaseUser{
				ID:         uuid.NewString(),
				UserID:     gitHubUser.UserID,
				CodebaseID: codebaseID,
				CreatedAt:  &t0,
			}); err != nil {
				logger.Warn("failed to create codebase-user relation in db", zap.Error(err))
				continue
			}

			if err := svc.postHogClient.Enqueue(posthog.Capture{
				Event:      "added user to codebase",
				DistinctId: gitHubUser.UserID,
				Properties: posthog.NewProperties().
					Set("github", true).
					Set("is_github_sender", false). // This event is not fired for the user that installed the GitHub app
					Set("codebase_id", codebaseID),
			}); err != nil {
				logger.Error("posthog failed", zap.Error(err))
			}

			// enqueue import pull requests for this user
			if err := svc.EnqueueGitHubPullRequestImport(ctx, codebaseID, gitHubUser.UserID); err != nil {
				logger.Error("failed to add to pr importer queue", zap.Error(err))
			}

			didInviteAny = true

		default:
			logger.Error("failed to get codebase-user relation from db", zap.Error(err))
			continue
		}
	}

	if didInviteAny {
		// Send events
		svc.eventsSender.Codebase(codebaseID, events.CodebaseUpdated, codebaseID)
	}

	return nil
}

func (svc *Service) authAsUserOrFallbackAsApp(userID *string, installationID int64) (client.RepositoriesClient, error) {
	// Prefer user auth
	if userID != nil {
		gitHubUser, err := svc.gitHubUserRepo.GetByUserID(*userID)
		// Auth as user if a user could be found
		if err == nil {
			personalClient, err := svc.gitHubPersonalClientProvider(gitHubUser.AccessToken)
			if err != nil {
				return nil, fmt.Errorf("failed to create personal github client: %w", err)
			}
			return personalClient.Repositories, nil
		}
	}

	// Fallback to authenticating as the app, note tough that this is a worse option. As requests from the app might not see all users.
	tokenClient, _, err := svc.gitHubClientProvider(
		svc.gitHubAppConfig,
		installationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}
	return tokenClient.Repositories, nil

}

// listAllCollaborators returns a list of collaborators that the authenticated user can _see_.
// Note that the app does not have collaborator access, and will only see users that have a _public_ membership.
// To get the full list of collaborators, authenticate as a user with a confirmed membership (the user that installed the app is a good candidate).
func listAllCollaborators(ctx context.Context, reposClient client.RepositoriesClient, owner, name string) ([]*github.User, error) {
	var users []*github.User
	page := 1
	for page != 0 {
		newUsers, nextPage, err := listCollaborators(ctx, reposClient, owner, name, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		users = append(users, newUsers...)
	}
	return users, nil
}

func listCollaborators(ctx context.Context, reposClient client.RepositoriesClient, owner, name string, page int) ([]*github.User, int, error) {
	users, rsp, err := reposClient.ListCollaborators(ctx, owner, name, &github.ListCollaboratorsOptions{
		Affiliation: "all",
		ListOptions: github.ListOptions{Page: page, PerPage: 50}},
	)
	if err != nil {
		return nil, 0, err
	}
	return users, rsp.NextPage, nil
}

func (svc *Service) AddUser(codebaseID string, gitHubUser *sturdy_github.GitHubUser, gitHubRepo *sturdy_github.GitHubRepository) error {
	// Add access to this user directly
	t := time.Now()
	err := svc.codebaseUserRepo.Create(codebase.CodebaseUser{
		ID:         uuid.NewString(),
		UserID:     gitHubUser.UserID,
		CodebaseID: codebaseID,
		CreatedAt:  &t,
	})
	if err != nil {
		return fmt.Errorf("failed to add sender to codebaseUserRepo: %w", err)
	}

	err = svc.postHogClient.Enqueue(posthog.Capture{
		Event:      "added user to codebase",
		DistinctId: gitHubUser.UserID,
		Properties: posthog.NewProperties().
			Set("github", true).
			Set("is_github_sender", true).
			Set("codebase_id", codebaseID),
	})
	if err != nil {
		svc.logger.Error("posthog failed", zap.Error(err))
	}

	svc.logger.Info("adding github sender to the codebase", zap.String("user_id", gitHubUser.UserID))

	// Track installation event on the user that installed it
	err = svc.postHogClient.Enqueue(posthog.Capture{
		Event:      "installed github repository",
		DistinctId: gitHubUser.UserID,
		Properties: posthog.NewProperties().
			Set("github", true).
			Set("is_github_sender", true).
			Set("codebase_id", codebaseID),
	})
	if err != nil {
		svc.logger.Error("posthog failed", zap.Error(err))
	}

	// Send events
	svc.eventsSender.Codebase(codebaseID, events.CodebaseUpdated, codebaseID)

	return nil
}

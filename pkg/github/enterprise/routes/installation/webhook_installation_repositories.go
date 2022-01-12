package installation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	db_codebase "mash/pkg/codebase/db"
	"mash/pkg/github"
	"mash/pkg/github/enterprise/db"
	service_github "mash/pkg/github/enterprise/service"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

func HandleInstallationRepositoriesEvent(
	ctx context.Context,
	logger *zap.Logger,
	event *gh.InstallationRepositoriesEvent,
	gitHubAppInstallationsRepository db.GitHubInstallationRepo,
	gitHubAppInstalledRepositoryRepository db.GitHubRepositoryRepo,
	postHogClient posthog.Client,
	codebaseRepo db_codebase.CodebaseRepository,
	gitHubService *service_github.Service,
) error {
	_, err := gitHubAppInstallationsRepository.GetByInstallationID(event.Installation.GetID())
	// If the original InstallationEvent webhook was missed (otherwise user has to remove and re-add app)
	if errors.Is(err, sql.ErrNoRows) {
		err := gitHubAppInstallationsRepository.Create(github.GitHubInstallation{
			ID:                     uuid.NewString(),
			InstallationID:         event.Installation.GetID(),
			Owner:                  event.Installation.Account.GetLogin(), // "sturdy-dev" or "zegl", etc.
			CreatedAt:              time.Now(),
			HasWorkflowsPermission: true,
		})
		if err != nil {
			return err
		}
	}

	if event.GetRepositorySelection() == "selected" || event.GetRepositorySelection() == "all" {
		// Add new repos
		for _, r := range event.RepositoriesAdded {
			err := handleInstalledRepository(
				ctx,
				logger,
				event.GetInstallation().GetID(),
				r,
				event.GetSender(),
				gitHubAppInstalledRepositoryRepository,
				postHogClient,
				codebaseRepo,
				gitHubService,
			)
			if err != nil {
				return err
			}
		}

		// Mark uninstalled repos as uninstalled
		for _, r := range event.RepositoriesRemoved {
			installedRepo, err := gitHubAppInstalledRepositoryRepository.GetByInstallationAndGitHubRepoID(
				event.GetInstallation().GetID(),
				r.GetID(),
			)
			if err != nil {
				log.Println("failed to mark as uninstalled", err)
				continue
			}

			t := time.Now()
			installedRepo.UninstalledAt = &t

			err = gitHubAppInstalledRepositoryRepository.Update(installedRepo)
			if err != nil {
				log.Println("failed to mark as uninstalled", err)
				continue
			}
		}
	}

	return nil
}

func handleInstalledRepository(
	ctx context.Context,
	logger *zap.Logger,
	installationID int64,
	ghRepo *gh.Repository,
	sender *gh.User,
	gitHubRepositoryRepo db.GitHubRepositoryRepo,
	postHogClient posthog.Client,
	codebaseRepo db_codebase.CodebaseRepository,
	gitHubService *service_github.Service,
) error {
	// Create a non-ready codebase (add the initiating user), and put the event on a queue
	logger = logger.With(zap.String("repo_name", ghRepo.GetName()), zap.Int64("installation_id", installationID))
	logger.Info("handleInstalledRepository")

	// Tracking on the GitHub installation itself, there is also some tracking on the user
	err := postHogClient.Enqueue(posthog.Capture{
		DistinctId: fmt.Sprintf("%d", installationID), // Using the installation ID as a person?
		Event:      "installed repository",
		Properties: posthog.NewProperties().
			Set("installation_id", installationID).
			Set("repo_name", ghRepo.GetName()),
	})
	if err != nil {
		logger.Error("posthog post failed", zap.Error(err))
	}

	existingGitHubRepo, err := gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, ghRepo.GetID())
	// If GitHub repo already exists (previously installed and then uninstalled) un-archive it
	if err == nil {
		logger.Info("handleInstalledRepository repository already exists", zap.Any("existing_github_repo", existingGitHubRepo))

		// un-archive the codebase if archived
		cb, err := codebaseRepo.GetAllowArchived(existingGitHubRepo.CodebaseID)
		if err != nil {
			logger.Error("failed to get codebase", zap.Error(err))
			return err
		}

		if cb.ArchivedAt != nil {
			cb.ArchivedAt = nil
			if err := codebaseRepo.Update(cb); err != nil {
				logger.Error("failed to un-archive codebase", zap.Error(err))
				return err
			}
		}

		if existingGitHubRepo.UninstalledAt != nil {
			existingGitHubRepo.UninstalledAt = nil
			err := gitHubRepositoryRepo.Update(existingGitHubRepo)
			if err != nil {
				logger.Error("failed to un-archive github repository entry", zap.Error(err))
				return err
			}
		}

		return nil
	}

	if err := gitHubService.CreateNonReadyCodebaseAndClone(ctx, ghRepo, installationID, sender); err != nil {
		return fmt.Errorf("failed to create non-ready codebase: %w", err)
	}

	return nil
}

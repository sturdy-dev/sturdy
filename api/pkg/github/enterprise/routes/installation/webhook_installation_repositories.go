package installation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func HandleInstallationRepositoriesEvent(
	ctx context.Context,
	logger *zap.Logger,
	event *gh.InstallationRepositoriesEvent,
	gitHubAppInstallationsRepository db.GitHubInstallationRepo,
	gitHubAppInstalledRepositoryRepository db.GitHubRepositoryRepo,
	analyticsService *service_analytics.Service,
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
				analyticsService,
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
	analyticsService *service_analytics.Service,
	codebaseRepo db_codebase.CodebaseRepository,
	gitHubService *service_github.Service,
) error {
	// Create a non-ready codebase (add the initiating user), and put the event on a queue
	logger = logger.With(zap.String("repo_name", ghRepo.GetName()), zap.Int64("installation_id", installationID))
	logger.Info("handleInstalledRepository")

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

		analyticsService.Capture(ctx, "installed repository", analytics.DistinctID(fmt.Sprintf("%d", installationID)),
			analytics.CodebaseID(cb.ID),
			analytics.Property("repo_name", ghRepo.GetName()),
		)

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

	return nil
}

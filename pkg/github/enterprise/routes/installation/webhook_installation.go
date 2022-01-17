package installation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"mash/pkg/analytics"
	db_codebase "mash/pkg/codebase/db"
	"mash/pkg/github"
	"mash/pkg/github/enterprise/db"
	service_github "mash/pkg/github/enterprise/service"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func HandleInstallationEvent(
	ctx context.Context,
	logger *zap.Logger,
	event *gh.InstallationEvent,
	gitHubInstallationRepo db.GitHubInstallationRepo,
	gitHubRepositoryRepo db.GitHubRepositoryRepo,
	analyticsClient analytics.Client,
	codebaseRepo db_codebase.CodebaseRepository,
	gitHubService *service_github.Service,
) error {
	log.Printf("GitHub Installed: %+v", event)

	// Identify the installation to analytics
	err := analyticsClient.Enqueue(analytics.Identify{
		DistinctId: fmt.Sprintf("%d", event.GetInstallation().GetID()), // Using the installation ID as a person?
		Properties: analytics.NewProperties().
			Set("installation_org", event.Installation.Account.GetLogin()).
			Set("email", event.Installation.Account.GetLogin()).
			Set("github_app_installation", true), // To differentiate between humans and installations
	})
	if err != nil {
		logger.Error("analytics post failed", zap.Error(err))
	}

	// Track with analytics
	err = analyticsClient.Enqueue(analytics.Capture{
		DistinctId: fmt.Sprintf("%d", event.GetInstallation().GetID()), // Using the installation ID as a person?
		Event:      fmt.Sprintf("github installation %s", event.GetAction()),
	})
	if err != nil {
		logger.Error("analytics post failed", zap.Error(err))
	}

	if event.GetAction() == "created" ||
		event.GetAction() == "deleted" ||
		event.GetAction() == "new_permissions_accepted" {

		t := time.Now()
		var uninstalledAt *time.Time
		if event.GetAction() == "deleted" {
			uninstalledAt = &t
		}

		// Check if it's already installed (stale data or whatever)
		existing, err := gitHubInstallationRepo.GetByInstallationID(event.Installation.GetID())
		if errors.Is(err, sql.ErrNoRows) {
			// Save new installation
			err := gitHubInstallationRepo.Create(github.GitHubInstallation{
				ID:                     uuid.NewString(),
				InstallationID:         event.Installation.GetID(),
				Owner:                  event.Installation.GetAccount().GetLogin(), // "sturdy-dev" or "zegl", etc.
				CreatedAt:              t,
				UninstalledAt:          uninstalledAt,
				HasWorkflowsPermission: true,
			})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// Update existing entry
			existing.UninstalledAt = uninstalledAt

			if event.GetAction() == "new_permissions_accepted" {
				existing.HasWorkflowsPermission = true
			}

			err = gitHubInstallationRepo.Update(existing)
			if err != nil {
				return err
			}
		}

		// Save all repositories
		if event.GetAction() == "created" {
			for _, r := range event.Repositories {
				err = handleInstalledRepository(
					ctx,
					logger,
					event.GetInstallation().GetID(),
					r,
					event.GetSender(),
					gitHubRepositoryRepo,
					analyticsClient,
					codebaseRepo,
					gitHubService,
				)
				if err != nil {
					return err
				}
			}
		}
		// TODO: Handle deleted - update installed repositories?

		return nil
	}

	// Unhandled actions:
	// suspend, unsuspend

	return nil
}

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

func HandleInstallationEvent(
	ctx context.Context,
	logger *zap.Logger,
	event *gh.InstallationEvent,
	gitHubInstallationRepo db.GitHubInstallationRepo,
	gitHubRepositoryRepo db.GitHubRepositoryRepo,
	analyticsServcie *service_analytics.Service,
	codebaseRepo db_codebase.CodebaseRepository,
	gitHubService *service_github.Service,
) error {
	log.Printf("GitHub Installed: %+v", event)

	analyticsServcie.IdentifyGitHubInstallation(ctx, event.GetInstallation())

	analyticsServcie.Capture(ctx, fmt.Sprintf("github installation %s", event.GetAction()),
		analytics.DistinctID(fmt.Sprintf("%d", event.GetInstallation().GetID())),
	)

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
					analyticsServcie,
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

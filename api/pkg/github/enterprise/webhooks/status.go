package webhooks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"getsturdy.com/api/pkg/statuses"
)

func getStatusType(event *StatusEvent) statuses.Type {
	state := event.GetState() // "pending", "success", "failure", "error"
	switch state {
	case "pending":
		return statuses.TypePending
	case "failure", "error":
		return statuses.TypeFailing
	case "success":
		return statuses.TypeHealthy
	default:
		return statuses.TypeUndefined
	}
}

func getStatusTime(event *StatusEvent) time.Time {
	if event.UpdatedAt != nil {
		return event.GetUpdatedAt().Time
	}
	return event.GetCreatedAt().Time
}

func (svc *Service) HandleStatusEvent(ctx context.Context, event *StatusEvent) error {
	gitHubRepoID := event.GetRepo().GetID()
	installationID := event.GetInstallation().GetID()
	repo, err := svc.gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, gitHubRepoID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil
	default:
		return fmt.Errorf("failed to get repository by id: %w", err)
	}

	status := &statuses.Status{
		ID:          uuid.New().String(),
		CommitSHA:   event.GetCommit().GetSHA(),
		CodebaseID:  repo.CodebaseID,
		Type:        getStatusType(event),
		Title:       event.GetContext(),
		Description: event.Description,
		DetailsURL:  event.TargetURL,
		Timestamp:   getStatusTime(event),
	}

	if err := svc.statusService.Set(ctx, status); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

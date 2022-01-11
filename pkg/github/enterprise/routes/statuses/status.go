package statuses

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	db_github "mash/pkg/github/db"
	"mash/pkg/statuses"
	service_statuses "mash/pkg/statuses/service"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func getStatusType(event *gh.StatusEvent) statuses.Type {
	state := event.GetState() // "pending", "success", "failure", "error"
	switch state {
	case "pending":
		return statuses.TypePending
	case "failure", "error":
		return statuses.TypeFailing
	case "success":
		return statuses.TypeHealty
	default:
		return statuses.TypeUndefined
	}
}

func getStatusTime(event *gh.StatusEvent) time.Time {
	if event.UpdatedAt != nil {
		return event.GetUpdatedAt().Time
	}
	return event.GetCreatedAt().Time
}

func HandleStatusEvent(
	ctx context.Context,
	logger *zap.Logger,
	event *gh.StatusEvent,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	statusesService *service_statuses.Service,
) error {
	gitHubRepoID := event.GetRepo().GetID()
	installationID := event.GetInstallation().GetID()
	repo, err := gitHubRepositoryRepo.GetByInstallationAndGitHubRepoID(installationID, gitHubRepoID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil
	default:
		return fmt.Errorf("failed to get repository by id: %w", err)
	}

	status := &statuses.Status{
		ID:          uuid.New().String(),
		CommitID:    event.GetCommit().GetSHA(),
		CodebaseID:  repo.CodebaseID,
		Type:        getStatusType(event),
		Title:       event.GetContext(),
		Description: event.Description,
		DetailsURL:  event.TargetURL,
		Timestamp:   getStatusTime(event),
	}

	if err := statusesService.Set(ctx, status); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

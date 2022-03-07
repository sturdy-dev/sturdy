package webhooks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/statuses"
)

var statusPending = map[string]bool{
	"queued":      true,
	"in_progress": true,
}

var statusCompleted = map[string]bool{
	"completed": true,
}

var conclusionPending = map[string]bool{
	"action_required": true,
}

var conclusionFailing = map[string]bool{
	"failure":   true,
	"timed_out": true,
}

var conclusionHealthy = map[string]bool{
	"success": true,
}

func getJobTime(job *WorkflowJob) time.Time {
	if job.CompletedAt != nil {
		return job.GetCompletedAt().Time
	}
	return job.GetStartedAt().Time
}

func getJobType(job *WorkflowJob) statuses.Type {
	status := job.GetStatus()         // queued, in_progress, completed
	conclution := job.GetConclusion() // success, failure, neutral, cancelled, timed_out, action_required, stale
	switch {
	case statusPending[status]:
		return statuses.TypePending
	case statusCompleted[status]:
		switch {
		case conclusionPending[conclution]:
			return statuses.TypePending
		case conclusionFailing[conclution]:
			return statuses.TypeFailing
		case conclusionHealthy[conclution]:
			return statuses.TypeHealty
		default:
			return statuses.TypeUndefined
		}
	default:
		return statuses.TypeUndefined
	}
}

func (svc *Service) HandleWorkflowJobEvent(ctx context.Context, event *WorkflowJobEvent) error {
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

	job := event.GetWorkflowJob()

	jobType := getJobType(job)
	if jobType == statuses.TypeUndefined {
		svc.logger.Warn(
			"failed to parse github job type",
			zap.String("repo_id", repo.ID),
			zap.String("status", job.GetStatus()),
			zap.String("conclution", job.GetConclusion()),
		)
		return nil
	}

	status := &statuses.Status{
		ID:         uuid.New().String(),
		CommitID:   job.GetHeadSHA(),
		CodebaseID: repo.CodebaseID,
		Type:       jobType,
		Title:      job.GetName(),
		Timestamp:  getJobTime(job),
		DetailsURL: job.HTMLURL,
	}

	if err := svc.statusService.Set(ctx, status); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

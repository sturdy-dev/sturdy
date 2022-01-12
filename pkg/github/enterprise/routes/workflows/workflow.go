package workflows

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	db_github "mash/pkg/github/enterprise/db"
	"mash/pkg/statuses"
	service_statuses "mash/pkg/statuses/service"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func getJobTime(job *gh.WorkflowJob) time.Time {
	if job.CompletedAt != nil {
		return job.GetCompletedAt().Time
	}
	return job.GetStartedAt().Time
}

func getJobType(job *gh.WorkflowJob) statuses.Type {
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

func HandleWorkflowJobEvent(
	ctx context.Context,
	logger *zap.Logger,
	event *gh.WorkflowJobEvent,
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

	job := event.GetWorkflowJob()

	jobType := getJobType(job)
	if jobType == statuses.TypeUndefined {
		logger.Warn(
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

	if err := statusesService.Set(ctx, status); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

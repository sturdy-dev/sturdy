package pr

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"mash/pkg/github/enterprise/db"
	"mash/pkg/view/events"
	db_workspace "mash/pkg/workspace/db"

	gh "github.com/google/go-github/v39/github"
	"go.uber.org/zap"
)

func HandlePullRequestEvent(logger *zap.Logger, event *gh.PullRequestEvent, workspaceReader db_workspace.WorkspaceReader, gitHubPRRepo db.GitHubPRRepo, eventSender events.EventSender, workspaceWriter db_workspace.WorkspaceWriter) error {
	apiPR := event.GetPullRequest()
	pr, err := gitHubPRRepo.GetByGitHubID(apiPR.GetID())
	if errors.Is(err, sql.ErrNoRows) {
		return nil // noop
	} else if err != nil {
		return fmt.Errorf("failed to get github pull request from db: %w", err)
	}

	if apiPR.GetState() == "closed" || apiPR.GetState() == "open" {
		pr.Open = apiPR.GetState() == "open"
		t0 := time.Now()
		pr.UpdatedAt = &t0
		pr.ClosedAt = apiPR.ClosedAt
		pr.Merged = apiPR.GetMerged()
		pr.MergedAt = apiPR.MergedAt
		if err := gitHubPRRepo.Update(pr); err != nil {
			return fmt.Errorf("failed to update github pull request in db: %w", err)
		}
	}

	ws, err := workspaceReader.Get(pr.WorkspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Warn("handled a github pull request webhook for non-existing workspace", zap.String("workspace_id", pr.WorkspaceID), zap.String("github_pr_id", pr.ID), zap.String("github_pr_link", apiPR.GetHTMLURL()))
		return nil // noop
	} else if err != nil {
		return fmt.Errorf("failed to get workspace from db: %w", err)
	}

	if err := eventSender.Codebase(ws.CodebaseID, events.GitHubPRUpdated, pr.WorkspaceID); err != nil {
		logger.Error("failed to send codebase event", zap.String("workspace_id", pr.WorkspaceID), zap.String("github_pr_id", pr.ID), zap.String("github_pr_link", apiPR.GetHTMLURL()), zap.Error(err))
		// do not fail
	}
	return nil
}

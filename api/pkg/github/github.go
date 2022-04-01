package github

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/users"
)

type Installation struct {
	ID                     string     `db:"id"`
	InstallationID         int64      `db:"installation_id"`
	Owner                  string     `db:"owner"`
	CreatedAt              time.Time  `db:"created_at"`
	UninstalledAt          *time.Time `db:"uninstalled_at"`
	HasWorkflowsPermission bool       `db:"has_workflows_permission"`
}

type Repository struct {
	ID                               string       `db:"id"`
	InstallationID                   int64        `db:"installation_id"`
	Name                             string       `db:"name"`
	GitHubRepositoryID               int64        `db:"github_repository_id"`
	CreatedAt                        time.Time    `db:"created_at"`
	UninstalledAt                    *time.Time   `db:"uninstalled_at"`
	InstallationAccessToken          *string      `db:"installation_access_token"`
	InstallationAccessTokenExpiresAt *time.Time   `db:"installation_access_token_expires_at"`
	TrackedBranch                    string       `db:"tracked_branch"`
	SyncedAt                         *time.Time   `db:"synced_at"`
	CodebaseID                       codebases.ID `db:"codebase_id"`

	// When true, all changes must be made through GitHub, workspaces create Pull Requests
	// When false, changes are made on Sturdy, and sturytrunk pushes to GitHub.
	GitHubSourceOfTruth bool `json:"-" db:"github_source_of_truth"`

	// If the GitHub integration is enabled or not
	IntegrationEnabled bool `json:"-" db:"integration_enabled"`

	LastPushErrorMessage *string    `json:"-" db:"last_push_error_message"`
	LastPushAt           *time.Time `json:"-" db:"last_push_at"`

	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type User struct {
	ID       string   `json:"id" db:"id"`
	UserID   users.ID `json:"user_id" db:"user_id"`
	Username string   `json:"username" db:"username"`
	// Might be nil if the user was imported from github, but never authenticated.
	// For example, via pull requests webhook.
	AccessToken                *string   `json:"-" db:"access_token"`
	AccessTokenLastValidatedAt time.Time `json:"-" db:"access_token_last_validated_at"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
}

type PullRequestState string

const (
	// PullRequestStateUnknown is the default value for a pull request status.
	PullRequestStateUnknown PullRequestState = ""
	// PullRequestStateOpen is the status for a pull request that is open.
	PullRequestStateOpen PullRequestState = "open"
	// PullRequestStateClosed is the status for a pull request that is closed, not merged.
	PullRequestStateClosed PullRequestState = "closed"
	// PullRequestStateMerging is the status for a pull request that was merged through Sturdy,
	// and Sturdy is waiting for GitHub's webhook to complete the merge.
	PullRequestStateMerging PullRequestState = "merging"
	// PullRequestStateMerged is the status for a pull request that is merged, closed.
	PullRequestStateMerged PullRequestState = "merged"
)

type PullRequest struct {
	ID                 string   `db:"id"`
	WorkspaceID        string   `db:"workspace_id"`
	GitHubID           int64    `db:"github_id"`
	GitHubRepositoryID int64    `db:"github_repository_id"`
	CreatedBy          users.ID `db:"created_by"`
	GitHubPRNumber     int      `db:"github_pr_number"`
	Head               string   `db:"head"`
	// HeadSHA is empty for older pull requests.
	HeadSHA    *string          `db:"head_sha"`
	CodebaseID codebases.ID     `db:"codebase_id"`
	Base       string           `db:"base"`
	CreatedAt  time.Time        `db:"created_at"`
	UpdatedAt  *time.Time       `db:"updated_at"`
	ClosedAt   *time.Time       `db:"closed_at"`
	MergedAt   *time.Time       `db:"merged_at"`
	State      PullRequestState `db:"state"`
}

type CloneRepositoryEvent struct {
	CodebaseID         codebases.ID `json:"codebase_id"`
	InstallationID     int64        `json:"installation_id"`
	GitHubRepositoryID int64        `json:"github_repository_id"`
	SenderUserID       users.ID     `json:"sender_user_id"`
}

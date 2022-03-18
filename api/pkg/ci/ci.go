package ci

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
)

// Commit is the link between a commit in the trunk repo, and a commit in the "fake" ci repository that is being
// used to trigger the ci service.
type Commit struct {
	ID             string       `db:"id"`
	CodebaseID     codebases.ID `db:"codebase_id"`
	CiRepoCommitID string       `db:"ci_repo_commit_id"`
	TrunkCommitID  string       `db:"trunk_commit_id"`
	CreatedAt      time.Time    `db:"created_at"`
}

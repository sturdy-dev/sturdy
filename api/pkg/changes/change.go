package changes

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/users"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type Change struct {
	ID                 ID           `db:"id"`
	CodebaseID         codebases.ID `db:"codebase_id"`
	Title              *string      `db:"title"`
	UpdatedDescription string       `db:"updated_description"`
	UserID             *users.ID    `db:"user_id"`

	// Contains id of the workspace that the chage was created in
	// For imported changes, WorkspaceID is nil.
	WorkspaceID *string `db:"workspace_id"`

	// For changes created directly through Sturdy, CreatedAt is the time of creation / share / landing.
	// For changes created with Sturdy via GitHub, CreatedAt is the time when Sturdy imported the Change (via GitHub webhooks)
	// For changes created outside of Sturdy, CreatedAt is null
	CreatedAt *time.Time `db:"created_at"`

	// GitCreatedAt is the timestamp when the change was created to git.
	// This timestamp might be earlier than CreatedAt.
	GitCreatedAt    *time.Time `db:"git_created_at"`
	GitCreatorName  *string    `db:"git_creator_name"`
	GitCreatorEmail *string    `db:"git_creator_email"`

	CommitID *string `db:"commit_id"` // Commit IDs are only unique within a codebase / repository

	// This changes parent.
	// Is null for the first change in a codebase, or if the changes parent hasn't been imported to Sturdy yet.
	ParentChangeID *ID `db:"parent_change_id"`
}

package change

import "time"

type ID string

type Change struct {
	ID                 ID      `db:"id"`
	CodebaseID         string  `db:"codebase_id"`
	Title              *string `db:"title"`
	UpdatedDescription string  `db:"updated_description"`
	UserID             *string `db:"user_id"`

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
}

type ChangeCommit struct {
	ChangeID   ID     `db:"change_id"` // Change IDs are globally unique
	CommitID   string `db:"commit_id"` // Commit IDs are only unique within a codebase / repository
	CodebaseID string `db:"codebase_id"`
	Trunk      bool   `db:"trunk"` // If this commit is on the trunk
}

package comments

import (
	"time"

	"getsturdy.com/api/pkg/change"
	"getsturdy.com/api/pkg/users"
)

type ID string

type Comment struct {
	ID         ID     `db:"id"`
	CodebaseID string `db:"codebase_id"`

	// Either ChangeID is set, and the comment belongs to a change
	// or WorkspaceID is set (when the comment belongs to a live change)
	ChangeID    *change.ID `db:"change_id"`
	WorkspaceID *string    `db:"workspace_id"`

	UserID    users.ID   `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Message   string     `db:"message"`

	// The file that is commented on
	Path    string  `db:"path"`
	OldPath *string `db:"old_path"`

	// The first and last line to be commented on (1-indexed), from the _start_ of the file.
	// Not relative to the start of the hunk.
	LineStart int `db:"line_start"`
	LineEnd   int `db:"line_end"`

	// If the line numbers provided are the _new_ line number, or the old one.
	LineIsNew bool `db:"line_is_new"`

	ContextStartsAtLine *int    `db:"context_starts_at_line"`
	Context             *string `db:"context"`

	ParentComment *ID `db:"parent_comment_id"`
}

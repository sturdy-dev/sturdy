package suggestions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/users"
	"github.com/lib/pq"
)

type ID string

type Suggestion struct {
	ID         ID           `db:"id"`
	CodebaseID codebases.ID `db:"codebase_id"`
	// The id of the workspce that contains suggestions.
	WorkspaceID string `db:"workspace_id"`
	// The id of the workspace that suggestion is being made for.
	ForWorkspaceID string `db:"for_workspace_id"`
	// The id of the snapshot that suggestion is being made for.
	ForSnapshotID snapshots.ID `db:"for_snapshot_id"`
	// Time when the suggestion was created.
	CreatedAt time.Time `db:"created_at"`
	// ID of the suggestion diff hunk ids that were applied.
	AppliedHunks pq.StringArray `db:"applied_hunks"`
	// ID of the suggestion diff hunk ids that were dismissed.
	DismissedHunks pq.StringArray `db:"dismissed_hunks"`
	// The id of the user that created the suggestion.
	UserID users.ID `db:"user_id"`
	// DismissedAt is set if the whole suggestion was dismissed.
	DismissedAt *time.Time `db:"dismissed_at"`
	// NotifiedAt is set if the user has been notified about the suggestion.
	NotifiedAt *time.Time `db:"notified_at"`
}

type Hunk struct {
	FileName string
	Index    int
}

func (a *Hunk) String() string {
	return fmt.Sprintf("%s#%d", a.FileName, a.Index)
}

func ParseAppliedHunkID(in string) (*Hunk, error) {
	idx := strings.LastIndexByte(in, '#')
	if idx < 0 {
		return nil, fmt.Errorf("invalid applied hunk id: %s", in)
	}

	index, err := strconv.Atoi(in[idx+1:])
	if err != nil {
		return nil, fmt.Errorf("invalid applied hunk id: %s, error parsing index: %w", in, err)
	}

	return &Hunk{
		FileName: in[0:idx],
		Index:    index,
	}, nil
}

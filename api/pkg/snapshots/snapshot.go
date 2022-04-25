package snapshots

import (
	"fmt"
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/jsontime"
)

type Snapshot struct {
	ID string `json:"id" db:"id"`
	// CommitSHA with the snapshotted content. I.E. to get a view state without the snapshot, soft reset it to the commitID's parent.
	CommitSHA          string       `json:"-" db:"commit_id"`
	CodebaseID         codebases.ID `json:"codebase_id" db:"codebase_id"`
	ViewID             string       `json:"view_id" db:"view_id"` // ViewID is optional. TODO: Make it nullable?
	WorkspaceID        string       `json:"workspace_id" db:"workspace_id"`
	CreatedAt          time.Time    `json:"-" db:"created_at"`
	PreviousSnapshotID *string      `json:"previous_snapshot_id" db:"previous_snapshot_id"`
	Action             Action       `json:"action" db:"action"`           // The action that triggered the snapshot creations
	DeletedAt          *time.Time   `json:"deleted_at" db:"deleted_at"`   // If the snapshot has been garbage collected
	DiffsCount         *int32       `json:"diffs_count" db:"diffs_count"` // number of diffs in a Snapshot
}

func (s *Snapshot) BranchName() string {
	return fmt.Sprintf("snapshot-%s", s.ID)
}

type SnapshotJSON struct {
	*Snapshot
	CreatedAt jsontime.Time `json:"created_at"`
}

type Action string

func (a Action) String() string {
	return string(a)
}

const (
	ActionViewSync                  Action = "view_sync"
	ActionSyncCompleted             Action = "sync_completed"
	ActionFileUndoPatch             Action = "file_undo_patch"
	ActionFileUndoChange            Action = "file_undo_change"
	ActionFileIgnore                Action = "file_ignore"
	ActionFileRevert                Action = "file_revert"
	ActionChangeLand                Action = "change_land"
	ActionPreChangeLand             Action = "pre_change_land"
	ActionPreCheckoutOtherView      Action = "pre_checkout_other_view"
	ActionPreCheckoutOtherWorkspace Action = "pre_checkout_other_workspace"
	ActionWorkspaceExtract          Action = "workspace_extract"
	ActionChangeReverted            Action = "change_reverted"
	ActionSuggestionApply           Action = "suggestion_apply"
	ActionCITrigger                 Action = "ci_trigger"
)

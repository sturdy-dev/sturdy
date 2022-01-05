package sync

import (
	"mash/pkg/unidiff"
	"time"
)

type Sync struct {
	ID            string `db:"id"`
	UserID        string `db:"user_id"`
	CodebaseID    string `db:"codebase_id"`
	WorkspaceID   string `db:"workspace_id"`
	ViewID        string `db:"view_id"`
	BaseCommit    string `db:"base_commit"`
	OntoCommit    string `db:"onto_commit"`
	CurrentCommit string `db:"current_commit"`

	CreatedAt   *time.Time `db:"created_at"`
	CompletedAt *time.Time `db:"completed_at"`
	AbortedAt   *time.Time `db:"aborted_at"`

	UnsavedCommit *string `db:"unsaved_commit"`
}

type RebaseStatusResponse struct {
	IsRebasing       bool              `json:"is_rebasing"`
	HaveConflicts    bool              `json:"have_conflicts"`
	ConflictingFiles []ConflictingFile `json:"conflicting_files"`
	CanContinue      bool              `json:"can_continue"`

	ProgressCurrent uint `json:"progress_current"`
	ProgressTotal   uint `json:"progress_total"`
}

type ConflictingFile struct {
	Path          string           `json:"path"`
	WorkspaceDiff unidiff.FileDiff `json:"workspace_diff"`
	TrunkDiff     unidiff.FileDiff `json:"trunk_diff"`
}

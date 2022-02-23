package workspaces

import (
	"fmt"
	"time"

	"getsturdy.com/api/pkg/change"
	"getsturdy.com/api/pkg/snapshots"
)

type Workspace struct {
	ID                     string  `db:"id" json:"id"`
	UserID                 string  `db:"user_id" json:"user_id"`
	CodebaseID             string  `db:"codebase_id" json:"codebase_id"`
	Name                   *string `db:"name" json:"name"`
	ReadyForReviewChangeID *string `db:"ready_for_review_change" json:"ready_for_review_change"`
	ApprovedChangeID       *string `db:"approved_change" json:"approved_change"`

	// These are shadowed by the values in WorkspaceWithMetadataJSON
	CreatedAt    *time.Time `db:"created_at" json:"-"`
	LastLandedAt *time.Time `db:"last_landed_at" json:"-"`
	UpdatedAt    *time.Time `db:"updated_at" json:"-"`
	ArchivedAt   *time.Time `db:"archived_at" json:"-"`
	UnarchivedAt *time.Time `db:"unarchived_at" json:"-"`

	DraftDescription string `db:"draft_description" json:"draft_description"`

	// The primary view of this workspace
	ViewID *string `db:"view_id" json:"-"`

	// The last snapshot of this workspace
	// Is used as the "live" diff if the workspace has no view connected
	// and to restore the contents when the workspace is opened again
	LatestSnapshotID *string `db:"latest_snapshot_id" json:"-"`
	DiffsCount       *int32  `db:"diffs_count" json:"-"`

	UpToDateWithTrunk *bool `db:"up_to_date_with_trunk"`

	HeadChangeID       *change.ID `db:"head_change_id" json:"-"`
	HeadChangeComputed bool       `db:"head_change_computed" json:"-"`
}

func (w *Workspace) SetSnapshot(snapshot *snapshots.Snapshot) {
	if snapshot == nil {
		w.LatestSnapshotID = nil
		w.DiffsCount = nil
	} else {
		w.LatestSnapshotID = &snapshot.ID
		w.DiffsCount = snapshot.DiffsCount
	}
}

func (w Workspace) IsArchived() bool {
	if w.UnarchivedAt != nil {
		return false
	}
	return w.ArchivedAt != nil
}

func (w Workspace) NameOrFallback() string {
	if w.Name != nil {
		return *w.Name
	}
	return fmt.Sprintf("Unnamed Workspace %s", w.ID[:8])
}

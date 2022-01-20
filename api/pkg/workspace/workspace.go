package workspace

import (
	"fmt"
	"mash/pkg/author"
	"mash/pkg/jsontime"
	"time"
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

	UpToDateWithTrunk *bool `db:"up_to_date_with_trunk"`

	HeadCommitID   *string `db:"head_commit_id" json:"-"`
	HeadCommitShow bool    `db:"head_commit_show" json:"-"`
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

// WorkspaceWithMetadata contains dynamically generated content
type WorkspaceWithMetadata struct {
	Workspace
	CreatedBy author.Author `json:"created_by"`
}

// time.Time replaced with jsontime.Time to marshal timestamps as unix timestamps
type WorkspaceWithMetadataJSON struct {
	WorkspaceWithMetadata
	CreatedAt    jsontime.Time `json:"created_at"`
	LastLandedAt jsontime.Time `json:"last_landed_at"`
	UpdatedAt    jsontime.Time `json:"updated_at"`
	ArchivedAt   jsontime.Time `json:"archived_at"`
	UnarchivedAt jsontime.Time `json:"unarchived_at"`
}

func ToJSON(in WorkspaceWithMetadata) WorkspaceWithMetadataJSON {
	return WorkspaceWithMetadataJSON{
		WorkspaceWithMetadata: in,
		CreatedAt:             jsontime.FromTimeZeroIfNil(in.CreatedAt),
		LastLandedAt:          jsontime.FromTimeZeroIfNil(in.LastLandedAt),
		UpdatedAt:             jsontime.FromTimeZeroIfNil(in.UpdatedAt),
		ArchivedAt:            jsontime.FromTimeZeroIfNil(in.ArchivedAt),
		UnarchivedAt:          jsontime.FromTimeZeroIfNil(in.UnarchivedAt),
	}
}

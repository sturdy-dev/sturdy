package workspaces

import (
	"html"
	"strings"
	"time"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/users"

	"github.com/microcosm-cc/bluemonday"
)

type Workspace struct {
	ID         string       `db:"id" json:"id"`
	UserID     users.ID     `db:"user_id" json:"user_id"`
	CodebaseID codebases.ID `db:"codebase_id" json:"codebase_id"`
	Name       *string      `db:"name" json:"name"`

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

	HeadChangeID       *changes.ID `db:"head_change_id" json:"-"`
	HeadChangeComputed bool        `db:"head_change_computed" json:"-"`

	// ChangeID is the last change id that was landed from this workspace.
	ChangeID *changes.ID `db:"change_id" json:"-"`
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

var newLiner = strings.NewReplacer(
	"<ul>", "<ul>\n",
	"<ol>", "<ol>\n",
	"</ol>", "</ol>\n",
	"<li><p>", "<li><p>\n* ",
	"<h1>", "\n<h1>\n",
	"<h2>", "\n<h2>\n",
	"<h3>", "\n<h3>\n",
	"<h4>", "\n<h4>\n",
	"<h5>", "\n<h5>\n",
	"<h6>", "\n<h6>\n",
	"<br>", "<br>\n",
	"<p>", "<p>\n",
)

func (w Workspace) NameOrFallback() string {
	replaced := newLiner.Replace(w.DraftDescription)
	sanitizedDescription := strings.TrimLeft(bluemonday.StrictPolicy().Sanitize(replaced), "\n")

	// UnescapeString replaces "&lt;" with "<" etc.
	sanitizedDescription = html.UnescapeString(sanitizedDescription)

	if sanitizedDescription == "" {
		if w.Name != nil {
			return *w.Name
		}
		return "Untitled draft"
	}
	first, _, found := strings.Cut(sanitizedDescription, "\n")
	if found {
		return strings.TrimSpace(first)
	}
	return strings.TrimSpace(sanitizedDescription)
}

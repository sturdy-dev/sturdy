package db

import (
	"context"
	"fmt"
	"strings"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(entity workspaces.Workspace) error {
	_, err := r.db.NamedExec(`INSERT INTO workspaces
		(id, user_id, codebase_id, name, created_at, view_id, latest_snapshot_id, draft_description, diffs_count)
		VALUES
		(:id, :user_id, :codebase_id, :name, :created_at, :view_id, :latest_snapshot_id, :draft_description, :diffs_count)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to insert workspace: %w", err)
	}
	return nil
}

func (r *repo) Get(id string) (*workspaces.Workspace, error) {
	var entity workspaces.Workspace
	err := r.db.Get(&entity, `SELECT id, user_id, codebase_id, name,  created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed, diffs_count, change_id
	FROM workspaces
	WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}
	return &entity, nil
}

func (r *repo) ListByCodebaseIDs(codebaseIDs []codebases.ID, includeArchived bool) ([]*workspaces.Workspace, error) {
	q := `SELECT id, user_id, codebase_id, name, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed, diffs_count, change_id
	FROM workspaces
	WHERE codebase_id IN(?)`

	if !includeArchived {
		q += "  AND archived_at IS NULL"
	}

	query, args, err := sqlx.In(q, codebaseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)

	var views []*workspaces.Workspace
	err = r.db.Select(&views, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to ListByCodebaseIDs: %w", err)
	}
	return views, nil
}

func (r *repo) ListByCodebaseIDsAndUserID(codebaseIDs []codebases.ID, userID string) ([]*workspaces.Workspace, error) {
	query, args, err := sqlx.In(`SELECT id, user_id, codebase_id, name, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, diffs_count, change_id
	FROM workspaces
	WHERE codebase_id IN(?)
	  AND user_id = ?
	  AND archived_at IS NULL`,
		codebaseIDs,
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)

	var views []*workspaces.Workspace
	err = r.db.Select(&views, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to ListByCodebaseIDsAndUserID: %w", err)
	}
	return views, nil
}

func (r *repo) UnsetUpToDateWithTrunkForAllInCodebase(codebaseID codebases.ID) error {
	_, err := r.db.Exec("UPDATE workspaces SET up_to_date_with_trunk = NULL WHERE codebase_id = $1 AND archived_at IS NULL", codebaseID)
	if err != nil {
		return fmt.Errorf("failed to UnsetUpToDateWithTrunkForAllInCodebase: %w", err)
	}
	return nil
}

func (r *repo) GetByViewID(viewID string, includeArchived bool) (*workspaces.Workspace, error) {
	var entity workspaces.Workspace

	q := `SELECT id, user_id, codebase_id, name, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed, diffs_count, change_id
		FROM workspaces
		WHERE view_id=$1`

	if !includeArchived {
		q += " AND archived_at IS NULL"
	}
	err := r.db.Get(&entity, q, viewID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetByViewID: %w", err)
	}
	return &entity, nil
}

func (r *repo) ListByUserID(ctx context.Context, userID users.ID) ([]*workspaces.Workspace, error) {
	var entities []*workspaces.Workspace
	if err := r.db.SelectContext(ctx, &entities, `SELECT 
		id,
		user_id, 
		codebase_id, 
		name, 
		created_at, 
		last_landed_at, 
		archived_at, 
		unarchived_at, 
		updated_at, 
		draft_description, 
		view_id, 
		latest_snapshot_id, 
		up_to_date_with_trunk, 
		head_change_id, 
		head_change_computed, 
		diffs_count, 
		change_id
	FROM workspaces
	WHERE user_id=$1
	AND archived_at IS NULL`, userID); err != nil {
		return nil, fmt.Errorf("failed to ListByUserID: %w", err)
	}
	return entities, nil
}

func (r *repo) GetBySnapshotID(snapshotID string) (*workspaces.Workspace, error) {
	var entity workspaces.Workspace
	if err := r.db.Get(&entity, `
		SELECT 
			id, 
			user_id,
			codebase_id,
			name,
			created_at,
			last_landed_at,
			archived_at,
			unarchived_at,
			updated_at,
			draft_description,
			view_id,
			latest_snapshot_id,
			up_to_date_with_trunk,
			head_change_id,
			head_change_computed,
			diffs_count,
			change_id
		FROM 
			workspaces
		WHERE
			latest_snapshot_id=$1
	`, snapshotID); err != nil {
		return nil, fmt.Errorf("failed to GetBySnapshotID: %w", err)
	}
	return &entity, nil
}

type updateQuery struct {
	raw  strings.Builder
	args map[string]any
}

func newQuery() *updateQuery {
	return &updateQuery{
		raw:  strings.Builder{},
		args: make(map[string]any),
	}
}

func (u *updateQuery) String(workspaceID string) string {
	u.args["workspace_id"] = workspaceID
	return strings.Join([]string{
		`UPDATE workspaces SET`,
		u.raw.String(),
		`WHERE id = :workspace_id`,
	}, " ")
}

func (u *updateQuery) Set(field string, value any) *updateQuery {
	if len(u.args) == 0 {
		u.raw.WriteString(fmt.Sprintf("%s = :%s", field, field))
	} else {
		u.raw.WriteString(fmt.Sprintf(", %s = :%s", field, field))
	}
	u.args[field] = value
	return u
}

func (r *repo) UpdateFields(ctx context.Context, workspaceID string, fields ...UpdateOption) error {
	if len(fields) == 0 {
		return nil
	}

	opts := Options(fields).Parse()
	query := newQuery()

	if opts.updatedAtSet {
		query.Set("updated_at", opts.updatedAt)
	}
	if opts.upToDateWithTrunkSet {
		query.Set("up_to_date_with_trunk", opts.upToDateWithTrunk)
	}
	if opts.headChangeIDSet {
		query.Set("head_change_id", opts.headChangeID)
	}
	if opts.headChangeComputedSet {
		query.Set("head_change_computed", opts.headChangeComputed)
	}
	if opts.latestSnapshotIDSet {
		query.Set("latest_snapshot_id", opts.latestSnapshotID)
	}
	if opts.diffsCountSet {
		query.Set("diffs_count", opts.diffsCount)
	}
	if opts.viewIDSet {
		query.Set("view_id", opts.viewID)
	}
	if opts.lastLandedAtSet {
		query.Set("last_landed_at", opts.lastLandedAt)
	}
	if opts.changeIDSet {
		query.Set("change_id", opts.changeID)
	}
	if opts.draftDescriptionSet {
		query.Set("draft_description", opts.draftDescription)
	}
	if opts.archivedAtSet {
		query.Set("archived_at", opts.archivedAt)
	}
	if opts.unarchivedAtSet {
		query.Set("unarchived_at", opts.unarchivedAt)
	}
	if opts.nameSet {
		query.Set("name", opts.name)
	}
	if opts.userIDSet {
		query.Set("user_id", opts.userID)
	}

	if _, err := r.db.NamedExecContext(ctx, query.String(workspaceID), query.args); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}
	return nil
}

func (r *repo) ListByIDs(ctx context.Context, ids ...string) ([]*workspaces.Workspace, error) {
	var entities []*workspaces.Workspace
	if err := r.db.SelectContext(ctx, &entities, `SELECT 
		id,
		user_id, 
		codebase_id, 
		name, 
		created_at, 
		last_landed_at, 
		archived_at, 
		unarchived_at, 
		updated_at, 
		draft_description, 
		view_id, 
		latest_snapshot_id, 
		up_to_date_with_trunk, 
		head_change_id, 
		head_change_computed, 
		diffs_count, 
		change_id
	FROM workspaces
	WHERE id IN (:ids)`, map[string]interface{}{"ids": ids}); err != nil {
		return nil, fmt.Errorf("failed to ListByIDs: %w", err)
	}
	return entities, nil
}

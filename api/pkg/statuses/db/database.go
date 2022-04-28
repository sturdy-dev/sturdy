package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/statuses"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, status *statuses.Status) error {
	if _, err := r.db.NamedExecContext(ctx, `
		INSERT INTO statuses (
			id,
			commit_id,
			codebase_id,
			title,
			description,
			type,
			timestamp,
			details_url
		) VALUES (
			:id,
			:commit_id,
			:codebase_id,
			:title,
			:description,
			:type,
			:timestamp,
			:details_url
		)
	`, status); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (r *repository) Get(ctx context.Context, id string) (*statuses.Status, error) {
	var s statuses.Status
	if err := r.db.GetContext(ctx, &s, `
		SELECT
			id,
			commit_id,
			codebase_id,
			title,
			description,
			type,
			timestamp,
			details_url
		FROM
			statuses
		WHERE
			id = $1 
	`, id); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return &s, nil
}

func (r *repository) ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*statuses.Status, error) {
	var ss []*statuses.Status
	if err := r.db.SelectContext(ctx, &ss, `
		WITH latest AS (
            SELECT
                statuses.commit_id, 
				statuses.codebase_id, 
				statuses.title,
				MAX(statuses.timestamp) AS timestamp
            FROM
                statuses
                    JOIN snapshots ON snapshots.commit_id = statuses.commit_id
            WHERE
                snapshots.workspace_id = $1
            GROUP BY
                statuses.commit_id, statuses.codebase_id, statuses.title
        )
        SELECT
			statuses.id,
			statuses.commit_id,
			statuses.codebase_id,
			statuses.title,
			statuses.description,
			statuses.type,
			statuses.timestamp,
			statuses.details_url
        FROM
            statuses 
				JOIN latest ON
					statuses.commit_id       = latest.commit_id
					AND statuses.codebase_id = latest.codebase_id
					AND statuses.title       = latest.title
					AND statuses.timestamp   = latest.timestamp
	`, workspaceID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return ss, nil
}

// ListByCodebaseIDAndCommitID returns a list of latest statuses for commit_id grouped by title.
func (r *repository) ListByCodebaseIDAndCommitID(ctx context.Context, codebaseID codebases.ID, commitID string) ([]*statuses.Status, error) {
	var ss []*statuses.Status
	if err := r.db.SelectContext(ctx, &ss, `
		WITH latest AS (
			SELECT
				commit_id, codebase_id, title, MAX(timestamp) AS timestamp
			FROM
				statuses
			WHERE
				commit_id = $1
				AND codebase_id = $2
			GROUP BY
				commit_id, codebase_id, title
		)

		SELECT
				statuses.id,
				statuses.commit_id,
				statuses.codebase_id,
				statuses.title,
				statuses.description,
				statuses.type,
				statuses.timestamp,
				statuses.details_url
		FROM
			statuses JOIN latest ON
				statuses.commit_id       = latest.commit_id
				AND statuses.codebase_id = latest.codebase_id
				AND statuses.title       = latest.title
				AND statuses.timestamp   = latest.timestamp
		`, commitID, codebaseID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return ss, nil
}

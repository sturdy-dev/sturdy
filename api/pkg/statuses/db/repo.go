package db

import (
	"context"
	"fmt"

	"mash/pkg/statuses"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, status *statuses.Status) error {
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

func (r *Repository) Get(ctx context.Context, id string) (*statuses.Status, error) {
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

// ListByCodebaseIDAndCommitID returns a list of latest statuses for commit_id grouped by title.
func (r *Repository) ListByCodebaseIDAndCommitID(ctx context.Context, codebaseID string, commitID string) ([]*statuses.Status, error) {
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

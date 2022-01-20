package db

import (
	"fmt"
	"mash/pkg/mutagen"

	"github.com/jmoiron/sqlx"
)

type ViewStatusRepository interface {
	Create(status mutagen.ViewStatus) error
	Update(status *mutagen.ViewStatus) error
	GetByViewID(viewID string) (*mutagen.ViewStatus, error)
}

func NewRepository(db *sqlx.DB) ViewStatusRepository {
	return &repo{db: db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) Create(status mutagen.ViewStatus) error {
	_, err := r.db.NamedExec(`INSERT INTO view_status (
				id,
				state,
				staging_status_path,
				staging_status_received,
				staging_status_total,
				sturdy_version,
				last_error)
			VALUES (
				:id,
				:state,
				:staging_status_path,
				:staging_status_received,
				:staging_status_total,
				:sturdy_version,
				:last_error)`, &status)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) Update(status *mutagen.ViewStatus) error {
	_, err := r.db.NamedExec(`UPDATE view_status 
				SET state = :state,
				staging_status_path = :staging_status_path,
				staging_status_received = :staging_status_received,
				staging_status_total = :staging_status_total,
				sturdy_version = :sturdy_version,
				last_error = :last_error,
				updated_at = :updated_at
				WHERE id = :id`, status)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}

func (r *repo) GetByViewID(viewID string) (*mutagen.ViewStatus, error) {
	var res mutagen.ViewStatus
	err := r.db.Get(&res, "SELECT * FROM view_status WHERE id=$1", viewID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetByViewID: %w", err)
	}
	return &res, nil
}

package db

import (
	"context"
	"fmt"

	"mash/pkg/notification"

	"github.com/jmoiron/sqlx"
)

type PreferenceRepository struct {
	db *sqlx.DB
}

func NewPeferenceRepository(db *sqlx.DB) *PreferenceRepository {
	return &PreferenceRepository{
		db: db,
	}
}

func (repo *PreferenceRepository) Upsert(ctx context.Context, preference *notification.Preference) error {
	if _, err := repo.db.NamedExecContext(ctx, `INSERT INTO notification_preferences
		(user_id, type, channel, enabled)
		VALUES
		(:user_id, :type, :channel, :enabled)
		ON CONFLICT ON CONSTRAINT user_id_type_channel_uq_ix
			DO UPDATE SET enabled = :enabled
		`, preference); err != nil {
		return fmt.Errorf("failed to upsert: %w", err)
	}
	return nil
}

func (repo *PreferenceRepository) ListByUserID(ctx context.Context, userID string) ([]*notification.Preference, error) {
	var res []*notification.Preference
	if err := repo.db.SelectContext(ctx, &res, `SELECT user_id, type, channel, enabled
		FROM notification_preferences
		WHERE user_id = $1`, userID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return res, nil
}

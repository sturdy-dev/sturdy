package db

import (
	"fmt"
	"getsturdy.com/api/pkg/notification"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(id string) (notification.Notification, error)
	Create(notification.Notification) error
	Update(notification.Notification) error
	// todo: use id based pagination instead of offset
	ListByUser(userID string, limit, offset int) ([]notification.Notification, error)
	ListByUserAndIds(userID string, ids []string) ([]notification.Notification, error)
	ArchiveByUserAndIds(userID string, ids []string) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Get(id string) (notification.Notification, error) {
	var res notification.Notification
	err := r.db.Get(&res, `SELECT id, codebase_id, user_id, type, reference_id, created_at, archived_at
		FROM notifications
		WHERE id = $1`, id)
	if err != nil {
		return notification.Notification{}, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}

func (r *repo) Create(notification notification.Notification) error {
	_, err := r.db.NamedExec(`INSERT INTO notifications (id, codebase_id, user_id, type, reference_id, created_at, archived_at)
		VALUES (:id, :codebase_id, :user_id, :type, :reference_id, :created_at, :archived_at)`, &notification)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) Update(notification notification.Notification) error {
	_, err := r.db.NamedExec(`UPDATE notifications
    	SET archived_at = :archived_at
    	WHERE id = :id`, &notification)
	if err != nil {
		return fmt.Errorf("failed to update change: %w", err)
	}
	return nil
}

func (r *repo) ListByUser(userID string, limit, offset int) ([]notification.Notification, error) {
	var res []notification.Notification
	err := r.db.Select(&res, `SELECT id, codebase_id, user_id, type, reference_id, created_at, archived_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}

func (r *repo) ListByUserAndIds(userID string, ids []string) ([]notification.Notification, error) {
	query, args, err := sqlx.In(`SELECT id, codebase_id, user_id, type, reference_id, created_at, archived_at
	FROM notifications
	WHERE user_id = ?
	  AND id IN(?)`,
		userID,
		ids)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)
	var entities []notification.Notification
	err = r.db.Select(&entities, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *repo) ArchiveByUserAndIds(userID string, ids []string) error {
	query, args, err := sqlx.In(`UPDATE notifications
	SET archived_at = NOW()
	WHERE user_id = ?
	  AND id IN(?)`,
		userID,
		ids)
	if err != nil {
		return fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}
	return nil
}

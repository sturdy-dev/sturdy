package db

import (
	"fmt"
	"mash/pkg/newsletter"

	"github.com/jmoiron/sqlx"
)

type NotificationSettingsRepository interface {
	GetByUser(userID string) (*newsletter.NotificationSettings, error)
	Insert(newsletter.NotificationSettings) error
	Update(*newsletter.NotificationSettings) error
}

type repo struct {
	db *sqlx.DB
}

func NewNotificationSettingsRepository(db *sqlx.DB) NotificationSettingsRepository {
	return &repo{db}
}

func (r *repo) Insert(settings newsletter.NotificationSettings) error {
	_, err := r.db.NamedExec(`INSERT INTO notification_settings (user_id, receive_newsletter)
		VALUES (:user_id, :receive_newsletter)`, settings)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) Update(settings *newsletter.NotificationSettings) error {
	_, err := r.db.NamedExec(`UPDATE notification_settings
    	SET receive_newsletter = :receive_newsletter
	 	WHERE user_id = :user_id`, settings)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) GetByUser(userID string) (*newsletter.NotificationSettings, error) {
	var res newsletter.NotificationSettings
	err := r.db.Get(&res, "SELECT user_id, receive_newsletter FROM notification_settings WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

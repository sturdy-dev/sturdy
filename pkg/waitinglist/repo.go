package waitinglist

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type WaitingListRepo interface {
	Insert(email string) error
	ToSendInvitesTo() ([]WaitingListEntry, error)
	MarkEmailAsInvited(email string) error
}

type waitingListRepo struct {
	db *sqlx.DB
}

func NewWaitingListRepo(db *sqlx.DB) WaitingListRepo {
	return &waitingListRepo{db}
}

func (r *waitingListRepo) Insert(email string) error {
	_, err := r.db.Exec(`INSERT INTO waitinglist (email, created_at)
		VALUES ($1, NOW())`, email)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

type WaitingListEntry struct {
	ID    int    `db:"id"`
	Email string `db:"email"`
}

func (r *waitingListRepo) ToSendInvitesTo() ([]WaitingListEntry, error) {
	var res []WaitingListEntry
	err := r.db.Select(&res, `SELECT id, email
		FROM waitinglist
		WHERE should_send_email = true
		  AND invited_at IS NULL
		  AND ignored IS NULL
		  LIMIT 10`)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *waitingListRepo) MarkEmailAsInvited(email string) error {
	_, err := r.db.Exec(`UPDATE waitinglist SET invited_at = NOW() WHERE email = $1`, email)
	if err != nil {
		return err
	}
	return nil
}

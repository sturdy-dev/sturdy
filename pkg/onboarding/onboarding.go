package onboarding

import "time"

type Step struct {
	ID        string    `db:"step_id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

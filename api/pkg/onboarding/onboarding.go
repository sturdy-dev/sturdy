package onboarding

import (
	"time"

	"getsturdy.com/api/pkg/users"
)

type Step struct {
	ID        string    `db:"step_id"`
	UserID    users.ID  `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

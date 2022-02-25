package pki

import (
	"time"

	"getsturdy.com/api/pkg/users"
)

type UserPublicKey struct {
	// PublicKey is the Primary Key, and is globally unique
	PublicKey string     `json:"public_key" db:"public_key"`
	UserID    users.ID   `json:"user_id" db:"user_id"`
	AddedAt   time.Time  `json:"added_at" db:"added_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
}

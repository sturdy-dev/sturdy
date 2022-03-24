package crypto

import (
	"time"

	"getsturdy.com/api/pkg/users"
)

type KeyPair struct {
	ID         KeyPairID  `db:"ID"`
	PublicKey  PublicKey  `db:"public_key"`
	PrivateKey PrivateKey `db:"private_key"`
	CreatedAt  time.Time  `db:"created_at"`
	CreatedBy  users.ID   `db:"created_by"`
	LastUsedAt time.Time  `db:"last_used_at"`
}

type PublicKey string

type PrivateKey string

type KeyPairID string

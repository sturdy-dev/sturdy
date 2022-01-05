package servicetokens

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	ID         string     `db:"id"`
	CodebaseID string     `db:"codebase_id"`
	Hash       []byte     `db:"hash"`
	Name       string     `db:"name"`
	CreatedAt  time.Time  `db:"created_at"`
	LastUsedAt *time.Time `db:"last_used_at"`
}

func (t *Token) Verify(c string) error {
	return bcrypt.CompareHashAndPassword(t.Hash, []byte(c))
}

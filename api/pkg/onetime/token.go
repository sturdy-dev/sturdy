package onetime

import (
	"fmt"
	"math/rand"
	"time"

	"getsturdy.com/api/pkg/users"
)

// Token is a short-lived, one-time use token.
type Token struct {
	Key       string    `db:"key"`
	UserID    users.ID  `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	// Number of times the token has been used.
	Clicks int `db:"clicks"`
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func New(userID users.ID) *Token {
	return &Token{
		Key:       fmt.Sprintf("%s%s", randSeq(3), randSeq(3)),
		UserID:    userID,
		CreatedAt: time.Now(),
	}
}

const expireAfter = time.Minute * 10

func (t *Token) IsExpired() bool {
	return time.Since(t.CreatedAt) > expireAfter
}

func (t *Token) IsReused() bool {
	return t.Clicks > 0
}

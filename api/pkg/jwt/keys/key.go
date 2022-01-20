package keys

import (
	"fmt"

	"github.com/google/uuid"
)

var ErrKeyIsEmpty = fmt.Errorf("key is empty")

type Key struct {
	ID        string `db:"id"`
	PublicDER []byte `db:"public_der"`
}

// New creates a new key with a public der payload.
func New(publicDER []byte) (*Key, error) {
	if len(publicDER) == 0 {
		return nil, ErrKeyIsEmpty
	}

	return &Key{
		ID:        uuid.New().String(),
		PublicDER: publicDER,
	}, nil
}

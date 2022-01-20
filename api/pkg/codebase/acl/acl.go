package acl

import (
	"time"
)

type ID string

type ACL struct {
	ID         ID        `json:"id" db:"id"`
	CodebaseID string    `json:"codebase_id" db:"codebase_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	RawPolicy  string    `json:"policy" db:"policy"`

	// Policy contains a policy parsed from RawPolicy
	// Note that changes from this field won't be persisted in the database
	Policy Policy `json:"-" db:"-"`
}

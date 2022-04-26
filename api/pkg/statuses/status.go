package statuses

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
)

type Type string

const (
	TypeUndefined Type = ""
	TypePending   Type = "pending"
	TypeHealthy   Type = "healthy"
	TypeFailing   Type = "failing"
)

var ValidType = map[Type]bool{
	TypePending: true,
	TypeHealthy: true,
	TypeFailing: true,
}

type Status struct {
	ID          string       `db:"id" json:"id"`
	CommitSHA   string       `db:"commit_id" json:"commit_id"`
	CodebaseID  codebases.ID `db:"codebase_id" json:"codebase_id"`
	Type        Type         `db:"type" json:"type"`
	Title       string       `db:"title" json:"title"`
	DetailsURL  *string      `db:"details_url" json:"details_url"`
	Description *string      `db:"description" json:"description,omitempty"`
	Timestamp   time.Time    `db:"timestamp" json:"timestamp"`
}

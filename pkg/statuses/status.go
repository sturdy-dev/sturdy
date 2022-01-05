package statuses

import "time"

type Type string

const (
	TypeUndefined Type = ""
	TypePending   Type = "pending"
	TypeHealty    Type = "healthy"
	TypeFailing   Type = "failing"
)

var ValidType = map[Type]bool{
	TypePending: true,
	TypeHealty:  true,
	TypeFailing: true,
}

type Status struct {
	ID          string    `db:"id" json:"id"`
	CommitID    string    `db:"commit_id" json:"commit_id"`
	CodebaseID  string    `db:"codebase_id" json:"codebase_id"`
	Type        Type      `db:"type" json:"type"`
	Title       string    `db:"title" json:"title"`
	DetailsURL  *string   `db:"details_url" json:"details_url"`
	Description *string   `db:"description" json:"description,omitempty"`
	Timestamp   time.Time `db:"timestamp" json:"timestamp"`
}

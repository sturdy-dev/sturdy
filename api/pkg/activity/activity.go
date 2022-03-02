package activity

import (
	"time"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/users"
)

type Activity struct {
	ID           string    `db:"id"`
	UserID       users.ID  `db:"user_id"`
	WorkspaceID  string    `db:"workspace_id"`
	CreatedAt    time.Time `db:"created_at"`
	ActivityType Type      `db:"activity_type"`
	Reference    string    `db:"reference"`

	// If change_id is set, this is a change activity
	ChangeID *changes.ID `db:"change_id"`
}

type Type string

const (
	TypeComment         Type = "comment"          // Reference is a Comment ID
	TypeCreatedChange   Type = "created_change"   // Reference is a Change ID
	TypeRequestedReview Type = "requested_review" // Reference is a Review ID
	TypeReviewed        Type = "reviewed"         // Reference is a Review ID
)

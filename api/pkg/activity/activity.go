package activity

import (
	"time"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/users"
)

type Activity struct {
	ID           string    `db:"id"`
	UserID       users.ID  `db:"user_id"`
	CreatedAt    time.Time `db:"created_at"`
	ActivityType Type      `db:"activity_type"`
	Reference    string    `db:"reference"`

	// if WorkspaceID is set, then the activity is for a workspace
	WorkspaceID *string `db:"workspace_id"`
	// if ChangeID is set, then the activity is for a change
	ChangeID *changes.ID `db:"change_id"`
}

type Type string

const (
	TypeComment         Type = "comment"          // Reference is a Comment ID
	TypeCreatedChange   Type = "created_change"   // Reference is a Change ID
	TypeRequestedReview Type = "requested_review" // Reference is a Review ID
	TypeReviewed        Type = "reviewed"         // Reference is a Review ID
)

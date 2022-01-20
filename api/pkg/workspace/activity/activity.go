package activity

import "time"

type WorkspaceActivity struct {
	ID           string                `db:"id"`
	UserID       string                `db:"user_id"`
	WorkspaceID  string                `db:"workspace_id"`
	CreatedAt    time.Time             `db:"created_at"`
	ActivityType WorkspaceActivityType `db:"activity_type"`
	Reference    string                `db:"reference"`
}

type WorkspaceActivityType string

const (
	WorkspaceActivityTypeComment         WorkspaceActivityType = "comment"          // Reference is a Comment ID
	WorkspaceActivityTypeCreatedChange   WorkspaceActivityType = "created_change"   // Reference is a Change ID
	WorkspaceActivityTypeRequestedReview WorkspaceActivityType = "requested_review" // Reference is a Review ID
	WorkspaceActivityTypeReviewed        WorkspaceActivityType = "reviewed"         // Reference is a Review ID
)

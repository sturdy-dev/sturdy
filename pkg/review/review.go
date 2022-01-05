package review

import "time"

type Review struct {
	ID          string      `db:"id"`
	UserID      string      `db:"user_id"`
	CodebaseID  string      `db:"codebase_id"`
	WorkspaceID string      `db:"workspace_id"`
	Grade       ReviewGrade `db:"grade"`
	CreatedAt   time.Time   `db:"created_at"`
	DismissedAt *time.Time  `db:"dismissed_at"`
	IsReplaced  bool        `db:"is_replaced"` // Is false for new reviews.
	RequestedBy *string     `db:"requested_by"`
}

type ReviewGrade string

const (
	ReviewGradeApprove   ReviewGrade = "Approve"
	ReviewGradeReject    ReviewGrade = "Reject"
	ReviewGradeRequested ReviewGrade = "Requested"
)

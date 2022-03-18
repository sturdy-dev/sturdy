package review

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/users"
)

type Review struct {
	ID          string       `db:"id"`
	UserID      users.ID     `db:"user_id"`
	CodebaseID  codebases.ID `db:"codebase_id"`
	WorkspaceID string       `db:"workspace_id"`
	Grade       ReviewGrade  `db:"grade"`
	CreatedAt   time.Time    `db:"created_at"`
	DismissedAt *time.Time   `db:"dismissed_at"`
	IsReplaced  bool         `db:"is_replaced"` // Is false for new reviews.
	RequestedBy *users.ID    `db:"requested_by"`
}

type ReviewGrade string

const (
	ReviewGradeApprove   ReviewGrade = "Approve"
	ReviewGradeReject    ReviewGrade = "Reject"
	ReviewGradeRequested ReviewGrade = "Requested"
)

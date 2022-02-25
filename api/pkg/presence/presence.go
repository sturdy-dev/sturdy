package presence

import (
	"time"

	"getsturdy.com/api/pkg/users"
)

type Presence struct {
	ID           string    `db:"id"`
	UserID       users.ID  `db:"user_id"`
	WorkspaceID  string    `db:"workspace_id"`
	LastActiveAt time.Time `db:"last_active_at"`
	State        State     `db:"state"`
}

type State string

const (
	StateIdle    State = "idle"
	StateViewing State = "viewing"
	StateCoding  State = "coding"
)

var StatePriority = map[State]uint8{
	StateIdle:    10, // Idle and Viewing are both set from the web, and have the same priority
	StateViewing: 10,
	StateCoding:  20,
}

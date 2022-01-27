package licenses

import (
	"time"
)

type ID string

type Status string

const (
	StatusValid   Status = "valid"
	StatusInvalid Status = "invalid"
)

type Level string

const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

type Type string

const (
	TypeNotification Type = "notification"
	TypeBanner       Type = "banner"
	TypeFullscreen   Type = "fullscreen"
)

type Message struct {
	Type  Type   `json:"type"`
	Level Level  `json:"level"`
	Text  string `json:"text"`
}

type License struct {
	ID             ID        `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organizationId"`
	Key            string    `db:"key" json:"key"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	ExpiresAt      time.Time `db:"expires_at" json:"expiresAt"`

	Status   Status     `db:"-" json:"status"`
	Messages []*Message `db:"-" json:"messages"`
}

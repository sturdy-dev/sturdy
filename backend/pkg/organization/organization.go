package organization

import (
	"time"
)

type ShortOrganizationID string

type Organization struct {
	ID        string              `db:"id"`
	ShortID   ShortOrganizationID `db:"short_id"` // Used in web slugs
	Name      string              `db:"name"`
	CreatedAt time.Time           `db:"created_at"`
	CreatedBy string              `db:"created_by"`
	DeletedAt *time.Time          `db:"deleted_at"`
	DeletedBy *string             `db:"deleted_by"`
}

type Member struct {
	ID             string     `db:"id"`
	UserID         string     `db:"user_id"`
	OrganizationID string     `db:"organization_id"`
	CreatedAt      time.Time  `db:"created_at"`
	CreatedBy      string     `db:"created_by"`
	DeletedAt      *time.Time `db:"deleted_at"`
	DeletedBy      *string    `db:"deleted_by"`
}

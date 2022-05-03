package organization

import (
	"fmt"
	"time"

	"getsturdy.com/api/pkg/users"
)

type ShortOrganizationID string

type Organization struct {
	ID        string              `db:"id"`
	ShortID   ShortOrganizationID `db:"short_id"` // Used in web slugs
	Name      string              `db:"name"`
	CreatedAt time.Time           `db:"created_at"`
	CreatedBy users.ID            `db:"created_by"`
	DeletedAt *time.Time          `db:"deleted_at"`
	DeletedBy *string             `db:"deleted_by"`
}

func (o Organization) Slug() string {
	return fmt.Sprintf("%s-%s", o.Name, o.ShortID)
}

type Member struct {
	ID             string     `db:"id"`
	UserID         users.ID   `db:"user_id"`
	OrganizationID string     `db:"organization_id"`
	CreatedAt      time.Time  `db:"created_at"`
	CreatedBy      users.ID   `db:"created_by"`
	DeletedAt      *time.Time `db:"deleted_at"`
	DeletedBy      *users.ID  `db:"deleted_by"`
}

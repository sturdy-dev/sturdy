package inmemory

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/organization"
	db_organization "getsturdy.com/api/pkg/organization/db"
)

type inMemoryOrganizationRepository struct {
	orgs []organization.Organization
}

func NewInMemoryOrganizationRepo() db_organization.Repository {
	return &inMemoryOrganizationRepository{orgs: make([]organization.Organization, 0)}
}

func (r *inMemoryOrganizationRepository) Get(ctx context.Context, id string) (*organization.Organization, error) {
	for _, org := range r.orgs {
		if org.ID == id {
			return &org, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryOrganizationRepository) GetByShortID(ctx context.Context, shortID organization.ShortOrganizationID) (*organization.Organization, error) {
	for _, org := range r.orgs {
		if org.ShortID == shortID {
			return &org, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryOrganizationRepository) GetFirst(ctx context.Context) (*organization.Organization, error) {
	if len(r.orgs) > 0 {
		return &r.orgs[0], nil
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryOrganizationRepository) Create(ctx context.Context, org organization.Organization) error {
	r.orgs = append(r.orgs, org)
	return nil
}

func (r *inMemoryOrganizationRepository) Update(ctx context.Context, org *organization.Organization) error {
	for k, v := range r.orgs {
		if v.ID == org.ID {
			r.orgs[k] = *org
		}
	}
	return nil
}

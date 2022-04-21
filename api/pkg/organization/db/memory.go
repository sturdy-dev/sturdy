package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/users"
)

type inMemoryOrganizationMemberRepository struct {
	users []organization.Member
}

func NewInMemoryOrganizationMemberRepository() MemberRepository {
	return &inMemoryOrganizationMemberRepository{users: make([]organization.Member, 0)}
}

func (r *inMemoryOrganizationMemberRepository) GetByID(_ context.Context, id string) (*organization.Member, error) {
	for _, u := range r.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryOrganizationMemberRepository) GetByUserIDAndOrganizationID(ctx context.Context, userID users.ID, organizationID string) (*organization.Member, error) {
	for _, u := range r.users {
		if u.UserID == userID && u.OrganizationID == organizationID {
			return &u, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryOrganizationMemberRepository) ListByOrganizationID(ctx context.Context, id string) ([]*organization.Member, error) {
	var res []*organization.Member
	for _, u := range r.users {
		if u.OrganizationID == id {
			u2 := u
			res = append(res, &u2)
		}
	}
	return res, nil
}

func (r *inMemoryOrganizationMemberRepository) ListByUserID(ctx context.Context, id users.ID) ([]*organization.Member, error) {
	var res []*organization.Member
	for _, u := range r.users {
		if u.UserID == id {
			u2 := u
			res = append(res, &u2)
		}
	}
	return res, nil
}

func (r *inMemoryOrganizationMemberRepository) Create(ctx context.Context, org *organization.Member) error {
	r.users = append(r.users, *org)
	return nil
}

func (r *inMemoryOrganizationMemberRepository) Update(ctx context.Context, org *organization.Member) error {
	for k, v := range r.users {
		if v.ID == org.ID {
			r.users[k] = *org
		}
	}
	return nil
}

type inMemoryOrganizationRepository struct {
	orgs []organization.Organization
}

func NewInMemoryOrganizationRepo() Repository {
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

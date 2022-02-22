package inmemory

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/organization"
	db_organization "getsturdy.com/api/pkg/organization/db"
)

type inMemoryOrganizationMemberRepository struct {
	users []organization.Member
}

func NewInMemoryOrganizationMemberRepository() db_organization.MemberRepository {
	return &inMemoryOrganizationMemberRepository{users: make([]organization.Member, 0)}
}

func (r *inMemoryOrganizationMemberRepository) GetByUserIDAndOrganizationID(ctx context.Context, userID, organizationID string) (*organization.Member, error) {
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

func (r *inMemoryOrganizationMemberRepository) ListByUserID(ctx context.Context, id string) ([]*organization.Member, error) {
	var res []*organization.Member
	for _, u := range r.users {
		if u.UserID == id {
			u2 := u
			res = append(res, &u2)
		}
	}
	return res, nil
}

func (r *inMemoryOrganizationMemberRepository) Create(ctx context.Context, org organization.Member) error {
	r.users = append(r.users, org)
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

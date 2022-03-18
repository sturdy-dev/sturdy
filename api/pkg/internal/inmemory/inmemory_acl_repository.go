package inmemory

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebases/acl"
	db_acl "getsturdy.com/api/pkg/codebases/acl/db"
)

type inMemoryAclRepo struct {
	acls []acl.ACL
}

func NewInMemoryAclRepo() db_acl.ACLRepository {
	return &inMemoryAclRepo{
		acls: make([]acl.ACL, 0),
	}
}

func (r *inMemoryAclRepo) Create(_ context.Context, a acl.ACL) error {
	r.acls = append(r.acls, a)
	return nil

}

func (r *inMemoryAclRepo) Update(_ context.Context, a acl.ACL) error {
	for k, v := range r.acls {
		if v.ID == a.ID {
			r.acls[k] = a
		}
	}
	return nil
}

func (r *inMemoryAclRepo) GetByCodebaseID(_ context.Context, codebaseID string) (acl.ACL, error) {
	for _, v := range r.acls {
		if v.CodebaseID == codebaseID {
			return v, nil
		}
	}
	return acl.ACL{}, sql.ErrNoRows
}

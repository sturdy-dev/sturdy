package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/codebase/acl"
	db_acl "getsturdy.com/api/pkg/codebase/acl/db"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"

	"github.com/google/uuid"
	"github.com/tailscale/hujson"
)

type Provider struct {
	aclDB          db_acl.ACLRepository
	usersDB        db_user.Repository
	codebaseUserDB db_codebase.CodebaseUserRepository
}

func New(
	aclRepo db_acl.ACLRepository,
	codebaseUserDB db_codebase.CodebaseUserRepository,
	usersDB db_user.Repository,
) *Provider {
	return &Provider{
		aclDB:          aclRepo,
		codebaseUserDB: codebaseUserDB,
		usersDB:        usersDB,
	}
}

func (p *Provider) GetByCodebaseID(ctx context.Context, codebaseID string) (acl.ACL, error) {
	entity, err := p.aclDB.GetByCodebaseID(ctx, codebaseID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		entity, err = p.createDefaultPolicy(ctx, codebaseID)
		if err != nil {
			return acl.ACL{}, err
		}
	default:
		return acl.ACL{}, err
	}

	if err := hujson.Unmarshal([]byte(entity.RawPolicy), &entity.Policy); err != nil {
		return acl.ACL{}, fmt.Errorf("failed to unmarshal policy: %w", err)
	}

	return entity, nil
}

func (p *Provider) createDefaultPolicy(ctx context.Context, codebaseID string) (acl.ACL, error) {
	emails, err := p.getUserEmailsForCodebase(ctx, codebaseID)
	if err != nil {
		return acl.ACL{}, err
	}

	a := acl.ACL{
		ID:         acl.ID(uuid.New().String()),
		CodebaseID: codebaseID,
		CreatedAt:  time.Now().UTC(),
	}

	policy, err := defaultACLFor(string(a.ID), emails)
	if err != nil {
		return acl.ACL{}, fmt.Errorf("failed to generage default policy: %w", err)
	}

	a.RawPolicy = policy

	if err := p.aclDB.Create(ctx, a); err != nil {
		return acl.ACL{}, fmt.Errorf("failed to create default policy: %w", err)
	}

	return a, nil
}

func (p *Provider) getUserEmailsForCodebase(ctx context.Context, codebaseID string) ([]string, error) {
	uu, err := p.codebaseUserDB.GetByCodebase(codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query codebase users: %w", err)
	}

	userIDs := make([]users.ID, 0, len(uu))
	for _, u := range uu {
		userIDs = append(userIDs, u.UserID)
	}

	users, err := p.usersDB.GetByIDs(ctx, userIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	emails := make([]string, 0, len(users))
	for _, user := range users {
		emails = append(emails, user.Email)
	}

	return emails, nil
}

func (p *Provider) Update(ctx context.Context, a acl.ACL) error {
	return p.aclDB.Update(ctx, a)
}

package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"mash/pkg/codebase"
	"mash/pkg/codebase/acl"
	"mash/pkg/user"

	"github.com/google/uuid"
	"github.com/tailscale/hujson"
)

type codebaseUserRepository interface {
	GetByCodebase(codebaseID string) ([]*codebase.CodebaseUser, error)
}

type aclRepository interface {
	Create(context.Context, acl.ACL) error
	Update(context.Context, acl.ACL) error
	GetByCodebaseID(ctx context.Context, codebaseID string) (acl.ACL, error)
}

type userRepository interface {
	GetByIDs(ctx context.Context, id ...string) ([]*user.User, error)
}

type Provider struct {
	aclDB          aclRepository
	usersDB        userRepository
	codebaseUserDB codebaseUserRepository
}

func New(
	aclRepo aclRepository,
	codebaseUserDB codebaseUserRepository,
	usersDB userRepository,
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

	userIDs := make([]string, 0, len(uu))
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

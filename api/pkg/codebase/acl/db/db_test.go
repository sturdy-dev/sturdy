package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"getsturdy.com/api/pkg/codebase/acl"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/internal/sturdytest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if os.Getenv("E2E_TEST") != "" {
		os.Exit(m.Run())
	}
	fmt.Println("Skipping, set E2E_TEST to run.")
}

func Test_Create(t *testing.T) {
	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	assert.NoError(t, err)

	repo := NewACLRepository(d)
	ctx := context.Background()

	entity := acl.ACL{
		ID:         acl.ID(uuid.New().String()),
		CodebaseID: uuid.New().String(),
		CreatedAt:  time.Now(),
		Policy:     acl.Policy{},
	}

	assert.NoError(t, repo.Create(ctx, entity))
}

func Test_Create_twice_fails(t *testing.T) {
	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	assert.NoError(t, err)

	repo := NewACLRepository(d)
	ctx := context.Background()

	entity := acl.ACL{
		ID:         acl.ID(uuid.New().String()),
		CodebaseID: uuid.New().String(),
		CreatedAt:  time.Now(),
		Policy:     acl.Policy{},
	}

	assert.NoError(t, repo.Create(ctx, entity))
	assert.Error(t, repo.Create(ctx, entity))
}

func Test_GetByCodebaseID_not_found(t *testing.T) {
	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	assert.NoError(t, err)

	repo := NewACLRepository(d)
	ctx := context.Background()

	fromDB, err := repo.GetByCodebaseID(ctx, "unknown")
	assert.Equal(t, acl.ACL{}, fromDB)
	assert.Error(t, err)
}

func Test_GetByCodebaseID(t *testing.T) {
	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	assert.NoError(t, err)

	repo := NewACLRepository(d)
	ctx := context.Background()

	entity := acl.ACL{
		ID:         acl.ID(uuid.New().String()),
		CodebaseID: uuid.New().String(),
		CreatedAt:  time.Now(),
		Policy:     acl.Policy{},
	}

	encoded, err := json.Marshal(entity.Policy)
	assert.NoError(t, err)
	entity.RawPolicy = string(encoded)

	assert.NoError(t, repo.Create(ctx, entity))

	fromDB, err := repo.GetByCodebaseID(ctx, entity.CodebaseID)
	assert.NoError(t, err)

	assert.Equal(t, entity.ID, fromDB.ID)
	assert.Equal(t, entity.CodebaseID, fromDB.CodebaseID)
	assert.True(t, entity.CreatedAt.Equal(fromDB.CreatedAt))
	assert.Equal(t, entity.Policy, fromDB.Policy)
}

func Test_Update(t *testing.T) {
	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	assert.NoError(t, err)

	repo := NewACLRepository(d)
	ctx := context.Background()

	entity := acl.ACL{
		ID:         acl.ID(uuid.New().String()),
		CodebaseID: uuid.New().String(),
		CreatedAt:  time.Now(),
		Policy:     acl.Policy{},
	}

	assert.NoError(t, repo.Create(ctx, entity))

	entity.Policy.Groups = append(entity.Policy.Groups, &acl.Group{
		ID: "admins",
		Members: []*acl.Identifier{
			{Type: acl.Users, Pattern: "1"},
		},
	})

	entity.Policy.Rules = append(entity.Policy.Rules, &acl.Rule{
		ID:     "rule 1",
		Action: acl.ActionWrite,
		Principals: []*acl.Identifier{
			{Type: acl.Users, Pattern: "123"},
		},
		Resources: []*acl.Identifier{
			{Type: acl.Codebases, Pattern: "324"},
		},
	})

	encoded, err := json.Marshal(entity.Policy)
	assert.NoError(t, err)
	entity.RawPolicy = string(encoded)

	assert.NoError(t, repo.Update(ctx, entity))

	fromDB, err := repo.GetByCodebaseID(ctx, entity.CodebaseID)
	assert.NoError(t, err)

	assert.Equal(t, entity.ID, fromDB.ID)
	assert.Equal(t, entity.CodebaseID, fromDB.CodebaseID)
	assert.True(t, entity.CreatedAt.Equal(fromDB.CreatedAt))
	assert.Equal(t, entity.RawPolicy, fromDB.RawPolicy)
}

func Test_Update_not_found(t *testing.T) {
	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	assert.NoError(t, err)

	repo := NewACLRepository(d)
	ctx := context.Background()

	entity := acl.ACL{
		ID:         acl.ID(uuid.New().String()),
		CodebaseID: uuid.New().String(),
		CreatedAt:  time.Now(),
		Policy:     acl.Policy{},
	}

	assert.Error(t, repo.Update(ctx, entity))
}

package db_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"getsturdy.com/api/pkg/internal/dbtest"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/organization/db"
	"getsturdy.com/api/pkg/users"
)

func TestMain(m *testing.M) {
	if os.Getenv("E2E_TEST") != "" {
		os.Exit(m.Run())
	}
	fmt.Println("Skipping, set E2E_TEST to run.")
}

func TestAddMember(t *testing.T) {
	d := dbtest.DB(t)
	repo := db.NewMember(d)
	ctx := context.Background()

	member := &organization.Member{
		ID:             uuid.NewString(),
		OrganizationID: uuid.NewString(),
		UserID:         users.ID(uuid.NewString()),
		CreatedAt:      time.Now(),
		CreatedBy:      users.ID(uuid.NewString()),
	}

	assert.NoError(t, repo.Create(ctx, member))
	memberToRemove, err := repo.GetByUserIDAndOrganizationID(ctx, member.UserID, member.OrganizationID)
	assert.NoError(t, err)

	time := time.Now()
	memberToRemove.DeletedAt = &time
	memberToRemove.DeletedBy = &member.CreatedBy

	assert.NoError(t, repo.Update(ctx, memberToRemove))
	assert.NoError(t, repo.Create(ctx, member))

	memberToTest, err := repo.GetByUserIDAndOrganizationID(ctx, member.UserID, member.OrganizationID)
	assert.NoError(t, err)
	assert.Equal(t, memberToRemove.ID, memberToTest.ID)
	assert.Nil(t, memberToTest.DeletedAt)
	assert.Nil(t, memberToTest.DeletedBy)
}

func TestAddTwoSameMembers_returns_same_id(t *testing.T) {
	d := dbtest.DB(t)
	repo := db.NewMember(d)
	ctx := context.Background()

	member := &organization.Member{
		ID:             uuid.NewString(),
		OrganizationID: uuid.NewString(),
		UserID:         users.ID(uuid.NewString()),
		CreatedAt:      time.Now(),
		CreatedBy:      users.ID(uuid.NewString()),
	}

	assert.NoError(t, repo.Create(ctx, member))

	time := time.Now()
	member.DeletedAt = &time
	member.DeletedBy = &member.CreatedBy

	assert.NoError(t, repo.Update(ctx, member))

	oldID := member.ID
	newID := uuid.NewString()
	member.ID = newID
	assert.NoError(t, repo.Create(ctx, member))

	assert.Equal(t, oldID, member.ID)
}

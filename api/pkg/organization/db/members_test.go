package db

import (
	"context"
	"fmt"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/internal/sturdytest"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if os.Getenv("E2E_TEST") != "" {
		os.Exit(m.Run())
	}
	fmt.Println("Skipping, set E2E_TEST to run.")
}

func TestAddMember(t *testing.T) {
	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)

	assert.NoError(t, err)
	repo := NewMember(d)
	ctx := context.Background()

	member := organization.Member{
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

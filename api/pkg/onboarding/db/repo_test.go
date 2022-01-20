package db

import (
	"context"
	"os"
	"testing"
	"time"

	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/internal/sturdytest"
	"getsturdy.com/api/pkg/onboarding"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func getDB(t *testing.T) *sqlx.DB {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	assert.NoError(t, err)
	return d
}

func Test_InsertCompletedStep__double_insert_doesnt_fail(t *testing.T) {
	repo := New(getDB(t))
	ctx := context.Background()

	step := &onboarding.Step{
		UserID:    "user-id",
		ID:        "step-id",
		CreatedAt: time.Now(),
	}

	assert.NoError(t, repo.InsertCompletedStep(ctx, step))
	assert.NoError(t, repo.InsertCompletedStep(ctx, step))
}

func Test_GetCompletedSteps(t *testing.T) {
	repo := New(getDB(t))
	ctx := context.Background()

	step := &onboarding.Step{
		UserID:    "user-id",
		ID:        "step-id",
		CreatedAt: time.Now(),
	}
	assert.NoError(t, repo.InsertCompletedStep(ctx, step))

	steps, err := repo.GetCompletedSteps(ctx, step.UserID)
	assert.NoError(t, err)
	if assert.Len(t, steps, 1) {
		assert.Equal(t, step.UserID, steps[0].UserID)
		assert.Equal(t, step.ID, steps[0].ID)
	}
}

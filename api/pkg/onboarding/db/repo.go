package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/onboarding"
	"getsturdy.com/api/pkg/users"

	"github.com/jmoiron/sqlx"
)

type CompletedOnboardingStepsRepository interface {
	GetCompletedSteps(context.Context, users.ID) ([]*onboarding.Step, error)
	InsertCompletedStep(context.Context, *onboarding.Step) error
}

func New(db *sqlx.DB) CompletedOnboardingStepsRepository {
	return &completedOnboardingStepsRepository{db}
}

type completedOnboardingStepsRepository struct {
	db *sqlx.DB
}

func (c *completedOnboardingStepsRepository) GetCompletedSteps(ctx context.Context, userID users.ID) ([]*onboarding.Step, error) {
	var steps []*onboarding.Step
	if err := c.db.SelectContext(ctx, &steps, `
		SELECT step_id, user_id, created_at FROM completed_onboarding_steps
			WHERE user_id = $1
	`, userID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return steps, nil
}

func (c *completedOnboardingStepsRepository) InsertCompletedStep(ctx context.Context, step *onboarding.Step) error {
	if _, err := c.db.NamedExecContext(ctx, `
		INSERT INTO completed_onboarding_steps(user_id, step_id, created_at)
			VALUES (:user_id, :step_id, :created_at)
			ON CONFLICT (user_id, step_id) DO NOTHING
	`, step); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
)

type OnboardingRootResolver interface {
	// Queries
	CompletedOnboardingSteps(ctx context.Context) ([]OnboardingStepResolver, error)

	// Mutations
	CompleteOnboardingStep(ctx context.Context, args CompleteOnboardingStepArgs) (OnboardingStepResolver, error)

	// Subscriptions
	CompletedOnboardingStep(ctx context.Context) (chan OnboardingStepResolver, error)
}

type OnboardingStepResolver interface {
	ID() graphql.ID
}

type CompleteOnboardingStepArgs struct {
	StepID graphql.ID
}

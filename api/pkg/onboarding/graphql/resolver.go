package graphql

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/auth"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/onboarding"
	"getsturdy.com/api/pkg/onboarding/db"
	"go.uber.org/zap"

	"github.com/graph-gophers/graphql-go"
)

type onboardingRootResolver struct {
	logger *zap.Logger
	repo   db.CompletedOnboardingStepsRepository

	eventsPublisher  *eventsv2.Publisher
	eventsSubscriber *eventsv2.Subscriber
}

func NewRootResolver(
	repo db.CompletedOnboardingStepsRepository,
	logger *zap.Logger,
	eventsPublisher *eventsv2.Publisher,
	eventsSubscriber *eventsv2.Subscriber,
) resolvers.OnboardingRootResolver {
	return &onboardingRootResolver{
		logger:           logger,
		repo:             repo,
		eventsPublisher:  eventsPublisher,
		eventsSubscriber: eventsSubscriber,
	}
}

func (r *onboardingRootResolver) CompleteOnboardingStep(ctx context.Context, args resolvers.CompleteOnboardingStepArgs) (resolvers.OnboardingStepResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	step := &onboarding.Step{
		ID:        string(args.StepID),
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := r.repo.InsertCompletedStep(ctx, step); err != nil {
		return nil, gqlerrors.Error(err)
	}
	if err := r.eventsPublisher.CompletedOnboardingStep(ctx, eventsv2.User(userID), step); err != nil {
		r.logger.Error("failed to publish completed onboarding step event", zap.Error(err))
		// do not fail
	}
	return &onboardingStepResolver{step: step}, nil
}

func (r *onboardingRootResolver) CompletedOnboardingSteps(ctx context.Context) ([]resolvers.OnboardingStepResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	steps, err := r.repo.GetCompletedSteps(ctx, userID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var stepResolvers []resolvers.OnboardingStepResolver
	for _, step := range steps {
		stepResolvers = append(stepResolvers, &onboardingStepResolver{step: step})
	}
	return stepResolvers, nil
}

func (r *onboardingRootResolver) CompletedOnboardingStep(ctx context.Context) (chan resolvers.OnboardingStepResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make(chan resolvers.OnboardingStepResolver, 100)

	r.eventsSubscriber.OnCompletedOnboardingStep(ctx, eventsv2.SubscribeUser(userID), func(_ context.Context, step *onboarding.Step) error {
		select {
		case res <- &onboardingStepResolver{step: step}:
			return nil
		default:
			r.logger.Error("dropped subscription event", zap.String("step_id", step.ID), zap.Int("count", len(res)))
			return nil
		}
	})

	return res, nil
}

type onboardingStepResolver struct {
	step *onboarding.Step
}

func (o *onboardingStepResolver) ID() graphql.ID {
	return graphql.ID(o.step.ID)
}

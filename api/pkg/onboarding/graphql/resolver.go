package graphql

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/onboarding"
	"getsturdy.com/api/pkg/onboarding/db"

	"github.com/graph-gophers/graphql-go"
)

type onboardingRootResolver struct {
	repo        db.CompletedOnboardingStepsRepository
	eventSender events.EventSender
	eventReader events.EventReader
}

func NewRootResolver(repo db.CompletedOnboardingStepsRepository, eventSender events.EventSender, eventReader events.EventReader) resolvers.OnboardingRootResolver {
	return &onboardingRootResolver{repo, eventSender, eventReader}
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
	r.eventSender.User(userID, events.CompletedOnboardingStep, step.ID)
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

	cancelFunc := r.eventReader.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		if eventType == events.CompletedOnboardingStep {
			select {
			case <-ctx.Done():
				return events.ErrClientDisconnected
			case res <- &onboardingStepResolver{
				step: &onboarding.Step{
					ID:        reference,
					UserID:    userID,
					CreatedAt: time.Now(),
				},
			}:
				return nil
			default:
				return nil
			}
		}
		return nil
	})

	go func() {
		<-ctx.Done()
		cancelFunc()
		close(res)
	}()

	return res, nil
}

type onboardingStepResolver struct {
	step *onboarding.Step
}

func (o *onboardingStepResolver) ID() graphql.ID {
	return graphql.ID(o.step.ID)
}

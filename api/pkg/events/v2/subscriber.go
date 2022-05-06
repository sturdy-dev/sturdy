package events

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/notification"
	"getsturdy.com/api/pkg/onboarding"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/review"
	"getsturdy.com/api/pkg/statuses"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/views"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/pkg/workspaces/watchers"
)

type Subscriber struct {
	pubsub *pubSub
}

func NewSubscriber(
	pubsub *pubSub,
) *Subscriber {
	return &Subscriber{
		pubsub: pubsub,
	}
}

func SubscribeUser(id users.ID) Topic {
	return userTopic(id)
}

func (s *Subscriber) OnCodebaseEvent(ctx context.Context, topic Topic, callback func(context.Context, *codebases.Codebase) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Codebase, callback)
	}, topic, CodebaseEvent)
}

func (s *Subscriber) OnCodebaseUpdated(ctx context.Context, topic Topic, callback func(context.Context, *codebases.Codebase) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Codebase, callback)
	}, topic, CodebaseUpdated)
}

func (s *Subscriber) OnViewUpdated(ctx context.Context, topic Topic, callback func(context.Context, *views.View) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.View, callback)
	}, topic, ViewUpdated)
}

func (s *Subscriber) OnViewStatusUpdated(ctx context.Context, topic Topic, callback func(context.Context, *views.View) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.View, callback)
	}, topic, ViewStatusUpdated)
}

func (s *Subscriber) OnWorkspaceUpdated(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Workspace, callback)
	}, topic, WorkspaceUpdated)
}

func (s *Subscriber) OnWorkspaceUpdatedComments(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Workspace, callback)
	}, topic, WorkspaceUpdatedComments)
}

func (s *Subscriber) OnWorkspaceUpdatedReviews(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Workspace, callback)
	}, topic, WorkspaceUpdatedReviews)
}

func (s *Subscriber) OnWorkspaceUpdatedActivity(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Workspace, callback)
	}, topic, WorkspaceUpdatedActivity)
}

func (s *Subscriber) OnWorkspaceUpdatedSnapshot(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Workspace, callback)
	}, topic, WorkspaceUpdatedSnapshot)
}

func (s *Subscriber) OnWorkspaceUpdatedPresence(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Workspace, callback)
	}, topic, WorkspaceUpdatedPresence)
}

func (s *Subscriber) OnWorkspaceUpdatedSuggestion(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Workspace, callback)
	}, topic, WorkspaceUpdatedSuggestion)
}

func (s *Subscriber) OnWorkspaceWatchingStatusUpdated(ctx context.Context, topic Topic, callback func(context.Context, *watchers.Watcher) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.WorkspaceWatcher, callback)
	}, topic, WorkspaceWatchingStatusUpdated)
}

func (s *Subscriber) OnReviewUpdated(ctx context.Context, topic Topic, callback func(context.Context, *review.Review) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Review, callback)
	}, topic, ReviewUpdated)
}

func (s *Subscriber) OnGitHubPRUpdated(ctx context.Context, topic Topic, callback func(context.Context, *github.PullRequest) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.GitHubPullRequest, callback)
	}, topic, GitHubPRUpdated)
}

func (s *Subscriber) OnNotificationEvent(ctx context.Context, topic Topic, callback func(context.Context, *notification.Notification) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Notification, callback)
	}, topic, NotificationEvent)
}

func (s *Subscriber) OnStatusUpdated(ctx context.Context, topic Topic, callback func(context.Context, *statuses.Status) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Status, callback)
	}, topic, StatusUpdated)
}

func (s *Subscriber) OnCompletedOnboardingStep(ctx context.Context, topic Topic, callback func(context.Context, *onboarding.Step) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.OnboardingStep, callback)
	}, topic, CompletedOnboardingStep)
}

func (s *Subscriber) OnOrganizationUpdated(ctx context.Context, topic Topic, callback func(context.Context, *organization.Organization) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callbackWithError(ctx, event.Organization, callback)
	}, topic, OrganizationUpdated)
}

func callbackWithError[T any](ctx context.Context, value T, callback func(context.Context, T) error) error {
	if err := callback(ctx, value); err != nil {
		return fmt.Errorf("%s: %w", functionName(callback), err)
	} else {
		return err
	}
}

func functionName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

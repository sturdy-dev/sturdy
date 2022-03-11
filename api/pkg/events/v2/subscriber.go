package events

import (
	"context"

	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/notification"
	"getsturdy.com/api/pkg/onboarding"
	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/review"
	"getsturdy.com/api/pkg/statuses"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/workspaces"
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

func (s *Subscriber) OnCodebaseEvent(ctx context.Context, topic Topic, callback func(context.Context, *codebase.Codebase) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Codebase)
	}, topic, CodebaseEvent)
}

func (s *Subscriber) OnCodebaseUpdated(ctx context.Context, topic Topic, callback func(context.Context, *codebase.Codebase) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Codebase)
	}, topic, CodebaseUpdated)
}

func (s *Subscriber) OnViewUpdated(ctx context.Context, topic Topic, callback func(context.Context, *view.View) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.View)
	}, topic, ViewUpdated)
}

func (s *Subscriber) OnViewStatusUpdated(ctx context.Context, topic Topic, callback func(context.Context, *view.View) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.View)
	}, topic, ViewStatusUpdated)
}

func (s *Subscriber) OnWorkspaceUpdated(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceUpdated)
}

func (s *Subscriber) OnWorkspaceUpdatedComments(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceUpdatedComments)
}

func (s *Subscriber) OnWorkspaceUpdatedReviews(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceUpdatedReviews)
}

func (s *Subscriber) OnWorkspaceUpdatedActivity(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceUpdatedActivity)
}

func (s *Subscriber) OnWorkspaceUpdatedSnapshot(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceUpdatedSnapshot)
}

func (s *Subscriber) OnWorkspaceUpdatedPresence(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceUpdatedPresence)
}

func (s *Subscriber) OnWorkspaceUpdatedSuggestion(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceUpdatedSuggestion)
}

func (s *Subscriber) OnWorkspaceWatchingStatusUpdated(ctx context.Context, topic Topic, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, topic, WorkspaceWatchingStatusUpdated)
}

func (s *Subscriber) OnReviewUpdated(ctx context.Context, topic Topic, callback func(context.Context, *review.Review) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Review)
	}, topic, ReviewUpdated)
}

func (s *Subscriber) OnGitHubPRUpdated(ctx context.Context, topic Topic, callback func(context.Context, *github.GitHubPullRequest) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.GitHubPullRequest)
	}, topic, GitHubPRUpdated)
}

func (s *Subscriber) OnNotificationEvent(ctx context.Context, topic Topic, callback func(context.Context, *notification.Notification) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Notification)
	}, topic, NotificationEvent)
}

func (s *Subscriber) OnStatusUpdated(ctx context.Context, topic Topic, callback func(context.Context, *statuses.Status) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Status)
	}, topic, StatusUpdated)
}

func (s *Subscriber) OnCompletedOnboardingStep(ctx context.Context, topic Topic, callback func(context.Context, *onboarding.Step) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.OnboardingStep)
	}, topic, CompletedOnboardingStep)
}

func (s *Subscriber) OnOrganizationUpdated(ctx context.Context, topic Topic, callback func(context.Context, *organization.Organization) error) {
	s.pubsub.sub(ctx, func(ctx context.Context, event *event) error {
		return callback(ctx, event.Organization)
	}, topic, OrganizationUpdated)
}

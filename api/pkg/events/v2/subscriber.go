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
	pubsub *PubSub
}

func NewSubscriber(pubsub *PubSub) *Subscriber {
	return &Subscriber{
		pubsub: pubsub,
	}
}

func (s *Subscriber) Workspace(ctx context.Context, workspaceID string) *sub {
	return &sub{
		pubsub: s.pubsub,
		topic:  Workspace(workspaceID),
	}
}

func (s *Subscriber) User(ctx context.Context, userID users.ID) *sub {
	return &sub{
		pubsub: s.pubsub,
		topic:  User(userID),
	}
}

type sub struct {
	topic  Topic
	pubsub *PubSub
}

func (s *sub) OnCodebaseEvent(ctx context.Context, callback func(context.Context, *codebase.Codebase) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Codebase)
	}, s.topic, CodebaseUpdated)
}

func (s *sub) OnCodebaseUpdated(ctx context.Context, callback func(context.Context, *codebase.Codebase) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Codebase)
	}, s.topic, CodebaseUpdated)
}

func (s *sub) OnViewUpdated(ctx context.Context, callback func(context.Context, *view.View) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.View)
	}, s.topic, ViewUpdated)
}

func (s *sub) OnViewStatusUpdated(ctx context.Context, callback func(context.Context, *view.View) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.View)
	}, s.topic, ViewStatusUpdated)
}

func (s *sub) OnWorkspaceUpdated(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceUpdated)
}

func (s *sub) OnWorkspaceUpdatedComments(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceUpdatedComments)
}

func (s *sub) OnWorkspaceUpdatedReviews(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceUpdatedReviews)
}

func (s *sub) OnWorkspaceUpdatedActivity(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceUpdatedActivity)
}

func (s *sub) OnWorkspaceUpdatedSnapshot(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceUpdatedSnapshot)
}

func (s *sub) OnWorkspaceUpdatedPresence(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceUpdatedPresence)
}

func (s *sub) OnWorkspaceUpdatedSuggestion(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceUpdatedSuggestion)
}

func (s *sub) OnWorkspaceWatchingStatusUpdated(ctx context.Context, callback func(context.Context, *workspaces.Workspace) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Workspace)
	}, s.topic, WorkspaceWatchingStatusUpdated)
}

func (s *sub) OnReviewUpdated(ctx context.Context, callback func(context.Context, *review.Review) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Review)
	}, s.topic, ReviewUpdated)
}

func (s *sub) OnGitHubPRUpdated(ctx context.Context, callback func(context.Context, *github.GitHubPullRequest) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.GitHubPullRequest)
	}, s.topic, GitHubPRUpdated)
}

func (s *sub) OnNotificationEvent(ctx context.Context, callback func(context.Context, *notification.Notification) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Notification)
	}, s.topic, NotificationEvent)
}

func (s *sub) OnStatusUpdated(ctx context.Context, callback func(context.Context, *statuses.Status) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Status)
	}, s.topic, StatusUpdated)
}

func (s *sub) OnCompletedOnboardingStep(ctx context.Context, callback func(context.Context, *onboarding.Step) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.OnboardingStep)
	}, s.topic, CompletedOnboardingStep)
}

func (s *sub) OnOrganizationUpdated(ctx context.Context, callback func(context.Context, *organization.Organization) error) {
	s.pubsub.sub(func(ctx context.Context, event *event) error {
		return callback(ctx, event.Organization)
	}, s.topic, OrganizationUpdated)
}

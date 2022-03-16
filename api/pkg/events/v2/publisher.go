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
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/workspaces"

	db_codebase "getsturdy.com/api/pkg/codebase/db"
	db_organization "getsturdy.com/api/pkg/organization/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

type Publisher struct {
	pubSub *pubSub

	codebaseUserRepo       db_codebase.CodebaseUserRepository
	organizationMemberRepo db_organization.MemberRepository
	workspaceRepo          db_workspaces.WorkspaceReader
}

func NewPublisher(
	pubsub *pubSub,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	workspaceRepo db_workspaces.WorkspaceReader,
	organizationMemberRepo db_organization.MemberRepository,
) *Publisher {
	return &Publisher{
		pubSub:                 pubsub,
		codebaseUserRepo:       codebaseUserRepo,
		workspaceRepo:          workspaceRepo,
		organizationMemberRepo: organizationMemberRepo,
	}
}

func (p *Publisher) CodebaseEvent(ctx context.Context, receiver *receiver, codebase *codebase.Codebase) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:     CodebaseEvent,
			Codebase: codebase,
		})
	}
	return nil
}

func (p *Publisher) CodebaseUpdated(ctx context.Context, receiver *receiver, codebase *codebase.Codebase) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:     CodebaseUpdated,
			Codebase: codebase,
		})
	}
	return nil
}

func (p *Publisher) ViewUpdated(ctx context.Context, receiver *receiver, view *view.View) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type: ViewUpdated,
			View: view,
		})
	}
	return nil
}

func (p *Publisher) ViewStatusUpdated(ctx context.Context, receiver *receiver, view *view.View) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type: ViewStatusUpdated,
			View: view,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceUpdated(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdated,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceUpdatedComments(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedComments,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceUpdatedReviews(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedReviews,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceUpdatedActivity(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedActivity,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceUpdatedSnapshot(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedSnapshot,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceUpdatedPresence(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedPresence,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceUpdatedSuggestion(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedSuggestion,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) WorkspaceWatchingStatusUpdated(ctx context.Context, receiver *receiver, workspace *workspaces.Workspace) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceWatchingStatusUpdated,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *Publisher) ReviewUpdated(ctx context.Context, receiver *receiver, review *review.Review) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:   ReviewUpdated,
			Review: review,
		})
	}
	return nil
}

func (p *Publisher) GitHubPRUpdated(ctx context.Context, receiver *receiver, pr *github.PullRequest) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:              GitHubPRUpdated,
			GitHubPullRequest: pr,
		})
	}
	return nil
}

func (p *Publisher) NotificationEvent(ctx context.Context, receiver *receiver, notification *notification.Notification) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:         NotificationEvent,
			Notification: notification,
		})
	}
	return nil
}

func (p *Publisher) StatusUpdated(ctx context.Context, receiver *receiver, status *statuses.Status) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:   StatusUpdated,
			Status: status,
		})
	}
	return nil
}

func (p *Publisher) CompletedOnboardingStep(ctx context.Context, receiver *receiver, step *onboarding.Step) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:           CompletedOnboardingStep,
			OnboardingStep: step,
		})
	}
	return nil
}

func (p *Publisher) OrganizationUpdated(ctx context.Context, receiver *receiver, organization *organization.Organization) error {
	topics, err := receiver.Topics(ctx, p.codebaseUserRepo, p.workspaceRepo, p.organizationMemberRepo)
	if err != nil {
		return err
	}
	for topic := range topics {
		p.pubSub.pub(topic, &event{
			Type:         OrganizationUpdated,
			Organization: organization,
		})
	}
	return nil
}

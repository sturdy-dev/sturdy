package events

import (
	"context"
	"fmt"

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

	db_codebase "getsturdy.com/api/pkg/codebase/db"
	db_organization "getsturdy.com/api/pkg/organization/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

type Publisher struct {
	pubSub *PubSub

	codebaseUserRepo       db_codebase.CodebaseUserRepository
	organizationMemberRepo db_organization.MemberRepository
	workspaceRepo          db_workspaces.WorkspaceReader
}

func NewPublisher(
	pubsub *PubSub,
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

type receiver struct {
	UserIDs         []users.ID
	WorkspaceIDs    []string
	CodebaseIDs     []string
	OrganizationIDs []string
}

func Workspace(workspaceID string) *receiver {
	return &receiver{
		WorkspaceIDs: []string{workspaceID},
	}
}

func User(userID users.ID) *receiver {
	return &receiver{
		UserIDs: []users.ID{userID},
	}
}

func Codebase(codebaseID string) *receiver {
	return &receiver{
		CodebaseIDs: []string{codebaseID},
	}
}

func Organization(organizationID string) *receiver {
	return &receiver{
		OrganizationIDs: []string{organizationID},
	}
}

func (r *receiver) Topics(ctx context.Context, publisher *Publisher) (map[Topic]bool, error) {
	topics := map[Topic]bool{}
	for _, userID := range r.UserIDs {
		topics[user(userID)] = true
	}

	for _, workspaceID := range r.WorkspaceIDs {
		ws, err := publisher.workspaceRepo.Get(workspaceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
		}
		r.CodebaseIDs = append(r.CodebaseIDs, ws.CodebaseID)
	}

	for _, codebaseID := range r.CodebaseIDs {
		members, err := publisher.codebaseUserRepo.GetByCodebase(codebaseID)
		if err != nil {
			return nil, fmt.Errorf("failed to get codebase members: %w", err)
		}
		for _, member := range members {
			topics[user(member.UserID)] = true
		}
	}

	for _, organizationID := range r.OrganizationIDs {
		members, err := publisher.organizationMemberRepo.ListByOrganizationID(ctx, organizationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization members: %w", err)
		}
		for _, member := range members {
			topics[user(member.UserID)] = true
		}
	}

	return topics, nil
}

func (p *Publisher) CodebaseEvent(ctx context.Context, receiver *receiver, codebase *codebase.Codebase) error {
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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

func (p *Publisher) GitHubPRUpdated(ctx context.Context, receiver *receiver, pr *github.GitHubPullRequest) error {
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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
	topics, err := receiver.Topics(ctx, p)
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

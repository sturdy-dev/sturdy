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

	db_codebase "getsturdy.com/api/pkg/codebase/db"
	db_organization "getsturdy.com/api/pkg/organization/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

type Publisher struct {
	pubsub *PubSub

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
		pubsub:                 pubsub,
		codebaseUserRepo:       codebaseUserRepo,
		workspaceRepo:          workspaceRepo,
		organizationMemberRepo: organizationMemberRepo,
	}
}

func (s *Publisher) User(ctx context.Context, userID ...users.ID) publisher {
	topics := []Topic{}
	for _, id := range userID {
		topics = append(topics, User(id))
	}
	return publisher{
		topics: topics,
		pubSub: s.pubsub,
	}
}

func (s *Publisher) Codebase(ctx context.Context, id string) publisher {
	members, err := s.codebaseUserRepo.GetByCodebase(id)
	if err != nil {
		return publisher{
			err: err,
		}
	}
	userIDs := []users.ID{}
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}
	return s.User(ctx, userIDs...)
}

func (s *Publisher) Workspace(ctx context.Context, id string) publisher {
	ws, err := s.workspaceRepo.Get(id)
	if err != nil {
		return publisher{
			err: err,
		}
	}
	publisher := s.Codebase(ctx, ws.CodebaseID)
	publisher.topics = append(publisher.topics, Workspace(ws.ID))
	return publisher
}

func (s *Publisher) Organization(ctx context.Context, id string) publisher {
	members, err := s.organizationMemberRepo.ListByOrganizationID(ctx, id)
	if err != nil {
		return publisher{
			err: err,
		}
	}
	userIDs := []users.ID{}
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}
	return s.User(ctx, userIDs...)
}

type publisher struct {
	err    error
	topics []Topic
	pubSub *PubSub
}

func (p *publisher) CodebaseEvent(codebase *codebase.Codebase) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:     CodebaseEvent,
			Codebase: codebase,
		})
	}
	return nil
}

func (p *publisher) CodebaseUpdated(codebase *codebase.Codebase) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:     CodebaseUpdated,
			Codebase: codebase,
		})
	}
	return nil
}

func (p *publisher) ViewUpdated(view *view.View) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type: ViewUpdated,
			View: view,
		})
	}
	return nil
}

func (p *publisher) ViewStatusUpdated(view *view.View) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type: ViewStatusUpdated,
			View: view,
		})
	}
	return nil
}

func (p *publisher) WorkspaceUpdated(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdated,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) WorkspaceUpdatedComments(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedComments,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) WorkspaceUpdatedReviews(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedReviews,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) WorkspaceUpdatedActivity(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedActivity,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) WorkspaceUpdatedSnapshot(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedSnapshot,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) WorkspaceUpdatedPresence(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedPresence,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) WorkspaceUpdatedSuggestion(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceUpdatedSuggestion,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) WorkspaceWatchingStatusUpdated(workspace *workspaces.Workspace) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:      WorkspaceWatchingStatusUpdated,
			Workspace: workspace,
		})
	}
	return nil
}

func (p *publisher) ReviewUpdated(review *review.Review) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:   ReviewUpdated,
			Review: review,
		})
	}
	return nil
}

func (p *publisher) GitHubPRUpdated(pr *github.GitHubPullRequest) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:              GitHubPRUpdated,
			GitHubPullRequest: pr,
		})
	}
	return nil
}

func (p *publisher) NotificationEvent(notification *notification.Notification) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:         NotificationEvent,
			Notification: notification,
		})
	}
	return nil
}

func (p *publisher) StatusUpdated(status *statuses.Status) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:   StatusUpdated,
			Status: status,
		})
	}
	return nil
}

func (p *publisher) CompletedOnboardingStep(step *onboarding.Step) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:           CompletedOnboardingStep,
			OnboardingStep: step,
		})
	}
	return nil
}

func (p *publisher) OrganizationUpdated(organization *organization.Organization) error {
	if p.err != nil {
		return p.err
	}
	for _, topic := range p.topics {
		p.pubSub.pub(topic, &event{
			Type:         OrganizationUpdated,
			Organization: organization,
		})
	}
	return nil
}

package events

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	db_organization "getsturdy.com/api/pkg/organization/db"
	"getsturdy.com/api/pkg/users"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

// TODO: support sending multiple events. Some users of this interface call methods in a loop.
type EventSender interface {
	// User sends this event to this user only
	User(id users.ID, eventType EventType, reference string)

	// Codebase sends this event to all members of this codebase
	Codebase(id codebases.ID, eventType EventType, reference string) error

	// Workspace sends this event to all members of the codebase of this workspace
	Workspace(id string, eventType EventType, reference string) error

	// Organization sends this event to all members of this Organization
	Organization(ctx context.Context, id string, eventType EventType, reference string) error
}

type eventsSender struct {
	codebaseUserRepo       db_codebases.CodebaseUserRepository
	workspaceRepo          db_workspaces.WorkspaceReader
	organizationMemberRepo db_organization.MemberRepository

	events eventWriter
}

func NewSender(
	codebaseUserRepo db_codebases.CodebaseUserRepository,
	workspaceRepo db_workspaces.WorkspaceReader,
	organizationMemberRepo db_organization.MemberRepository,
	events EventReadWriter,
) EventSender {
	return &eventsSender{
		codebaseUserRepo:       codebaseUserRepo,
		workspaceRepo:          workspaceRepo,
		organizationMemberRepo: organizationMemberRepo,
		events:                 events,
	}
}

func (s *eventsSender) User(id users.ID, eventType EventType, reference string) {
	s.events.UserEvent(id, eventType, reference)
}

func (s *eventsSender) Codebase(id codebases.ID, eventType EventType, reference string) error {
	members, err := s.codebaseUserRepo.GetByCodebase(id)
	if err != nil {
		return err
	}

	for _, m := range members {
		s.User(m.UserID, eventType, reference)
	}

	return nil
}

func (s *eventsSender) Workspace(id string, eventType EventType, reference string) error {
	ws, err := s.workspaceRepo.Get(id)
	if err != nil {
		return err
	}
	return s.Codebase(ws.CodebaseID, eventType, reference)
}

func (s *eventsSender) Organization(ctx context.Context, id string, eventType EventType, reference string) error {
	members, err := s.organizationMemberRepo.ListByOrganizationID(ctx, id)
	if err != nil {
		return err
	}

	for _, m := range members {
		s.User(m.UserID, eventType, reference)
	}

	return nil
}

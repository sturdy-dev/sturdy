package events

import (
	db_codebase "mash/pkg/codebase/db"
	db_workspace "mash/pkg/workspace/db"
)

// TODO: support sending multiple events. Some users of this interface call methods in a loop.
type EventSender interface {
	// User sends this event to this user only
	User(id string, eventType EventType, reference string)

	// Codebase sends this event to all members of this codebase
	Codebase(id string, eventType EventType, reference string) error

	// Workspace sends this event to all members of the codebase of this workspace
	Workspace(id string, eventType EventType, reference string) error
}

type eventsSender struct {
	codebaseUserRepo db_codebase.CodebaseUserRepository
	workspaceRepo    db_workspace.WorkspaceReader

	events eventWriter
}

func NewSender(
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	workspaceRepo db_workspace.WorkspaceReader,
	events EventReadWriter,
) EventSender {
	return &eventsSender{
		codebaseUserRepo: codebaseUserRepo,
		workspaceRepo:    workspaceRepo,
		events:           events,
	}
}

func (s *eventsSender) User(id string, eventType EventType, reference string) {
	s.events.UserEvent(id, eventType, reference)
}

func (s *eventsSender) Codebase(id string, eventType EventType, reference string) error {
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
	s.events.WorkspaceEvent(id, eventType, reference)

	ws, err := s.workspaceRepo.Get(id)
	if err != nil {
		return err
	}
	return s.Codebase(ws.CodebaseID, eventType, reference)
}

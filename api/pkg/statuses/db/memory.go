package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/statuses"
)

type memory struct {
	byID          map[string]*statuses.Status
	byWorkspaceID map[string]*statuses.Status
}

func NewMemory() Repository {
	return &memory{
		byID:          make(map[string]*statuses.Status),
		byWorkspaceID: make(map[string]*statuses.Status),
	}
}

func (m *memory) Create(_ context.Context, status *statuses.Status) error {
	m.byID[status.ID] = status
	return nil
}

func (m *memory) Get(ctx context.Context, id string) (*statuses.Status, error) {
	if status, ok := m.byID[id]; ok {
		return status, nil
	}
	return nil, sql.ErrNoRows
}

func (m *memory) ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*statuses.Status, error) {
	// todo: implement
	return nil, nil
}

func (m *memory) ListByCodebaseIDAndCommitID(ctx context.Context, codebaseID codebases.ID, commitID string) ([]*statuses.Status, error) {
	// todo: implement
	return nil, nil
}

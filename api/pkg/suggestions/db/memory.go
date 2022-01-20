package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"getsturdy.com/api/pkg/suggestions"
)

var _ Repository = &memory{}

type memory struct {
	byID             map[suggestions.ID]*suggestions.Suggestion
	byWorkspaceID    map[string]*suggestions.Suggestion
	byForWorkspaceID map[string][]*suggestions.Suggestion
	byForSnapshotID  map[string][]*suggestions.Suggestion
}

func NewMemory() *memory {
	return &memory{
		byID:             make(map[suggestions.ID]*suggestions.Suggestion),
		byWorkspaceID:    make(map[string]*suggestions.Suggestion),
		byForWorkspaceID: make(map[string][]*suggestions.Suggestion),
		byForSnapshotID:  make(map[string][]*suggestions.Suggestion),
	}
}

func (m *memory) Create(ctx context.Context, s *suggestions.Suggestion) error {
	if _, err := m.GetByWorkspaceID(ctx, s.WorkspaceID); !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("suggestion already exists")
	}
	m.byID[s.ID] = s
	m.byWorkspaceID[s.WorkspaceID] = s
	m.byForWorkspaceID[s.ForWorkspaceID] = append(m.byForWorkspaceID[s.ForWorkspaceID], s)
	m.byForSnapshotID[s.ForSnapshotID] = append(m.byForSnapshotID[s.ForSnapshotID], s)
	return nil
}

func (m *memory) Update(ctx context.Context, suggestion *suggestions.Suggestion) error {
	if _, err := m.GetByWorkspaceID(ctx, suggestion.WorkspaceID); errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("suggestion not found")
	}
	m.byID[suggestion.ID] = suggestion
	m.byWorkspaceID[suggestion.WorkspaceID] = suggestion
	replaced := false
	for _, s := range m.byForWorkspaceID[suggestion.ForWorkspaceID] {
		if s.ID == suggestion.ID {
			*s = *suggestion
			replaced = true
		}
	}
	if !replaced {
		m.byForWorkspaceID[suggestion.ForWorkspaceID] = append(m.byForWorkspaceID[suggestion.ForWorkspaceID], suggestion)
	}
	return nil
}

func (m *memory) GetByID(_ context.Context, id suggestions.ID) (*suggestions.Suggestion, error) {
	if found, ok := m.byID[id]; ok {
		return found, nil
	}
	return nil, sql.ErrNoRows
}

func (m *memory) GetByWorkspaceID(_ context.Context, workspaceID string) (*suggestions.Suggestion, error) {
	if found, ok := m.byWorkspaceID[workspaceID]; ok {
		return found, nil
	}
	return nil, sql.ErrNoRows
}

func (m *memory) ListBySnapshotID(_ context.Context, snapshotID string) ([]*suggestions.Suggestion, error) {
	list := m.byForSnapshotID[snapshotID]
	return list, nil
}

func (m *memory) ListForWorkspaceID(_ context.Context, forWorkspaceID string) ([]*suggestions.Suggestion, error) {
	list := m.byForWorkspaceID[forWorkspaceID]
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

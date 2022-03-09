package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/review"
	"getsturdy.com/api/pkg/users"
)

var _ ReviewRepository = &memory{}

type memory struct {
	byID              map[string]*review.Review
	byWorkspaceByUser map[string]map[users.ID][]*review.Review
}

func NewMemory() *memory {
	return &memory{
		byID:              map[string]*review.Review{},
		byWorkspaceByUser: map[string]map[users.ID][]*review.Review{},
	}
}

func (m *memory) store(r *review.Review) {
	m.byID[r.ID] = r

	byWorkspace := m.byWorkspaceByUser[r.WorkspaceID]
	if byWorkspace == nil {
		byWorkspace = map[users.ID][]*review.Review{}
		m.byWorkspaceByUser[r.WorkspaceID] = byWorkspace
	}
	byWorkspace[r.UserID] = append(byWorkspace[r.UserID], r)
}

func (m *memory) Create(ctx context.Context, r review.Review) error {
	m.store(&r)
	return nil
}

func (m *memory) Update(ctx context.Context, r *review.Review) error {
	m.store(r)
	return nil
}
func (m *memory) Get(ctx context.Context, id string) (*review.Review, error) {
	r, ok := m.byID[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return r, nil
}

func (m *memory) GetLatestByUserAndWorkspace(ctx context.Context, userID users.ID, workspaceID string) (*review.Review, error) {
	byWorkspace := m.byWorkspaceByUser[workspaceID]
	if byWorkspace == nil {
		return nil, sql.ErrNoRows
	}

	for _, review := range byWorkspace[userID] {
		if !review.IsReplaced {
			return review, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *memory) ListLatestByWorkspace(ctx context.Context, workspaceID string) ([]*review.Review, error) {
	byWorkspace := m.byWorkspaceByUser[workspaceID]
	if byWorkspace == nil {
		return nil, nil
	}

	rr := []*review.Review{}
	for _, reviews := range byWorkspace {
		for _, review := range reviews {
			if review.DismissedAt != nil {
				continue
			}
			if review.IsReplaced {
				continue
			}
			rr = append(rr, review)
		}
	}
	return rr, nil
}

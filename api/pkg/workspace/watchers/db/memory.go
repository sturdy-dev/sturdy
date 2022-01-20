package db

import (
	"context"
	"database/sql"
	"sort"

	"mash/pkg/workspace/watchers"
)

var _ Repository = &inMemory{}

type inMemory struct {
	byWorkspaceID map[string][]*watchers.Watcher
}

func NewInMemory() *inMemory {
	return &inMemory{
		byWorkspaceID: make(map[string][]*watchers.Watcher),
	}
}

func (i *inMemory) Create(ctx context.Context, w *watchers.Watcher) error {
	i.byWorkspaceID[w.WorkspaceID] = append(i.byWorkspaceID[w.WorkspaceID], w)
	return nil
}

func (i *inMemory) ListWatchingByWorkspaceID(ctx context.Context, workspaceID string) ([]*watchers.Watcher, error) {
	// group by userID
	latestByUserID := map[string]*watchers.Watcher{}
	for _, watcher := range i.byWorkspaceID[workspaceID] {
		latestByUserID[watcher.UserID] = watcher
	}

	// filter non ignored
	nonIgnored := make([]*watchers.Watcher, 0, len(latestByUserID))
	for _, watcher := range latestByUserID {
		if watcher.Status == watchers.StatusIgnored {
			continue
		}
		nonIgnored = append(nonIgnored, watcher)
	}

	// sort by CreatedAt
	sort.Slice(nonIgnored, func(i, j int) bool {
		return nonIgnored[i].CreatedAt.Before(nonIgnored[j].CreatedAt)
	})

	return nonIgnored, nil
}

func (i *inMemory) GetByUserIDAndWorkspaceID(ctx context.Context,userID string, workspaceID string) (*watchers.Watcher, error) {
	byWorkspaceID := i.byWorkspaceID[workspaceID]
	byUserID := make([]*watchers.Watcher, 0, len(byWorkspaceID))
	for _, watcher := range byWorkspaceID {
		if watcher.UserID != userID {
			continue
		}
		byUserID = append(byUserID, watcher)
	}
	if len(byUserID) == 0 {
		return nil, sql.ErrNoRows
	}
	// sort by CreatedAt
	sort.Slice(byUserID, func(i, j int) bool{
		return byUserID[i].CreatedAt.Before(byUserID[j].CreatedAt)
	})
	return byUserID[0], nil
}

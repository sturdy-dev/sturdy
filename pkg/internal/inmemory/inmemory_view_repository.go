package inmemory

import (
	"context"
	"database/sql"
	"fmt"
	"mash/pkg/view"
	db_view "mash/pkg/view/db"
	"sort"
)

type inMemoryViewRepo struct {
	views []view.View
}

func NewInMemoryViewRepo() db_view.Repository {
	return &inMemoryViewRepo{
		views: make([]view.View, 0),
	}
}

func (f *inMemoryViewRepo) Create(entity view.View) error {
	f.views = append(f.views, entity)
	return nil
}

func (f *inMemoryViewRepo) Get(id string) (*view.View, error) {
	for _, v := range f.views {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *inMemoryViewRepo) ListByCodebase(codebaseID string) ([]*view.View, error) {
	var res []*view.View
	for _, v := range f.views {
		if v.CodebaseID == codebaseID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) ListByUser(userID string) ([]*view.View, error) {
	var res []*view.View
	for _, v := range f.views {
		if v.UserID == userID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) LastUsedByCodebaseAndUser(ctx context.Context, codebaseID, userID string) (*view.View, error) {
	views, _ := f.ListByCodebaseAndUser(codebaseID, userID)
	if len(views) < 1 {
		return nil, sql.ErrNoRows
	}
	sort.Slice(views, func(i, j int) bool {
		a := views[i]
		b := views[j]
		return a.LastUsedAt.Before(*b.LastUsedAt)
	})
	return views[0], nil
}

func (f *inMemoryViewRepo) ListByCodebaseAndUser(codebaseID, userID string) ([]*view.View, error) {
	var res []*view.View
	for _, v := range f.views {
		if v.CodebaseID == codebaseID && v.UserID == userID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) ListByCodebaseAndWorkspace(codebaseID, workspaceID string) ([]*view.View, error) {
	var res []*view.View
	for _, v := range f.views {
		if v.CodebaseID == codebaseID && v.WorkspaceID == workspaceID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) Update(e *view.View) error {
	for i, v := range f.views {
		if v.ID == e.ID {
			f.views[i] = *e
			return nil
		}
	}
	return fmt.Errorf("view not found in update")
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/views"
)

type inMemoryViewRepo struct {
	views []views.View
}

func NewInMemoryViewRepo() Repository {
	return &inMemoryViewRepo{
		views: make([]views.View, 0),
	}
}

func (f *inMemoryViewRepo) Create(entity views.View) error {
	f.views = append(f.views, entity)
	return nil
}

func (f *inMemoryViewRepo) Get(id string) (*views.View, error) {
	for _, v := range f.views {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *inMemoryViewRepo) ListByCodebase(codebaseID codebases.ID) ([]*views.View, error) {
	var res []*views.View
	for _, v := range f.views {
		if v.CodebaseID == codebaseID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) ListByUser(userID users.ID) ([]*views.View, error) {
	var res []*views.View
	for _, v := range f.views {
		if v.UserID == userID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) LastUsedByCodebaseAndUser(ctx context.Context, codebaseID codebases.ID, userID users.ID) (*views.View, error) {
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

func (f *inMemoryViewRepo) ListByCodebaseAndUser(codebaseID codebases.ID, userID users.ID) ([]*views.View, error) {
	var res []*views.View
	for _, v := range f.views {
		if v.CodebaseID == codebaseID && v.UserID == userID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) ListByCodebaseAndWorkspace(codebaseID codebases.ID, workspaceID string) ([]*views.View, error) {
	var res []*views.View
	for _, v := range f.views {
		if v.CodebaseID == codebaseID && v.WorkspaceID == workspaceID {
			vv := v
			res = append(res, &vv)
		}
	}
	return res, nil
}

func (f *inMemoryViewRepo) Update(e *views.View) error {
	for i, v := range f.views {
		if v.ID == e.ID {
			f.views[i] = *e
			return nil
		}
	}
	return fmt.Errorf("view not found in update")
}

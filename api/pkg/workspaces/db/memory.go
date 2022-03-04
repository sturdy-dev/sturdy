package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/workspaces"
)

var _ Repository = &memory{}

type memory struct {
	workspaces []workspaces.Workspace
}

func NewMemory() *memory {
	return &memory{workspaces: make([]workspaces.Workspace, 0)}
}

func (f *memory) Create(entity workspaces.Workspace) error {
	f.workspaces = append(f.workspaces, entity)
	return nil
}
func (f *memory) Get(id string) (*workspaces.Workspace, error) {
	for _, ws := range f.workspaces {
		if ws.ID == id {
			return &ws, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *memory) ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspaces.Workspace, error) {
	panic("not implemented")
}

func (f *memory) ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspaces.Workspace, error) {
	panic("not implemented")
}

func (f *memory) Update(_ context.Context, entity *workspaces.Workspace) error {
	for idx, ws := range f.workspaces {
		if ws.ID == entity.ID {
			f.workspaces[idx] = *entity
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *memory) GetByViewID(viewId string, includeArchived bool) (*workspaces.Workspace, error) {
	for _, ws := range f.workspaces {
		if ws.ViewID != nil && *ws.ViewID == viewId &&
			(includeArchived || ws.ArchivedAt == nil) {
			return &ws, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *memory) GetBySnapshotID(id string) (*workspaces.Workspace, error) {
	panic("not implemented")
}

func (f *memory) UnsetUpToDateWithTrunkForAllInCodebase(codebaseID string) error {
	for idx, ws := range f.workspaces {
		if ws.CodebaseID == codebaseID {
			f.workspaces[idx].UpToDateWithTrunk = nil
		}
	}
	return nil
}

func (f *memory) SetUpToDateWithTrunk(_ context.Context, workspaceID string, upToDateWithTrunk bool) error {
	for _, ws := range f.workspaces {
		if ws.ID == workspaceID {
			ws.UpToDateWithTrunk = &upToDateWithTrunk
			return nil
		}
	}
	return nil
}

func (f *memory) SetHeadChange(_ context.Context, workspaceID string, changeID changes.ID) error {
	for _, ws := range f.workspaces {
		if ws.ID == workspaceID {
			ws.HeadChangeComputed = true
			ws.HeadChangeID = &changeID
			return nil
		}
	}
	return nil
}

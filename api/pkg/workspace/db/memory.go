package db

import (
	"database/sql"

	"getsturdy.com/api/pkg/workspace"
)

var _ Repository = &memory{}

type memory struct {
	workspaces []workspace.Workspace
}

func NewMemory() *memory {
	return &memory{workspaces: make([]workspace.Workspace, 0)}
}

func (f *memory) Create(entity workspace.Workspace) error {
	f.workspaces = append(f.workspaces, entity)
	return nil
}
func (f *memory) Get(id string) (*workspace.Workspace, error) {
	for _, ws := range f.workspaces {
		if ws.ID == id {
			return &ws, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *memory) ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspace.Workspace, error) {
	panic("not implemented")
}

func (f *memory) ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspace.Workspace, error) {
	panic("not implemented")
}

func (f *memory) Update(entity *workspace.Workspace) error {
	for idx, ws := range f.workspaces {
		if ws.ID == entity.ID {
			f.workspaces[idx] = *entity
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *memory) GetByViewID(viewId string, includeArchived bool) (*workspace.Workspace, error) {
	for _, ws := range f.workspaces {
		if ws.ViewID != nil && *ws.ViewID == viewId &&
			(includeArchived || ws.ArchivedAt == nil) {
			return &ws, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *memory) GetBySnapshotID(id string) (*workspace.Workspace, error) {
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

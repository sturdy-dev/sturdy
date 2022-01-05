package inmemory

import (
	"database/sql"
	"mash/pkg/workspace"
	db_workspace "mash/pkg/workspace/db"
)

type inMemoryWorkspaceRepo struct {
	workspaces []workspace.Workspace
}

func NewInMemoryWorkspaceRepo() db_workspace.Repository {
	return &inMemoryWorkspaceRepo{workspaces: make([]workspace.Workspace, 0)}
}

func (f *inMemoryWorkspaceRepo) Create(entity workspace.Workspace) error {
	f.workspaces = append(f.workspaces, entity)
	return nil
}
func (f *inMemoryWorkspaceRepo) Get(id string) (*workspace.Workspace, error) {
	for _, ws := range f.workspaces {
		if ws.ID == id {
			return &ws, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *inMemoryWorkspaceRepo) ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspace.Workspace, error) {
	panic("not implemented")
}

func (f *inMemoryWorkspaceRepo) ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspace.Workspace, error) {
	panic("not implemented")
}

func (f *inMemoryWorkspaceRepo) Update(entity *workspace.Workspace) error {
	for idx, ws := range f.workspaces {
		if ws.ID == entity.ID {
			f.workspaces[idx] = *entity
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *inMemoryWorkspaceRepo) GetByViewID(viewId string, includeArchived bool) (*workspace.Workspace, error) {
	for _, ws := range f.workspaces {
		if ws.ViewID != nil && *ws.ViewID == viewId &&
			(includeArchived || ws.ArchivedAt == nil) {
			return &ws, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *inMemoryWorkspaceRepo) GetBySnapshotID(id string) (*workspace.Workspace, error) {
	panic("not implemented")
}

func (f *inMemoryWorkspaceRepo) UnsetUpToDateWithTrunkForAllInCodebase(codebaseID string) error {
	for idx, ws := range f.workspaces {
		if ws.CodebaseID == codebaseID {
			f.workspaces[idx].UpToDateWithTrunk = nil
		}
	}
	return nil
}

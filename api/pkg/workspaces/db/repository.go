package db

import (
	"context"

	"getsturdy.com/api/pkg/workspaces"
)

type Repository interface {
	WorkspaceWriter
	WorkspaceReader
}

type WorkspaceWriter interface {
	Create(workspace workspaces.Workspace) error
	Update(context.Context, *workspaces.Workspace) error
	UnsetUpToDateWithTrunkForAllInCodebase(codebaseID string) error
}

type WorkspaceReader interface {
	Get(id string) (*workspaces.Workspace, error)
	ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspaces.Workspace, error)
	ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspaces.Workspace, error)
	GetByViewID(viewID string, includeArchived bool) (*workspaces.Workspace, error)
	GetBySnapshotID(snapshotID string) (*workspaces.Workspace, error)
}

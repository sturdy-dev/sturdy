package db

import (
	"context"

	"getsturdy.com/api/pkg/workspace"
)

type Repository interface {
	WorkspaceWriter
	WorkspaceReader
}

type WorkspaceWriter interface {
	Create(workspace workspace.Workspace) error
	Update(context.Context, *workspace.Workspace) error
	UnsetUpToDateWithTrunkForAllInCodebase(codebaseID string) error
}

type WorkspaceReader interface {
	Get(id string) (*workspace.Workspace, error)
	ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspace.Workspace, error)
	ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspace.Workspace, error)
	GetByViewID(viewID string, includeArchived bool) (*workspace.Workspace, error)
	GetBySnapshotID(snapshotID string) (*workspace.Workspace, error)
}

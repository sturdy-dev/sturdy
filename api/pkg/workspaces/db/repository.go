package db

import (
	"context"

	"getsturdy.com/api/pkg/changes"
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
	SetUpToDateWithTrunk(context.Context, string, bool) error
	SetHeadChange(context.Context, string, *changes.ID) error
}

type WorkspaceReader interface {
	Get(id string) (*workspaces.Workspace, error)
	ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspaces.Workspace, error)
	ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspaces.Workspace, error)
	GetByViewID(viewID string, includeArchived bool) (*workspaces.Workspace, error)
	GetBySnapshotID(snapshotID string) (*workspaces.Workspace, error)
}

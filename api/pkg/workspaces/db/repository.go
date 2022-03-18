package db

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/workspaces"
)

type Repository interface {
	WorkspaceWriter
	WorkspaceReader
}

type WorkspaceWriter interface {
	Create(workspace workspaces.Workspace) error
	UnsetUpToDateWithTrunkForAllInCodebase(codebases.ID) error

	UpdateFields(ctx context.Context, workspaceID string, fields ...UpdateOption) error
}

type WorkspaceReader interface {
	Get(id string) (*workspaces.Workspace, error)
	ListByCodebaseIDs(codebaseIDs []codebases.ID, includeArchived bool) ([]*workspaces.Workspace, error)
	ListByCodebaseIDsAndUserID(codebaseIDs []codebases.ID, userID string) ([]*workspaces.Workspace, error)
	GetByViewID(viewID string, includeArchived bool) (*workspaces.Workspace, error)
	GetBySnapshotID(snapshotID string) (*workspaces.Workspace, error)
}

type UpdateOptions struct {
	updatedAt    *time.Time
	updatedAtSet bool

	upToDateWithTrunk    *bool
	upToDateWithTrunkSet bool

	headChangeID    *changes.ID
	headChangeIDSet bool

	headChangeComputed    bool
	headChangeComputedSet bool

	latestSnapshotID    *string
	latestSnapshotIDSet bool

	diffsCount    *int32
	diffsCountSet bool

	viewID    *string
	viewIDSet bool

	lastLandedAt    *time.Time
	lastLandedAtSet bool

	changeID    *changes.ID
	changeIDSet bool

	draftDescription    string
	draftDescriptionSet bool

	archivedAt    *time.Time
	archivedAtSet bool

	unarchivedAt    *time.Time
	unarchivedAtSet bool

	name    *string
	nameSet bool
}

type UpdateOption func(*UpdateOptions)

type Options []UpdateOption

func (o Options) Parse() *UpdateOptions {
	opts := &UpdateOptions{}
	for _, opt := range o {
		opt(opts)
	}
	return opts
}

func SetUpdatedAt(updatedAt *time.Time) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.updatedAt = updatedAt
		opts.updatedAtSet = true
	}
}

func SetUpToDateWithTrunk(upToDateWithTrunk *bool) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.upToDateWithTrunk = upToDateWithTrunk
		opts.upToDateWithTrunkSet = true
	}
}

func SetHeadChangeID(headChangeID *changes.ID) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.headChangeID = headChangeID
		opts.headChangeIDSet = true
	}
}

func SetHeadChangeComputed(headChangeComputed bool) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.headChangeComputed = headChangeComputed
		opts.headChangeComputedSet = true
	}
}

func SetLatestSnapshotID(latestSnapshotID *string) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.latestSnapshotID = latestSnapshotID
		opts.latestSnapshotIDSet = true
	}
}

func SetDiffsCount(diffsCount *int32) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.diffsCount = diffsCount
		opts.diffsCountSet = true
	}
}

func SetViewID(viewID *string) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.viewID = viewID
		opts.viewIDSet = true
	}
}

func SetLastLandedAt(lastLandedAt *time.Time) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.lastLandedAt = lastLandedAt
		opts.lastLandedAtSet = true
	}
}

func SetChangeID(changeID *changes.ID) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.changeID = changeID
		opts.changeIDSet = true
	}
}

func SetDraftDescription(draftDescription string) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.draftDescription = draftDescription
		opts.draftDescriptionSet = true
	}
}

func SetArchivedAt(archivedAt *time.Time) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.archivedAt = archivedAt
		opts.archivedAtSet = true
	}
}

func SetUnarchivedAt(unarchivedAt *time.Time) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.unarchivedAt = unarchivedAt
		opts.unarchivedAtSet = true
	}
}

func SetName(name *string) UpdateOption {
	return func(opts *UpdateOptions) {
		opts.name = name
		opts.nameSet = true
	}
}

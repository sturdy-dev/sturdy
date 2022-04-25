package db

import (
	"context"
	"database/sql"
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
)

// snapshotRepo implements snapshot.Repository
type snapshotRepo struct {
	byID                     map[string]*snapshots.Snapshot
	latestInView             map[string]*snapshots.Snapshot
	latestInViewAndWorkspace map[string]*snapshots.Snapshot
}

func NewInMemorySnapshotRepo() Repository {
	return &snapshotRepo{
		byID:                     make(map[string]*snapshots.Snapshot),
		latestInView:             make(map[string]*snapshots.Snapshot),
		latestInViewAndWorkspace: make(map[string]*snapshots.Snapshot),
	}
}

func (f *snapshotRepo) Create(snapshot *snapshots.Snapshot) error {
	f.byID[snapshot.ID] = snapshot
	f.latestInView[snapshot.ViewID] = snapshot
	f.latestInViewAndWorkspace[snapshot.ViewID+snapshot.WorkspaceID] = snapshot
	return nil
}
func (f *snapshotRepo) ListByView(viewID string) ([]*snapshots.Snapshot, error) {
	panic("not implemented")
}
func (f *snapshotRepo) LatestInView(viewID string) (*snapshots.Snapshot, error) {
	if snap, ok := f.latestInView[viewID]; ok {
		return snap, nil
	}
	return nil, sql.ErrNoRows
}
func (f *snapshotRepo) LatestInViewAndWorkspace(viewID, workspaceID string) (*snapshots.Snapshot, error) {
	if snap, ok := f.latestInViewAndWorkspace[viewID+workspaceID]; ok {
		return snap, nil
	}
	return nil, sql.ErrNoRows
}
func (f *snapshotRepo) Get(ID string) (*snapshots.Snapshot, error) {
	if snap, ok := f.byID[ID]; ok {
		return snap, nil
	}
	return nil, sql.ErrNoRows
}

func (f *snapshotRepo) ListUndeletedInCodebase(_ codebases.ID, _ time.Time) ([]*snapshots.Snapshot, error) {
	panic("not implemented")
}

func (f *snapshotRepo) Update(*snapshots.Snapshot) error {
	panic("not implemented")
}

func (f *snapshotRepo) ListByViewCopiedFromBranchName(copiedFromBranchName string) ([]*snapshots.Snapshot, error) {
	panic("not implemented")
}

func (f *snapshotRepo) GetByCommitSHA(_ context.Context, sha string) (*snapshots.Snapshot, error) {
	panic("not implemented")
}

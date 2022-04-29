package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/snapshots"
)

// snapshotRepo implements snapshot.Repository
type snapshotRepo struct {
	byID              map[snapshots.ID]*snapshots.Snapshot
	latestInWorkspace map[string]*snapshots.Snapshot
}

func NewInMemorySnapshotRepo() Repository {
	return &snapshotRepo{
		byID:              make(map[snapshots.ID]*snapshots.Snapshot),
		latestInWorkspace: make(map[string]*snapshots.Snapshot),
	}
}

func (f *snapshotRepo) Create(snapshot *snapshots.Snapshot) error {
	f.byID[snapshot.ID] = snapshot
	f.latestInWorkspace[snapshot.WorkspaceID] = snapshot
	return nil
}

func (f *snapshotRepo) LatestInWorkspace(_ context.Context, workspace_id string) (*snapshots.Snapshot, error) {
	if snap, ok := f.latestInWorkspace[workspace_id]; ok {
		return snap, nil
	}
	return nil, sql.ErrNoRows
}

func (f *snapshotRepo) Get(id snapshots.ID) (*snapshots.Snapshot, error) {
	if snap, ok := f.byID[id]; ok {
		return snap, nil
	}
	return nil, sql.ErrNoRows
}

func (f *snapshotRepo) Update(s *snapshots.Snapshot) error {
	f.byID[s.ID] = s
	return nil
}

func (f *snapshotRepo) GetByCommitSHA(_ context.Context, sha string) (*snapshots.Snapshot, error) {
	panic("not implemented")
}

func (f *snapshotRepo) ListByIDs(ctx context.Context, ids []snapshots.ID) ([]*snapshots.Snapshot, error) {
	res := []*snapshots.Snapshot{}
	for _, id := range ids {
		if snap, ok := f.byID[id]; ok {
			res = append(res, snap)
		}
	}
	return res, nil
}

func (r *snapshotRepo) GetByPreviousSnapshotID(ctx context.Context, previousSnapshotID snapshots.ID) (*snapshots.Snapshot, error) {
	for _, snap := range r.byID {
		if snap.PreviousSnapshotID != nil && *snap.PreviousSnapshotID == previousSnapshotID {
			return snap, nil
		}
	}
	return nil, sql.ErrNoRows
}

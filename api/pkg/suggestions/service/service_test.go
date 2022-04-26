package service_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"testing"

	"getsturdy.com/api/pkg/analytics/disabled"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	db_changes "getsturdy.com/api/pkg/changes/db"
	service_change "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	vcs_codebase "getsturdy.com/api/pkg/codebases/vcs"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/notification/sender"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/suggestions"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
	"getsturdy.com/api/vcs/testutil"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// todo:
// * test new file on both sides
// * construct diffs in the code maybe? there are too many testdata files already

var (
	//go:embed testdata/move_and_add_chunk_at_the_beginning.diff
	moveAndAddChunkAtTheBeginning []byte

	//go:embed testdata/original.txt
	original []byte
	//go:embed testdata/minus_original.diff
	minusOriginalDiff []byte
	//go:embed testdata/plus_original.diff
	plusOriginalDiff []byte

	//go:embed testdata/plus_start_chunk.txt
	plusStartChunk []byte
	//go:embed testdata/plus_start_chunk.diff
	plusStartChunkDiff []byte

	//go:embed testdata/plus_middle_chunk.txt
	plusMiddleChunk []byte
	//go:embed testdata/plus_middle_chunk.diff
	plusMiddleChunkDiff []byte

	//go:embed testdata/plus_end_chunk.txt
	plusEndChunk []byte
	//go:embed testdata/plus_end_chunk.diff
	plusEndChunkDiff []byte

	//go:embed testdata/plus_two_chunks.txt
	plusTwoChunks []byte
	//go:embed testdata/plus_two_chunks.1.diff
	plusTwoChunksHunk1 []byte
	//go:embed testdata/plus_two_chunks.2.diff
	plusTwoChunksHunk2 []byte
)

type diffs struct {
	write  map[string][]byte
	delete []string
}

type operation struct {
	applyHunks   []string
	dismissHunks []string
	dismiss      bool

	writeOriginal map[string][]byte
	suggesting    *diffs

	result []unidiff.FileDiff
}

func (o *operation) openSuggestion(t *testing.T, test *test) {
	if test.suggestion != nil {
		return
	}

	suggestion, err := test.suggestionService.Create(context.Background(), test.suggestingUserID, test.originalWorkspace)
	assert.NoError(t, err)
	test.suggestion = suggestion

	suggestingWorkspace, err := test.workspaceDB.Get(suggestion.WorkspaceID)
	assert.NoError(t, err)
	test.suggestingWorkspace = suggestingWorkspace

	// connect a workspace to suggestingView
	test.suggestingViewID = fmt.Sprintf("%s-view", test.suggestingUserID)
	assert.NoError(t, test.viewDB.Create(view.View{
		ID:         test.suggestingViewID,
		UserID:     test.suggestingUserID,
		CodebaseID: test.codebaseID,
	}))
	assert.NoError(t, test.workspaceDB.UpdateFields(context.TODO(), suggestingWorkspace.ID, db_workspaces.SetViewID(&test.suggestingViewID)))
	assert.NoError(t, vcs_view.Create(test.codebaseID, suggestingWorkspace.ID, test.suggestingViewID)(test.repoProvider))
}

func (o *operation) validate(t *testing.T, test *test) {
	if test.suggestion == nil {
		return
	}

	result, err := test.suggestionService.Diffs(context.Background(), test.suggestion)
	if assert.NoError(t, err) {
		assert.Equal(t, o.result, result)
	}
}

func (o *operation) snapshotSuggesting(t *testing.T, test *test) {
	suggestingSnapshot, err := test.gitSnapshotter.Snapshot(test.codebaseID, test.suggestingWorkspace.ID, snapshots.ActionViewSync, service_snapshots.WithOnView(test.suggestingViewID))
	assert.NoError(t, err)
	test.suggestingWorkspace.LatestSnapshotID = &suggestingSnapshot.ID
	assert.NoError(t, test.workspaceDB.UpdateFields(context.TODO(), test.suggestingWorkspace.ID, db_workspaces.SetLatestSnapshotID(&suggestingSnapshot.ID)))
}

func (o *operation) snapshotOriginal(t *testing.T, test *test) {
	snapshot, err := test.gitSnapshotter.Snapshot(test.codebaseID, test.originalWorkspace.ID, snapshots.ActionSuggestionApply, service_snapshots.WithOnView(test.originalViewID))
	assert.NoError(t, err)
	test.originalWorkspace.LatestSnapshotID = &snapshot.ID
	assert.NoError(t, test.workspaceDB.UpdateFields(context.TODO(), test.originalWorkspace.ID, db_workspaces.SetLatestSnapshotID(&snapshot.ID)))
}

func (o *operation) setupOriginal(t *testing.T, test *test) {
	if test.originalWorkspace != nil {
		return
	}
	o.setupOriginalWorkspace(t, test)
	o.setupOriginalView(t, test)
}

func (o *operation) setupOriginalWorkspace(t *testing.T, test *test) {
	if test.originalWorkspace == nil {
		originalWorkspace, err := test.workspaceService.Create(context.TODO(), service_workspace.CreateWorkspaceRequest{
			UserID:     test.originalUserID,
			CodebaseID: test.codebaseID,
			Name:       fmt.Sprintf("%s's workspace", test.originalUserID),
		})
		assert.NoError(t, err)
		test.originalWorkspace = originalWorkspace
	}
}

func (o *operation) setupOriginalView(t *testing.T, test *test) {
	test.originalViewID = fmt.Sprintf("%s-view", test.originalUserID)
	assert.NoError(t, test.viewDB.Create(view.View{
		ID:         test.originalViewID,
		UserID:     test.originalUserID,
		CodebaseID: test.codebaseID,
	}))
	test.originalWorkspace.ViewID = &test.originalViewID
	assert.NoError(t, test.workspaceDB.UpdateFields(context.TODO(), test.originalWorkspace.ID, db_workspaces.SetViewID(&test.originalViewID)))
	vcs_view.Create(test.codebaseID, test.originalWorkspace.ID, test.originalViewID)(test.repoProvider)
}

func (o *operation) run(t *testing.T, test *test) {
	o.setupOriginal(t, test)

	switch {
	case o.writeOriginal != nil:
		// make some changes
		viewPath := test.repoProvider.ViewPath(test.codebaseID, test.originalViewID)
		for filepath, content := range o.writeOriginal {
			t.Logf("original: writing %s", filepath)
			assert.NoError(t, os.WriteFile(path.Join(viewPath, filepath), content, 0777))
		}

		o.snapshotOriginal(t, test)
		o.validate(t, test)
	case o.suggesting != nil:
		o.snapshotOriginal(t, test)
		o.openSuggestion(t, test)

		suggestingViewPath := test.repoProvider.ViewPath(test.codebaseID, test.suggestingViewID)
		// delete
		for _, filepath := range o.suggesting.delete {
			t.Logf("suggesting: deleting %s", filepath)
			fp := path.Join(suggestingViewPath, filepath)
			assert.NoError(t, os.RemoveAll(fp), fp)
		}
		// write
		for filepath, content := range o.suggesting.write {
			t.Logf("suggesting: writing %s", filepath)
			fp := path.Join(suggestingViewPath, filepath)
			d := path.Dir(fp)
			assert.NoError(t, os.MkdirAll(d, 0777))
			assert.NoError(t, os.WriteFile(fp, content, 0777))
		}
		o.snapshotSuggesting(t, test)
		o.validate(t, test)
	case o.applyHunks != nil:
		t.Logf("applying hunks")
		if assert.NoError(t, test.suggestionService.ApplyHunks(context.Background(), test.suggestion, o.applyHunks...)) {
			o.validate(t, test)
		}
	case o.dismissHunks != nil:
		t.Logf("dismissing hunks")
		if assert.NoError(t, test.suggestionService.DismissHunks(context.Background(), test.suggestion, o.dismissHunks...)) {
			o.validate(t, test)
		}
	case o.dismiss:
		t.Logf("dismissing sugestion")
		if assert.NoError(t, test.suggestionService.Dismiss(context.Background(), test.suggestion)) {
			o.validate(t, test)
		}
	}
}

type test struct {
	repoProvider      provider.RepoProvider
	executorProvider  executor.Provider
	suggestionRepo    db_suggestions.Repository
	viewDB            db_view.Repository
	workspaceDB       db_workspaces.Repository
	snapshotsDB       db_snapshots.Repository
	codebaseUserRepo  db_codebases.CodebaseUserRepository
	gitSnapshotter    *service_snapshots.Service
	workspaceService  *service_workspace.Service
	suggestionService *service_suggestions.Service

	codebaseID codebases.ID

	originalUserID    users.ID
	originalViewID    string
	originalWorkspace *workspaces.Workspace

	suggestingUserID    users.ID
	suggestingViewID    string
	suggestingWorkspace *workspaces.Workspace
	suggestion          *suggestions.Suggestion

	operations []*operation
}

func newTest(t *testing.T, operations []*operation) *test {
	repoProvider := testutil.TestingRepoProvider(t)
	executorProvider := executor.NewProvider(zap.NewNop(), repoProvider)
	suggestionRepo := db_suggestions.NewMemory()

	viewDB := db_view.NewInMemoryViewRepo()
	workspaceDB := db_workspaces.NewMemory()
	snapshotsDB := db_snapshots.NewInMemorySnapshotRepo()
	codebaseUserRepo := db_codebases.NewInMemoryCodebaseUserRepo()
	logger := zap.NewNop()
	eventsSender := events.NewSender(codebaseUserRepo, workspaceDB, nil, events.NewInMemory(logger))
	changeRepo := db_changes.NewInMemoryRepo()

	analyticsService := service_analytics.New(zap.NewNop(), disabled.NewClient(zap.NewNop()))
	gitSnapshotter := service_snapshots.New(snapshotsDB, workspaceDB, workspaceDB, viewDB, suggestionRepo, eventsSender, nil, executorProvider, zap.NewNop(), analyticsService)
	changeService := service_change.New(changeRepo, nil, zap.NewNop(), executorProvider, gitSnapshotter)
	workspaceService := service_workspace.New(zap.NewNop(), analyticsService, workspaceDB, workspaceDB, changeService, nil /*viewService*/, nil /*usersService*/, executorProvider, nil /*eventsSender*/, nil /*eventsSernderv2*/, gitSnapshotter)
	suggestionService := service_suggestions.New(zap.NewNop(), suggestionRepo, workspaceService, executorProvider, gitSnapshotter, analyticsService, sender.NewNoopNotificationSender(), eventsSender)
	return &test{
		repoProvider:      repoProvider,
		executorProvider:  executorProvider,
		suggestionRepo:    suggestionRepo,
		viewDB:            viewDB,
		workspaceDB:       workspaceDB,
		snapshotsDB:       snapshotsDB,
		codebaseUserRepo:  codebaseUserRepo,
		gitSnapshotter:    gitSnapshotter,
		workspaceService:  workspaceService,
		suggestionService: suggestionService,

		originalUserID:   "user1",
		codebaseID:       "codebaseID",
		suggestingUserID: "user2",

		operations: operations,
	}
}

func (test *test) run(t *testing.T) {
	// create codebase trunk repo
	assert.NoError(t, vcs_codebase.Create(test.codebaseID)(test.repoProvider))

	for _, operation := range test.operations {
		operation.run(t, test)
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name       string
		operations []*operation
	}{
		{
			name: "rename and add chunk at the beginning",
			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						delete: []string{"file"},
						write:  map[string][]byte{"file.new": plusStartChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file.new",
							PreferredName: "file.new",
							IsMoved:       true,
							Hunks: []unidiff.Hunk{
								{
									ID:    "f949828cc55fa98a1c4376bfe30272346088abbedd56bbda0f5eb219cdb6a0d9",
									Patch: string(moveAndAddChunkAtTheBeginning),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "new file",
			operations: []*operation{
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": original},
					},

					result: []unidiff.FileDiff{
						{
							OrigName:      "/dev/null",
							NewName:       "file",
							PreferredName: "file",
							IsNew:         true,
							Hunks: []unidiff.Hunk{
								{
									ID:    "b6e93f63ca1afe4d1f5c70b187ae2197b296e895f72635f2f9f8115392d104f9",
									Patch: string(plusOriginalDiff),
								},
							},
						},
					},
				},
			},
		},

		{
			name: "move file",
			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write:  map[string][]byte{"file.new": original},
						delete: []string{"file"},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file.new",
							PreferredName: "file.new",
							IsMoved:       true,
							Hunks: []unidiff.Hunk{
								{
									ID:    "4dd77342d6b087367f5708a3039ca8a33db2cca34fae4797aad9f0f66d9c8e28",
									Patch: "diff --git \"a/file\" \"b/file.new\"\nsimilarity index 100%\nrename from \"file\"\nrename to \"file.new\"\n",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "delete file",
			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						delete: []string{"file"},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "/dev/null",
							PreferredName: "file",
							IsDeleted:     true,
							Hunks: []unidiff.Hunk{
								{
									ID:    "f4b61e064338d2edf7fb902634553e0553adff59ae88616bbf29ba180e256958",
									Patch: string(minusOriginalDiff),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "apply add chunk at the beginning",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusStartChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125",
									Patch: string(plusStartChunkDiff),
								},
							},
						},
					},
				},

				{
					applyHunks: []string{"337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125"},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125",
									Patch:     string(plusStartChunkDiff),
									IsApplied: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "dismiss chunk at the beginning",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusStartChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125",
									Patch: string(plusStartChunkDiff),
								},
							},
						},
					},
				},

				{
					dismissHunks: []string{"337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:          "337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125",
									Patch:       string(plusStartChunkDiff),
									IsDismissed: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "apply chunk in the middle",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusMiddleChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "22b8220ff705f16d21f05b04d993f69a9ea740ea058da43cf364d2e68480d831",
									Patch: string(plusMiddleChunkDiff),
								},
							},
						},
					},
				},
				{
					applyHunks: []string{"22b8220ff705f16d21f05b04d993f69a9ea740ea058da43cf364d2e68480d831"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "22b8220ff705f16d21f05b04d993f69a9ea740ea058da43cf364d2e68480d831",
									Patch:     string(plusMiddleChunkDiff),
									IsApplied: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "dismiss chunk in the middle",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusMiddleChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "22b8220ff705f16d21f05b04d993f69a9ea740ea058da43cf364d2e68480d831",
									Patch: string(plusMiddleChunkDiff),
								},
							},
						},
					},
				},
				{
					dismissHunks: []string{"22b8220ff705f16d21f05b04d993f69a9ea740ea058da43cf364d2e68480d831"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:          "22b8220ff705f16d21f05b04d993f69a9ea740ea058da43cf364d2e68480d831",
									Patch:       string(plusMiddleChunkDiff),
									IsDismissed: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "apply chunk in the end",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusEndChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "5a1176f7865aa965c54b82ba60c28e484cc65a5bee8f4bad137a630b195e85c9",
									Patch: string(plusEndChunkDiff),
								},
							},
						},
					},
				},
				{
					applyHunks: []string{"5a1176f7865aa965c54b82ba60c28e484cc65a5bee8f4bad137a630b195e85c9"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "5a1176f7865aa965c54b82ba60c28e484cc65a5bee8f4bad137a630b195e85c9",
									Patch:     string(plusEndChunkDiff),
									IsApplied: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "dismiss chunk in the end",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusEndChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "5a1176f7865aa965c54b82ba60c28e484cc65a5bee8f4bad137a630b195e85c9",
									Patch: string(plusEndChunkDiff),
								},
							},
						},
					},
				},
				{
					dismissHunks: []string{"5a1176f7865aa965c54b82ba60c28e484cc65a5bee8f4bad137a630b195e85c9"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:          "5a1176f7865aa965c54b82ba60c28e484cc65a5bee8f4bad137a630b195e85c9",
									Patch:       string(plusEndChunkDiff),
									IsDismissed: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "apply two chunks one by one",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusTwoChunks},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},
				{
					applyHunks: []string{"44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:     string(plusTwoChunksHunk1),
									IsApplied: true,
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},

				{
					applyHunks: []string{"03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:     string(plusTwoChunksHunk1),
									IsApplied: true,
								},
								{
									ID:        "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:     string(plusTwoChunksHunk2),
									IsApplied: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "apply two chunks one by one, backwards",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusTwoChunks},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},
				{
					applyHunks: []string{"03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:        "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:     string(plusTwoChunksHunk2),
									IsApplied: true,
								},
							},
						},
					},
				},
				{
					applyHunks: []string{"44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:     string(plusTwoChunksHunk1),
									IsApplied: true,
								},
								{
									ID:        "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:     string(plusTwoChunksHunk2),
									IsApplied: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "dismiss two chunks one by one",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusTwoChunks},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},
				{
					dismissHunks: []string{"44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:          "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:       string(plusTwoChunksHunk1),
									IsDismissed: true,
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},

				{
					dismissHunks: []string{"03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:          "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:       string(plusTwoChunksHunk1),
									IsDismissed: true,
								},
								{
									ID:          "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:       string(plusTwoChunksHunk2),
									IsDismissed: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "dismiss two chunks one by one backwards",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusTwoChunks},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},
				{
					dismissHunks: []string{"03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:          "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:       string(plusTwoChunksHunk2),
									IsDismissed: true,
								},
							},
						},
					},
				},
				{
					dismissHunks: []string{"44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:          "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:       string(plusTwoChunksHunk1),
									IsDismissed: true,
								},
								{
									ID:          "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:       string(plusTwoChunksHunk2),
									IsDismissed: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "apply two chunks",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusTwoChunks},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},
				{
					applyHunks: []string{"44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b", "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:     string(plusTwoChunksHunk1),
									IsApplied: true,
								},
								{
									ID:        "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:     string(plusTwoChunksHunk2),
									IsApplied: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "apply two chunks",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusTwoChunks},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},
				{
					applyHunks: []string{"44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b", "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:        "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:     string(plusTwoChunksHunk1),
									IsApplied: true,
								},
								{
									ID:        "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:     string(plusTwoChunksHunk2),
									IsApplied: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "dismiss two chunks",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusTwoChunks},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch: string(plusTwoChunksHunk1),
								},
								{
									ID:    "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch: string(plusTwoChunksHunk2),
								},
							},
						},
					},
				},
				{
					dismissHunks: []string{"44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b", "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837"},

					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:          "44da12f43fb46d324ba04f6fa7f4ea431f8e1102e6a74017259e00a33e63e69b",
									Patch:       string(plusTwoChunksHunk1),
									IsDismissed: true,
								},
								{
									ID:          "03bd6bd159d9e3fe531c073e39c24b8bcee6ded85eee933fa8018f5bd0c19837",
									Patch:       string(plusTwoChunksHunk2),
									IsDismissed: true,
								},
							},
						},
					},
				},
			},
		},

		{
			name: "outdated hunk suggestion",

			operations: []*operation{
				{writeOriginal: map[string][]byte{"file": original}},
				{
					suggesting: &diffs{
						write: map[string][]byte{"file": plusStartChunk},
					},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:    "337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125",
									Patch: string(plusStartChunkDiff),
								},
							},
						},
					},
				},
				{
					writeOriginal: map[string][]byte{"file": plusStartChunk},
					result: []unidiff.FileDiff{
						{
							OrigName:      "file",
							NewName:       "file",
							PreferredName: "file",
							Hunks: []unidiff.Hunk{
								{
									ID:         "337ae837f000b0c8b20bcd63df73f9062967056e3ed5d9506468dfaf993a8125",
									Patch:      string(plusStartChunkDiff),
									IsOutdated: true,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test := newTest(t, test.operations)
			test.run(t)
		})
	}
}

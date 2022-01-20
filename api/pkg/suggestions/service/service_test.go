package service_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"testing"

	"getsturdy.com/api/pkg/analytics/disabled"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	vcs_codebase "getsturdy.com/api/pkg/codebase/vcs"
	"getsturdy.com/api/pkg/internal/inmemory"
	"getsturdy.com/api/pkg/notification/sender"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/suggestions"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/view/events"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspace"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
	"getsturdy.com/api/vcs/testutil"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// todo:
// suggest new file
// new file on both sides
// suggest rename file
// suggest rename and edit file

var (
	//go:embed testdata/original.txt
	original []byte

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

type operation struct {
	applyHunks   []string
	dismissHunks []string
	dismiss      bool

	writeOriginal   *[]byte
	writeSuggesting *[]byte

	result []unidiff.FileDiff
}

func (o *operation) run(t *testing.T, test *test, suggestion *suggestions.Suggestion) {
	switch {
	case o.writeOriginal != nil:
		t.Logf("writing original")
		// create workspace
		if test.originalWorkspace == nil {
			originalWorkspace, err := test.workspaceService.Create(service_workspace.CreateWorkspaceRequest{
				UserID:     test.originalUserID,
				CodebaseID: test.codebaseID,
				Name:       fmt.Sprintf("%s's workspace", test.originalUserID),
			})
			assert.NoError(t, err)
			test.originalWorkspace = originalWorkspace
		}

		// connect workspace to a view
		test.originalViewID = fmt.Sprintf("%s-view", test.originalUserID)
		assert.NoError(t, test.viewDB.Create(view.View{
			ID:         test.originalViewID,
			UserID:     test.originalUserID,
			CodebaseID: test.codebaseID,
		}))
		test.originalWorkspace.ViewID = &test.originalViewID
		assert.NoError(t, test.workspaceDB.Update(test.originalWorkspace))
		vcs_view.Create(test.repoProvider, test.codebaseID, test.originalWorkspace.ID, test.originalViewID)

		// make some changes
		viewPath := test.repoProvider.ViewPath(test.codebaseID, test.originalViewID)
		assert.NoError(t, os.WriteFile(path.Join(viewPath, "file"), *o.writeOriginal, 0777))

		// take a workspace snapshot

		if test.originalWorkspace != nil {
			fmt.Printf("\noriginal: %+v\n\n", test.originalWorkspace.ID)
		}
		if test.suggestingWorkspace != nil {
			fmt.Printf("\nsuggesting: %+v\n\n", test.suggestingWorkspace.ID)
		}

		snapshot, err := test.gitSnapshotter.Snapshot(test.codebaseID, test.originalWorkspace.ID, snapshots.ActionViewSync, snapshotter.WithOnView(test.originalViewID))
		assert.NoError(t, err)
		test.originalWorkspace.LatestSnapshotID = &snapshot.ID
		assert.NoError(t, test.workspaceDB.Update(test.originalWorkspace))

		if suggestion != nil {
			result, err := test.suggestionService.Diffs(context.Background(), suggestion)
			if assert.NoError(t, err) {
				assert.Equal(t, o.result, result)
			}
		}
	case o.writeSuggesting != nil:
		t.Logf("writing suggestion")

		// start suggesting
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
		suggestingWorkspace.ViewID = &test.suggestingViewID
		assert.NoError(t, test.workspaceDB.Update(suggestingWorkspace))
		assert.NoError(t, vcs_view.Create(test.repoProvider, test.codebaseID, suggestingWorkspace.ID, test.suggestingViewID))

		// make some suggestions
		suggestingViewPath := test.repoProvider.ViewPath(test.codebaseID, test.suggestingViewID)
		assert.NoError(t, os.WriteFile(path.Join(suggestingViewPath, "file"), *o.writeSuggesting, 0777))

		// take a workspace snapshot
		suggestingSnapshot, err := test.gitSnapshotter.Snapshot(test.codebaseID, test.suggestingWorkspace.ID, snapshots.ActionViewSync, snapshotter.WithOnView(test.suggestingViewID))
		assert.NoError(t, err)
		test.suggestingWorkspace.LatestSnapshotID = &suggestingSnapshot.ID
		assert.NoError(t, test.workspaceDB.Update(test.suggestingWorkspace))

		result, err := test.suggestionService.Diffs(context.Background(), suggestion)
		if assert.NoError(t, err) {
			assert.Equal(t, o.result, result)
		}
	case o.applyHunks != nil:
		t.Logf("applying hunks")

		if assert.NoError(t, test.suggestionService.ApplyHunks(context.Background(), suggestion, o.applyHunks...)) {
			result, err := test.suggestionService.Diffs(context.Background(), suggestion)
			if assert.NoError(t, err) {
				assert.Equal(t, o.result, result)
			}
		}
	case o.dismissHunks != nil:
		t.Logf("dismissing hunks")
		if assert.NoError(t, test.suggestionService.DismissHunks(context.Background(), suggestion, o.dismissHunks...)) {
			result, err := test.suggestionService.Diffs(context.Background(), suggestion)
			if assert.NoError(t, err) {
				assert.Equal(t, o.result, result)
			}
		}
	case o.dismiss:
		t.Logf("dismissing sugestion")
		if assert.NoError(t, test.suggestionService.Dismiss(context.Background(), suggestion)) {
			result, err := test.suggestionService.Diffs(context.Background(), suggestion)
			if assert.NoError(t, err) {
				assert.Equal(t, o.result, result)
			}
		}
	}
}

type test struct {
	repoProvider      provider.RepoProvider
	executorProvider  executor.Provider
	suggestionRepo    db_suggestions.Repository
	viewDB            db_view.Repository
	workspaceDB       db_workspace.Repository
	snapshotsDB       db_snapshots.Repository
	codebaseUserRepo  db_codebase.CodebaseUserRepository
	gitSnapshotter    snapshotter.Snapshotter
	workspaceService  service_workspace.Service
	suggestionService *service_suggestions.Service

	codebaseID string

	originalUserID    string
	originalViewID    string
	originalWorkspace *workspace.Workspace

	suggestingUserID    string
	suggestingViewID    string
	suggestingWorkspace *workspace.Workspace
	suggestion          *suggestions.Suggestion

	operations []*operation
}

func newTest(t *testing.T, operations []*operation) *test {
	repoProvider := testutil.TestingRepoProvider(t)
	executorProvider := executor.NewProvider(zap.NewNop(), repoProvider)
	suggestionRepo := db_suggestions.NewMemory()

	viewDB := inmemory.NewInMemoryViewRepo()
	workspaceDB := db_workspace.NewMemory()
	snapshotsDB := inmemory.NewInMemorySnapshotRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	eventsSender := events.NewSender(codebaseUserRepo, workspaceDB, events.NewInMemory())

	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotsDB, workspaceDB, workspaceDB, viewDB, nil, executorProvider, zap.NewNop())
	workspaceService := service_workspace.New(zap.NewNop(), disabled.NewClient(), workspaceDB, workspaceDB, nil, nil, nil, nil, nil, executorProvider, nil, nil, gitSnapshotter, nil)
	suggestionService := service_suggestions.New(zap.NewNop(), suggestionRepo, workspaceService, executorProvider, gitSnapshotter, disabled.NewClient(), sender.NewNoopNotificationSender(), eventsSender)
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
	assert.NoError(t, vcs_codebase.Create(test.repoProvider, test.codebaseID))

	for _, operation := range test.operations {
		operation.run(t, test, test.suggestion)
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name       string
		operations []*operation
	}{
		{
			name: "apply add chunk at the beginning",

			operations: []*operation{
				{writeOriginal: &original},
				{
					writeSuggesting: &plusStartChunk,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusStartChunk,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusMiddleChunk,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusMiddleChunk,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusEndChunk,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusEndChunk,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusTwoChunks,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusTwoChunks,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusTwoChunks,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusTwoChunks,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusTwoChunks,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusTwoChunks,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusTwoChunks,
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
				{writeOriginal: &original},
				{
					writeSuggesting: &plusStartChunk,
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
					writeOriginal: &plusStartChunk,
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

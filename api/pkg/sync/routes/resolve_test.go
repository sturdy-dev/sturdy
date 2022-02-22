//nolint:bodyclose
package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"go.uber.org/dig"

	service_analytics "getsturdy.com/api/pkg/analytics/service"
	module_api "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebase"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	module_github "getsturdy.com/api/pkg/github/module"
	"getsturdy.com/api/pkg/graphql/resolvers"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	module_snapshots "getsturdy.com/api/pkg/snapshots/module"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/sync"
	routes_v3_sync "getsturdy.com/api/pkg/sync/routes"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	routes_v3_view "getsturdy.com/api/pkg/view/routes"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	vcsvcs "getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func str(s string) *string {
	return &s
}

func module(c *di.Container) {
	ctx := context.Background()
	c.Register(func() context.Context {
		return ctx
	})

	c.Import(module_api.Module)
	c.Import(module_configuration.TestingModule)
	c.Import(module_snapshots.TestingModule)

	// OSS version
	c.Import(module_github.Module)
}

func TestResolveHighLevelV2(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type nameContents struct {
		path     string
		contents *string
		copyFrom *string
	}

	cases := []struct {
		name                          string
		commonHistoryFiles            []nameContents
		trunkFiles                    []nameContents
		workspaceFiles                []nameContents // "unsaved changes" in the workspace
		expectedConflicts             bool
		customWorkspaceResolutions    []nameContents
		resolves                      []routes_v3_sync.ResolveFileRequest
		expectedContentsBeforeResolve []nameContents // data on disk when in the conflicting state
		expectedContentsAfterResolve  []nameContents
		tryToLandWithConflicts        bool
	}{
		{
			name:              "no-conflicts",
			trunkFiles:        []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
			workspaceFiles:    []nameContents{{path: "bar.txt", contents: str("bar-workspace")}},
			expectedConflicts: false,
			expectedContentsAfterResolve: []nameContents{
				{path: "foo.txt", contents: str("foo-trunk")},
				{path: "bar.txt", contents: str("bar-workspace")},
			},
		},

		{
			name:               "conflict-and-new-files",
			commonHistoryFiles: []nameContents{{path: "foo.txt", contents: str("foo-history")}},
			trunkFiles: []nameContents{
				{path: "foo.txt", contents: str("foo-trunk")},
				{path: "a.txt", contents: str("new-trunk")},
			},
			workspaceFiles: []nameContents{
				{path: "foo.txt", contents: str("foo-workspace")},
				{path: "b.txt", contents: str("new-workspace")},
			},
			expectedConflicts: true,
			resolves:          []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{
				{path: "foo.txt", contents: str("foo-workspace")},
				{path: "a.txt", contents: str("new-trunk")},
				{path: "b.txt", contents: str("new-workspace")},
			},
		},

		{
			name:               "conflict-and-new-files-try-to-land",
			commonHistoryFiles: []nameContents{{path: "foo.txt", contents: str("foo-history")}},
			trunkFiles: []nameContents{
				{path: "foo.txt", contents: str("foo-trunk")},
				{path: "a.txt", contents: str("new-trunk")},
			},
			workspaceFiles: []nameContents{
				{path: "foo.txt", contents: str("foo-workspace")},
				{path: "b.txt", contents: str("new-workspace")},
			},
			expectedConflicts: true,
			resolves:          []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{
				{path: "foo.txt", contents: str("foo-workspace")},
				{path: "a.txt", contents: str("new-trunk")},
				{path: "b.txt", contents: str("new-workspace")},
			},
			tryToLandWithConflicts: true,
		},

		{
			name:                         "no-changes",
			commonHistoryFiles:           []nameContents{},
			trunkFiles:                   []nameContents{{path: "README.md", contents: str("# Hello")}},
			workspaceFiles:               []nameContents{},
			expectedConflicts:            false,
			expectedContentsAfterResolve: []nameContents{{path: "README.md", contents: str("# Hello")}},
		},

		{
			name:                         "pick-workspace",
			trunkFiles:                   []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("foo-workspace")}},
			expectedConflicts:            true,
			resolves:                     []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt", contents: str("foo-workspace")}},
		},
		{
			name:                         "pick-trunk",
			trunkFiles:                   []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("foo-workspace")}},
			expectedConflicts:            true,
			resolves:                     []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
		},
		{
			name:                         "pick-custom",
			trunkFiles:                   []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("foo-workspace")}},
			expectedConflicts:            true,
			customWorkspaceResolutions:   []nameContents{{path: "foo.txt", contents: str("foo-custom")}},
			resolves:                     []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "custom"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt", contents: str("foo-custom")}},
		},
		{
			name: "pick-mixed-resolutions",
			trunkFiles: []nameContents{
				{path: "a.txt", contents: str("a-trunk")},
				{path: "b.txt", contents: str("b-trunk")},
				{path: "c.txt", contents: str("c-trunk")},
			},
			workspaceFiles: []nameContents{
				{path: "a.txt", contents: str("a-workspace")},
				{path: "b.txt", contents: str("b-workspace")},
				{path: "c.txt", contents: str("c-workspace")},
			},
			expectedConflicts:          true,
			customWorkspaceResolutions: []nameContents{{path: "c.txt", contents: str("c-custom")}},
			resolves: []routes_v3_sync.ResolveFileRequest{
				{FilePath: "a.txt", Version: "trunk"},
				{FilePath: "b.txt", Version: "workspace"},
				{FilePath: "c.txt", Version: "custom"},
			},
			expectedContentsAfterResolve: []nameContents{
				{path: "a.txt", contents: str("a-trunk")},
				{path: "b.txt", contents: str("b-workspace")},
				{path: "c.txt", contents: str("c-custom")},
			},
		},
		{
			name:                         "conflict-trunk-deleted-file-pick-workspace",
			commonHistoryFiles:           []nameContents{{path: "foo.txt", contents: str("common-history")}},
			trunkFiles:                   []nameContents{{path: "foo.txt"}}, // deleted on trunk
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("modified-workspace")}},
			expectedConflicts:            true,
			resolves:                     []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt", contents: str("modified-workspace")}},
		},
		{
			name:                         "conflict-trunk-deleted-file-pick-trunk",
			commonHistoryFiles:           []nameContents{{path: "foo.txt", contents: str("common-history")}},
			trunkFiles:                   []nameContents{{path: "foo.txt"}}, // deleted on trunk
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("modified-workspace")}},
			expectedConflicts:            true,
			resolves:                     []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt"}},
		},

		{
			name: "conflict-add-delete-extra-files",
			commonHistoryFiles: []nameContents{
				{path: "foo.txt", contents: str("common-history")},
				{path: "to-delete.txt", contents: str("to-be-deleted")},
			},
			trunkFiles:        []nameContents{{path: "foo.txt", contents: str("modified-trunk")}},
			workspaceFiles:    []nameContents{{path: "foo.txt", contents: str("modified-workspace")}},
			expectedConflicts: true,
			customWorkspaceResolutions: []nameContents{
				{path: "new.txt", contents: str("added")}, // previously untracked file
				{path: "to-delete.txt"},                   // delete this file, mid-sync
			},
			resolves: []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
			expectedContentsAfterResolve: []nameContents{
				{path: "foo.txt", contents: str("modified-trunk")},
				{path: "new.txt", contents: str("added")},
				{path: "to-delete.txt"},
			},
		},

		{
			name: "large-file-no-conflict", // test having a large file going through the syncing, no conflict in the large file
			commonHistoryFiles: []nameContents{
				{path: "foo.txt", contents: str("common-history")},
				{path: "large.jpg", copyFrom: str("testdata/large-img.jpg")},
			},
			trunkFiles:        []nameContents{{path: "foo.txt", contents: str("modified-trunk")}},
			workspaceFiles:    []nameContents{{path: "foo.txt", contents: str("modified-workspace")}},
			expectedConflicts: true,
			resolves:          []routes_v3_sync.ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
			expectedContentsBeforeResolve: []nameContents{
				{path: "foo.txt", contents: str(
					"<<<<<<< ........................................\n" +
						"modified-trunk\n" +
						"=======\n" +
						"modified-workspace\n" +
						">>>>>>> Unsaved workspace changes\n")},
				{path: "large.jpg", copyFrom: str("testdata/large-img.jpg")},
			},
			expectedContentsAfterResolve: []nameContents{
				{path: "foo.txt", contents: str("modified-trunk")},
				{path: "large.jpg", copyFrom: str("testdata/large-img.jpg")},
			},
		},
		{
			name: "large-file-with-conflict",
			commonHistoryFiles: []nameContents{
				{path: "foo.txt", contents: str("common-history")},
			},
			trunkFiles: []nameContents{
				{path: "large.jpg", copyFrom: str("testdata/large-img.jpg")},
			},
			workspaceFiles: []nameContents{
				{path: "large.jpg", copyFrom: str("testdata/large-img-2.jpg")},
			},
			expectedConflicts: true,
			resolves:          []routes_v3_sync.ResolveFileRequest{{FilePath: "large.jpg", Version: "workspace"}},
			expectedContentsBeforeResolve: []nameContents{
				// TODO: Today the LFS data is visible during conflicts in large files. Is this a good or bad idea?
				{path: "large.jpg", contents: str("version https://git-lfs.github.com/spec/v1\n" +
					"<<<<<<< ........................................\n" + // Regex matching
					"oid sha256:a540a47bdcc6f5af7cb8f1f1075d2d28848b97663502ca1cb3dfca2384361e6a\n" +
					"size 3911013\n" +
					"=======\n" +
					"oid sha256:2b60dee02b1eccfb66000973ff752ec6ffd8d5670b99174738948e7dd7ac71e6\n" +
					"size 3295872\n" +
					">>>>>>> Unsaved workspace changes")},
			},
			expectedContentsAfterResolve: []nameContents{
				{path: "large.jpg", copyFrom: str("testdata/large-img-2.jpg")},
			},
		},

		{
			name: "move-large-file-no-conflict",
			commonHistoryFiles: []nameContents{
				{path: "large.jpg", contents: str("testdata/large-img.jpg")},
			},
			trunkFiles: []nameContents{
				{path: "large.jpg"}, // delete
				{path: "large-moved.jpg", copyFrom: str("testdata/large-img.jpg")},
			},
			workspaceFiles: []nameContents{
				{path: "some-other.txt", contents: str("hello")},
			},
			expectedConflicts: false,
			expectedContentsAfterResolve: []nameContents{
				{path: "large-moved.jpg", copyFrom: str("testdata/large-img.jpg")},
				{path: "some-other.txt", contents: str("hello")},
			},
		},

		{
			name: "large-file-not-in-lfs",
			commonHistoryFiles: []nameContents{
				{path: "not-exists.jpg", contents: str(
					"version https://git-lfs.github.com/spec/v1\n" +
						"oid sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n" +
						"size 13337\n")},
			},
			trunkFiles: []nameContents{
				{path: "not-exists.jpg", contents: str(
					"version https://git-lfs.github.com/spec/v1\n" +
						"oid sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n" +
						"size 13337\n")},
			},
			workspaceFiles: []nameContents{
				{path: "not-exists.jpg", contents: str(
					"version https://git-lfs.github.com/spec/v1\n" +
						"oid sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\n" +
						"size 13337\n")},
			},
			expectedConflicts: true,
			resolves:          []routes_v3_sync.ResolveFileRequest{{FilePath: "not-exists.jpg", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{
				{path: "not-exists.jpg", contents: str(
					"version https://git-lfs.github.com/spec/v1\n" +
						"oid sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\n" +
						"size 13337\n")},
			},
		},
	}

	type deps struct {
		dig.In
		UserRepo              db_user.Repository
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		ViewRootResolver      resolvers.ViewRootResolver
		WorkspaceService      service_workspace.Service
		GitSnapshotter        snapshotter.Snapshotter
		RepoProvider          provider.RepoProvider
		CodebaseUserRepo      db_codebase.CodebaseUserRepository
		WorkspaceRepo         db_workspaces.Repository
		ViewRepo              db_view.Repository
		SnapshotRepo          db_snapshots.Repository
		ExecutorProvider      executor.Provider
		EventsSender          events.EventSender
		Logger                *zap.Logger
		AnalyticsService      *service_analytics.Service
		CodebaseRepo          db_codebase.CodebaseRepository
		SyncService           *service_sync.Service
	}

	var d deps
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	}

	userRepo := d.UserRepo
	workspaceRootResolver := d.WorkspaceRootResolver
	viewRootResolver := d.ViewRootResolver
	workspaceService := d.WorkspaceService
	repoProvider := d.RepoProvider
	codebaseUserRepo := d.CodebaseUserRepo
	codebaseRepo := d.CodebaseRepo

	createViewRoute := routes_v3_view.Create(d.Logger, d.ViewRepo, d.CodebaseUserRepo, d.AnalyticsService, d.WorkspaceRepo, d.GitSnapshotter, d.SnapshotRepo, d.WorkspaceRepo, d.ExecutorProvider, d.EventsSender)
	startRoutev2 := routes_v3_sync.StartV2(d.Logger, d.SyncService)
	resolveRoutev2 := routes_v3_sync.ResolveV2(d.Logger, d.SyncService)
	statusRoute := routes_v3_sync.Status(d.ViewRepo, d.ExecutorProvider, d.Logger)

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			userID := uuid.NewString()
			err := userRepo.Create(&users.User{ID: userID, Name: "Test Test", Email: userID + "@test.com"})
			assert.NoError(t, err)

			authenticatedUserContext := auth.NewContext(context.Background(), &auth.Subject{ID: userID, Type: auth.SubjectUser})

			codebaseID := uuid.NewString()

			err = codebaseRepo.Create(codebase.Codebase{
				ID:              codebaseID,
				ShortCodebaseID: codebase.ShortCodebaseID(codebaseID), // not realistic
			})
			assert.NoError(t, err)
			err = codebaseUserRepo.Create(codebase.CodebaseUser{
				ID:         uuid.NewString(),
				CodebaseID: codebaseID,
				UserID:     userID,
			})
			assert.NoError(t, err)

			basePath := repoProvider.TrunkPath(codebaseID)
			_, err = vcsvcs.CreateBareRepoWithRootCommit(basePath)
			assert.NoError(t, err)
			t.Logf("repo=%s", basePath)

			// Create workspace
			firstWorkspaceResolver, err := workspaceRootResolver.CreateWorkspace(
				authenticatedUserContext,
				resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{
					CodebaseID: graphql.ID(codebaseID),
				}},
			)
			assert.NoError(t, err)
			assert.NotNil(t, firstWorkspaceResolver)

			// Create view, attached to the first workspace
			var viewRes view.View
			requestWithParams(t, userID, createViewRoute, routes_v3_view.CreateRequest{
				CodebaseID:  codebaseID,
				WorkspaceID: string(firstWorkspaceResolver.ID()),
			}, &viewRes, nil)

			viewPath := repoProvider.ViewPath(codebaseID, viewRes.ID)

			makeChanges := func(changes []nameContents) {
				// Create change in workspace
				for _, f := range changes {
					if f.contents != nil {
						err = ioutil.WriteFile(path.Join(viewPath, f.path), []byte(*f.contents), 0o644)
					} else if f.copyFrom != nil {
						d, err := ioutil.ReadFile(*f.copyFrom)
						assert.NoError(t, err)
						err = ioutil.WriteFile(path.Join(viewPath, f.path), d, 0o644)
						assert.NoError(t, err)
					} else {
						err = os.Remove(path.Join(viewPath, f.path))
					}
					assert.NoError(t, err)
				}
			}

			makeAndLandChanges := func(changes []nameContents, workspaceID graphql.ID) error {
				makeChanges(changes)

				// Get diff
				diffs, _, err := workspaceService.Diffs(authenticatedUserContext, string(workspaceID))
				assert.NoError(t, err)

				var patchIds []string
				for _, d := range diffs {
					for _, h := range d.Hunks {
						patchIds = append(patchIds, h.ID)
					}
				}

				// Land the changes in the first workspace
				_, err = workspaceRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
					WorkspaceID: workspaceID,
					PatchIDs:    patchIds,
				}})
				return err
			}

			// Create common history (if any)
			if len(tc.commonHistoryFiles) > 0 {
				err = makeAndLandChanges(tc.commonHistoryFiles, firstWorkspaceResolver.ID())
				assert.NoError(t, err)
			}

			// Create second workspace
			secondWorkspaceResolver, err := workspaceRootResolver.CreateWorkspace(
				authenticatedUserContext,
				resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{
					CodebaseID: graphql.ID(codebaseID),
				}},
			)
			assert.NoError(t, err)

			// Create change in workspace
			err = makeAndLandChanges(tc.trunkFiles, firstWorkspaceResolver.ID())
			assert.NoError(t, err)

			// Open the second workspace
			_, err = viewRootResolver.OpenWorkspaceOnView(authenticatedUserContext, resolvers.OpenViewArgs{Input: resolvers.OpenWorkspaceOnViewInput{
				WorkspaceID: secondWorkspaceResolver.ID(),
				ViewID:      graphql.ID(viewRes.ID),
			}})
			assert.NoError(t, err)

			// Make the new changes
			if len(tc.workspaceFiles) > 0 {
				if tc.tryToLandWithConflicts {
					err = makeAndLandChanges(tc.workspaceFiles, secondWorkspaceResolver.ID())
					if tc.expectedConflicts {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				} else {
					makeChanges(tc.workspaceFiles)
				}
			}

			// Start sync
			viewIDParams := []gin.Param{{"viewID", viewRes.ID}}
			var startRebaseRes sync.RebaseStatusResponse
			requestWithParams(t, userID, startRoutev2, routes_v3_sync.InitSyncRequest{WorkspaceID: string(secondWorkspaceResolver.ID())}, &startRebaseRes, viewIDParams)
			assert.Equal(t, tc.expectedConflicts, startRebaseRes.HaveConflicts)

			// start sync again??
			{
				viewIDParams := []gin.Param{{"viewID", viewRes.ID}}
				var startRebaseRes sync.RebaseStatusResponse
				requestWithParams(t, userID, startRoutev2, routes_v3_sync.InitSyncRequest{WorkspaceID: string(secondWorkspaceResolver.ID())}, &startRebaseRes, viewIDParams)
				// assert.Equal(t, tc.expectedConflicts, startRebaseRes.HaveConflicts)
			}

			// Conflict resolution if we had a conflict
			if tc.expectedConflicts {
				assert.True(t, startRebaseRes.IsRebasing)
				assert.Len(t, startRebaseRes.ConflictingFiles, len(tc.resolves))

				// Verify contents before we resolve conflicts
				for _, f := range tc.expectedContentsBeforeResolve {
					d, err := ioutil.ReadFile(path.Join(viewPath, f.path))

					if f.contents != nil {
						assert.NoError(t, err)
						if !assert.Regexp(t, *f.contents, string(d)) {
							assert.Equal(t, *f.contents, string(d))
						}
					} else if f.copyFrom != nil {
						assert.NoError(t, err)
						expected, err := ioutil.ReadFile(*f.copyFrom)
						assert.NoError(t, err)
						// Compare len instead of full file, to be able to make sense of errors
						if !assert.Equal(t, len(expected), len(d)) {
							if len(d) < 1000 {
								t.Logf("contents: %s", string(d))
							}
						}

					} else {
						// expect to be deleted
						assert.Error(t, err)
						assert.Nil(t, d)
					}
				}

				// get status, expect to say that we're syncing
				{
					var syncStatusRes sync.RebaseStatusResponse
					requestWithParams(t, userID, statusRoute, struct{}{}, &syncStatusRes, viewIDParams)
					assert.True(t, syncStatusRes.IsRebasing)
					assert.True(t, syncStatusRes.HaveConflicts)
					assert.Len(t, syncStatusRes.ConflictingFiles, len(tc.resolves))
				}

				// Write custom resolutions
				for _, f := range tc.customWorkspaceResolutions {
					if f.contents != nil {
						err = ioutil.WriteFile(path.Join(viewPath, f.path), []byte(*f.contents), 0o644)
					} else if f.copyFrom != nil {
						d, err := ioutil.ReadFile(*f.copyFrom)
						assert.NoError(t, err)
						err = ioutil.WriteFile(path.Join(viewPath, f.path), d, 0o644)
						assert.NoError(t, err)
					} else {
						err = os.Remove(path.Join(viewPath, f.path))
					}
					assert.NoError(t, err)
				}

				// get status, expect to say that we're syncing
				{
					var syncStatusRes sync.RebaseStatusResponse
					requestWithParams(t, userID, statusRoute, struct{}{}, &syncStatusRes, viewIDParams)
					assert.True(t, syncStatusRes.IsRebasing)
					assert.True(t, syncStatusRes.HaveConflicts)
					assert.Len(t, syncStatusRes.ConflictingFiles, len(tc.resolves))
				}

				// Resolve conflict
				var resolveRebaseRes sync.RebaseStatusResponse
				requestWithParams(t, userID, resolveRoutev2, routes_v3_sync.ResolveRequest{Files: tc.resolves}, &resolveRebaseRes, viewIDParams)

				// The final resolve should leave leave state that is not rebasing and has no conflicts
				assert.False(t, resolveRebaseRes.IsRebasing)
				assert.False(t, resolveRebaseRes.HaveConflicts)
			}

			// Verify contents after resolve
			for _, f := range tc.expectedContentsAfterResolve {
				d, err := ioutil.ReadFile(path.Join(viewPath, f.path))

				if f.contents != nil {
					assert.NoError(t, err)
					if !assert.Regexp(t, *f.contents, string(d)) {
						assert.Equal(t, *f.contents, string(d))
					}
				} else if f.copyFrom != nil {
					assert.NoError(t, err)
					expected, err := ioutil.ReadFile(*f.copyFrom)
					assert.NoError(t, err)
					assert.Equal(t, len(expected), len(d)) // Compare len instead of full file, to be able to make sense of errors
				} else {
					// expect to be deleted
					assert.Error(t, err)
					assert.Nil(t, d)
				}
			}

			trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
			assert.NoError(t, err)
			trunkHead, err := trunkRepo.BranchCommitID("sturdytrunk")
			assert.NoError(t, err)
			workspaceHead, err := trunkRepo.BranchCommitID(string(secondWorkspaceResolver.ID()))
			assert.NoError(t, err)

			assert.Equal(t, trunkHead, workspaceHead, "workspace branch is not sturdytrunk")

			// TODO: Test landing
			// TODO: Test marked as conflicting
			// TODO: Test workspace head change
		})
	}
}

func requestWithParams(t *testing.T, userID string, route func(*gin.Context), request, response interface{}, params []gin.Param) {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Params = params

	data, err := json.Marshal(request)
	assert.NoError(t, err)

	c.Request, err = http.NewRequest("POST", "/", bytes.NewReader(data))
	c.Request = c.Request.WithContext(auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID}))
	assert.NoError(t, err)
	route(c)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	content, err := ioutil.ReadAll(res.Result().Body)
	assert.NoError(t, err)

	if len(content) > 0 {
		err = json.Unmarshal(content, response)
		assert.NoError(t, err)
	}
}

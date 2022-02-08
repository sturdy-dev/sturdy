package routes

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

	"getsturdy.com/api/pkg/analytics/disabled"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_change "getsturdy.com/api/pkg/change/db"
	service_change "getsturdy.com/api/pkg/change/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebase"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	db_comments "getsturdy.com/api/pkg/comments/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/internal/sturdytest"
	"getsturdy.com/api/pkg/queue"
	db_review "getsturdy.com/api/pkg/review/db"
	graphql_review "getsturdy.com/api/pkg/review/graphql"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	"getsturdy.com/api/pkg/sync"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	graphql_view "getsturdy.com/api/pkg/view/graphql"
	routes_v3_view "getsturdy.com/api/pkg/view/routes"
	db_activity "getsturdy.com/api/pkg/workspace/activity/db"
	activity_sender "getsturdy.com/api/pkg/workspace/activity/sender"
	service_activity "getsturdy.com/api/pkg/workspace/activity/service"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	graphql_workspace "getsturdy.com/api/pkg/workspace/graphql"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	db_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/db"
	graphql_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/graphql"
	service_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/service"
	vcsvcs "getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	allFilesAllowed, _ = unidiff.NewAllower("*")
)

func str(s string) *string {
	return &s
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
		resolves                      []ResolveFileRequest
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
			resolves:          []ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
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
			resolves:          []ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
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
			resolves:                     []ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt", contents: str("foo-workspace")}},
		},
		{
			name:                         "pick-trunk",
			trunkFiles:                   []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("foo-workspace")}},
			expectedConflicts:            true,
			resolves:                     []ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
		},
		{
			name:                         "pick-custom",
			trunkFiles:                   []nameContents{{path: "foo.txt", contents: str("foo-trunk")}},
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("foo-workspace")}},
			expectedConflicts:            true,
			customWorkspaceResolutions:   []nameContents{{path: "foo.txt", contents: str("foo-custom")}},
			resolves:                     []ResolveFileRequest{{FilePath: "foo.txt", Version: "custom"}},
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
			resolves: []ResolveFileRequest{
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
			resolves:                     []ResolveFileRequest{{FilePath: "foo.txt", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{{path: "foo.txt", contents: str("modified-workspace")}},
		},
		{
			name:                         "conflict-trunk-deleted-file-pick-trunk",
			commonHistoryFiles:           []nameContents{{path: "foo.txt", contents: str("common-history")}},
			trunkFiles:                   []nameContents{{path: "foo.txt"}}, // deleted on trunk
			workspaceFiles:               []nameContents{{path: "foo.txt", contents: str("modified-workspace")}},
			expectedConflicts:            true,
			resolves:                     []ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
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
			resolves: []ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
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
			resolves:          []ResolveFileRequest{{FilePath: "foo.txt", Version: "trunk"}},
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
			resolves:          []ResolveFileRequest{{FilePath: "large.jpg", Version: "workspace"}},
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
			resolves:          []ResolveFileRequest{{FilePath: "not-exists.jpg", Version: "workspace"}},
			expectedContentsAfterResolve: []nameContents{
				{path: "not-exists.jpg", contents: str(
					"version https://git-lfs.github.com/spec/v1\n" +
						"oid sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\n" +
						"size 13337\n")},
			},
		},
	}
	reposBasePath, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)
	repoProvider := provider.New(reposBasePath, "localhost:8888")

	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewDevelopment()
	viewRepo := db_view.NewRepo(d)
	viewUpdates := events.NewInMemory()
	userRepo := db_user.NewRepo(d)
	codebaseRepo := db_codebase.NewRepo(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	workspaceRepo := db_workspace.NewRepo(d)
	snapshotsRepo := db_snapshots.NewRepo(d)

	executorProvider := executor.NewProvider(logger, repoProvider)

	changeRepo := db_change.NewRepo(d)
	changeCommitRepo := db_change.NewCommitRepository(d)

	reviewRepo := db_review.NewReviewRepository(d)
	viewEvents := events.NewInMemory()
	workspaceActivityRepo := db_activity.NewActivityRepo(d)
	workspaceActivityReadsRepo := db_activity.NewActivityReadsRepo(d)
	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, viewUpdates)
	activityService := service_activity.New(workspaceActivityReadsRepo, eventsSender)
	activitySender := activity_sender.NewActivitySender(codebaseUserRepo, workspaceActivityRepo, activityService, eventsSender)
	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotsRepo, workspaceRepo, workspaceRepo, viewRepo, eventsSender, executorProvider, logger)
	commentRepo := db_comments.NewRepo(d)
	queue := queue.NewNoop()
	snapshotPublisher := worker_snapshots.New(logger, queue, gitSnapshotter)
	commentsService := service_comments.New(commentRepo)

	buildQueue := workers_ci.New(zap.NewNop(), queue, nil)

	changeService := service_change.New(nil, userRepo, changeRepo, changeCommitRepo)
	workspaceService := service_workspace.New(
		logger,
		disabled.NewClient(),

		workspaceRepo,
		workspaceRepo,

		userRepo,
		reviewRepo,

		commentsService,
		changeService,

		activitySender,
		executorProvider,
		eventsSender,
		snapshotPublisher,
		gitSnapshotter,
		buildQueue,
	)

	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, workspaceService, nil, logger, executorProvider, nil, eventsSender)
	authService := service_auth.New(codebaseService, nil, workspaceService, nil, nil)

	workspaceWatchersRootResolver := new(resolvers.WorkspaceWatcherRootResolver)
	workspaceWatcherRepo := db_workspace_watchers.NewInMemory()
	workspaceWatchersService := service_workspace_watchers.New(workspaceWatcherRepo, eventsSender)

	reviewRootResolver := graphql_review.New(
		logger,
		reviewRepo,
		nil,
		authService,
		nil,
		nil,
		eventsSender,
		viewEvents,
		nil,
		activitySender,
		workspaceWatchersService,
	)

	workspaceRootResolver := graphql_workspace.NewResolver(
		workspaceRepo,
		codebaseRepo,
		viewRepo,
		nil, // commentRepo
		nil, // snapshotRepo
		nil, // codebaseResolver
		nil, // authorResolver
		nil, // viewResolver
		nil, // commentResolver
		nil, // prResolver
		nil, // changeResolver
		nil, // workspaceActivityResolver
		reviewRootResolver,
		nil, // presenceRootResolver
		nil, // suggestitonsRootResolver
		nil, // statusesRootResolver
		*workspaceWatchersRootResolver,
		nil, // suggestionsService
		workspaceService,
		authService,
		logger,
		viewUpdates,
		workspaceRepo,
		executorProvider,
		eventsSender,
		gitSnapshotter,
	)

	*workspaceWatchersRootResolver = graphql_workspace_watchers.NewRootResolver(
		logger,
		workspaceWatchersService,
		workspaceService,
		authService,
		viewEvents,
		nil,
		&workspaceRootResolver,
	)

	viewRootResolver := graphql_view.NewResolver(
		viewRepo,
		workspaceRepo,
		gitSnapshotter,
		snapshotsRepo,
		nil,
		nil,
		workspaceRepo,
		viewEvents,
		eventsSender,
		executorProvider,
		logger,
		nil,
		workspaceWatchersService,
		nil,
		nil,
		authService,
	)

	syncService := service_sync.New(logger, executorProvider, viewRepo, workspaceRepo, workspaceRepo, gitSnapshotter)

	startRoutev2 := StartV2(logger, syncService)
	resolveRoutev2 := ResolveV2(logger, syncService)
	statusRoute := Status(viewRepo, executorProvider, logger)

	createViewRoute := routes_v3_view.Create(
		logger,
		viewRepo,
		codebaseUserRepo,
		disabled.NewClient(),
		workspaceRepo,
		gitSnapshotter,
		snapshotsRepo,
		workspaceRepo,
		executorProvider,
		eventsSender,
	)

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			userID := uuid.NewString()
			err = userRepo.Create(&users.User{ID: userID, Name: "Test Test", Email: userID + "@test.com"})
			assert.NoError(t, err)

			authenticatedUserContext := auth.NewContext(context.Background(), &auth.Subject{ID: userID, Type: auth.SubjectUser})

			codebaseID := uuid.NewString()

			err := codebaseRepo.Create(codebase.Codebase{
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
			requestWithParams(t, userID, startRoutev2, InitSyncRequest{WorkspaceID: string(secondWorkspaceResolver.ID())}, &startRebaseRes, viewIDParams)
			assert.Equal(t, tc.expectedConflicts, startRebaseRes.HaveConflicts)

			// start sync again??
			{
				viewIDParams := []gin.Param{{"viewID", viewRes.ID}}
				var startRebaseRes sync.RebaseStatusResponse
				requestWithParams(t, userID, startRoutev2, InitSyncRequest{WorkspaceID: string(secondWorkspaceResolver.ID())}, &startRebaseRes, viewIDParams)
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

				// Trigger snapshot! (just because, it usually happens when new files are written)
				assert.NoError(t, snapshotPublisher.Enqueue(context.Background(), codebaseID, viewRes.ID, viewRes.WorkspaceID, []string{"."}, snapshots.ActionViewSync))

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
					} else {
						err = os.Remove(path.Join(viewPath, f.path))
					}
					assert.NoError(t, err)
				}

				// Trigger snapshot! (just because, it usually happens when new files are written)
				assert.NoError(t, snapshotPublisher.Enqueue(context.Background(), codebaseID, viewRes.ID, viewRes.WorkspaceID, []string{"."}, snapshots.ActionViewSync))

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
				requestWithParams(t, userID, resolveRoutev2, ResolveRequest{Files: tc.resolves}, &resolveRebaseRes, viewIDParams)

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

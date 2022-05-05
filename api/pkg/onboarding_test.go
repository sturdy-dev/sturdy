//nolint:bodyclose
package pkg_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"go.uber.org/dig"

	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/auth"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	routes_v3_codebase "getsturdy.com/api/pkg/codebases/routes"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/di"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	service_gc "getsturdy.com/api/pkg/gc/service"
	gqldataloader "getsturdy.com/api/pkg/graphql/dataloader"
	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/internal/dbtest"
	queue "getsturdy.com/api/pkg/queue/module"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	workers_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	service_users "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	routes_v3_view "getsturdy.com/api/pkg/view/routes"
	service_view "getsturdy.com/api/pkg/view/service"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	routes_v3_workspace "getsturdy.com/api/pkg/workspaces/routes"
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

func testModule(t *testing.T) di.Module {
	return func(c *di.Container) {
		c.Import(api.Module)
		c.ImportWithForce(configuration.TestModule)
		c.ImportWithForce(queue.TestModule)
		c.Register(func() *testing.T { return t })
		c.RegisterWithForce(dbtest.DB)
	}
}

func TestCreate(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		di.In
		UserRepo              db_user.Repository
		UserService           service_users.Service
		CodebaseRootResolver  resolvers.CodebaseRootResolver
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		UserRootResolver      resolvers.UserRootResolver
		CommentsRootResolver  resolvers.CommentRootResolver
		ViewRootResolver      resolvers.ViewRootResolver
		LandRootResolver      resolvers.LandRootResovler
		GcService             *service_gc.Service
		CodebaseService       *service_codebase.Service
		WorkspaceService      *service_workspace.Service
		GitSnapshotter        *service_snapshots.Service
		RepoProvider          provider.RepoProvider
		SnapshotsQueue        workers_snapshots.Queue
		CIQueue               *workers_ci.BuildQueue

		// Dependencies of Gin Routes
		CodebaseUserRepo db_codebases.CodebaseUserRepository
		WorkspaceRepo    db_workspaces.Repository
		ViewRepo         db_view.Repository
		SnapshotRepo     db_snapshots.Repository
		ExecutorProvider executor.Provider
		EventsSender     *eventsv2.Publisher
		ViewService      *service_view.Service

		Logger           *zap.Logger
		AnalyticsService *service_analytics.Service
	}

	var d deps
	if !assert.NoError(t, di.Init(testModule(t)).To(&d)) {
		t.FailNow()
	}

	go func() {
		assert.NoError(t, d.SnapshotsQueue.Start(context.TODO()))
	}()
	go func() {
		assert.NoError(t, d.CIQueue.Start(context.TODO()))
	}()

	userRepo := d.UserRepo
	codebaseRootResolver := d.CodebaseRootResolver
	workspaceRootResolver := d.WorkspaceRootResolver
	landRootResolver := d.LandRootResolver
	userRootResolver := d.UserRootResolver
	commentsRootResolver := d.CommentsRootResolver
	viewRootResolver := d.ViewRootResolver
	gcService := d.GcService
	codebaseService := d.CodebaseService
	workspaceService := d.WorkspaceService
	repoProvider := d.RepoProvider

	logger := d.Logger
	analyticsService := d.AnalyticsService
	codebaseUserRepo := d.CodebaseUserRepo
	workspaceRepo := d.WorkspaceRepo
	viewRepo := d.ViewRepo
	executorProvider := d.ExecutorProvider

	createCodebaseRoute := routes_v3_codebase.Create(logger, codebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(logger, workspaceService, codebaseUserRepo)
	createViewRoute := routes_v3_view.Create(logger, viewRepo, codebaseUserRepo, analyticsService, workspaceRepo, executorProvider, d.ViewService)

	createUser := users.User{ID: users.ID(uuid.New().String()), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
	assert.NoError(t, userRepo.Create(&createUser))

	authenticatedUserContext := gqldataloader.NewContext(auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: createUser.ID.String()}))

	// Create a codebase
	var codebaseRes codebases.Codebase
	request(t, createUser.ID, createCodebaseRoute, routes_v3_codebase.CreateRequest{Name: "testrepo"}, &codebaseRes)
	assert.Len(t, codebaseRes.ID, 36)
	assert.Equal(t, "testrepo", codebaseRes.Name)
	assert.True(t, codebaseRes.IsReady, "codebase is ready")

	// Create a workspace
	var workspaceResult workspaces.Workspace
	request(t, createUser.ID, createWorkspaceRoute, routes_v3_workspace.CreateRequest{
		CodebaseID: codebaseRes.ID,
	}, &workspaceResult)
	assert.Len(t, workspaceResult.ID, 36)
	workspaceID := workspaceResult.ID

	// Create a view
	var viewRes view.View
	request(t, createUser.ID, createViewRoute, routes_v3_view.CreateRequest{
		CodebaseID:    codebaseRes.ID,
		WorkspaceID:   workspaceID,
		MountPath:     "~/testing",
		MountHostname: "testing.ftw",
	}, &viewRes)
	assert.Len(t, viewRes.ID, 36)
	assert.True(t, viewRes.CreatedAt.After(time.Now().Add(time.Second*-5)))

	// Make more changes to test.txt
	viewPath := repoProvider.ViewPath(codebaseRes.ID, viewRes.ID)

	t.Logf("viewPath=%s", viewPath)

	// Make changes in the view
	err := ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte("hello\n"), 0o666)
	assert.NoError(t, err)

	// Get diff
	diffs, _, err := workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	expectedDiffs := []unidiff.FileDiff{{OrigName: "/dev/null", NewName: "test.txt", PreferredName: "test.txt", IsNew: true, Hunks: []unidiff.Hunk{
		{
			ID:    "edc5f8dc6b69a14eefbdc56d830c44faf08d41ea6a370f4e0252b02906946991",
			Patch: "diff --git /dev/null \"b/test.txt\"\nnew file mode 100644\nindex 0000000..ce01362\n--- /dev/null\n+++ \"b/test.txt\"\n@@ -0,0 +1,1 @@\n+hello\n",
		},
	}}}
	assert.Equal(t, expectedDiffs, diffs)

	// Set workspace draft description
	_, err = workspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
		ID:               graphql.ID(workspaceID),
		DraftDescription: str("This is my first change"),
	}})
	assert.NoError(t, err)

	// Apply and land
	_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(workspaceID),
	}})
	assert.NoError(t, err)

	viewResolver, err := viewRootResolver.View(authenticatedUserContext, resolvers.ViewArgs{
		ID: graphql.ID(viewRes.ID),
	})
	assert.NoError(t, err)

	// get a reference to the new workspace after the change has been applied
	newWorkspace, err := viewResolver.Workspace(authenticatedUserContext)
	assert.NoError(t, err)
	workspaceID = string(newWorkspace.ID())

	// Get changelog in codebase
	cid := graphql.ID(codebaseRes.ID)
	codebaseResolver, err := codebaseRootResolver.Codebase(authenticatedUserContext, resolvers.CodebaseArgs{ID: &cid})
	assert.NoError(t, err)
	changes, err := codebaseResolver.Changes(authenticatedUserContext, nil)
	assert.NoError(t, err)
	if assert.Len(t, changes, 1) {
		assert.Equal(t, "This is my first change", changes[0].Description())
		author, err := changes[0].Author(authenticatedUserContext)
		assert.NoError(t, err)
		assert.Equal(t, createUser.ID, users.ID(author.ID()))
	}

	err = ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte("a\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\n"), 0o666)
	assert.NoError(t, err)

	// Get diff
	diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	expectedDiffs = []unidiff.FileDiff{{OrigName: "test.txt", NewName: "test.txt", PreferredName: "test.txt", Hunks: []unidiff.Hunk{
		{
			ID:    "fc85a16f432f111d2fc38572a4207c28547b03efcc629aabbd96021d773d9460",
			Patch: "diff --git \"a/test.txt\" \"b/test.txt\"\nindex ce01362..0edb856 100644\n--- \"a/test.txt\"\n+++ \"b/test.txt\"\n@@ -1,1 +1,26 @@\n-hello\n+a\n+b\n+c\n+d\n+e\n+f\n+g\n+h\n+i\n+j\n+k\n+l\n+m\n+n\n+o\n+p\n+q\n+r\n+s\n+t\n+u\n+v\n+w\n+x\n+y\n+z\n",
		},
	}}}
	assert.Equal(t, expectedDiffs, diffs)

	// Check if conflicts
	{
		ws, err := workspaceService.GetByID(context.Background(), workspaceID)
		assert.NoError(t, err)
		conflicts, err := workspaceService.HasConflicts(context.Background(), ws)
		assert.NoError(t, err)
		assert.False(t, conflicts)
	}

	// Apply and land
	_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(workspaceID),
	}})
	assert.NoError(t, err)

	newWorkspace, err = viewResolver.Workspace(authenticatedUserContext)
	assert.NoError(t, err)
	workspaceID = string(newWorkspace.ID())

	// Make changes to two parts of the file (early and late), expect two hunks
	// The row "d" is deleted, and "t" is replaced with "ttt"
	err = ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte("a\nb\nc\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nttt\nu\nv\nw\nx\ny\nz\n"), 0o666)
	assert.NoError(t, err)

	// Get diff
	diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	expectedDiffs = []unidiff.FileDiff{{OrigName: "test.txt", NewName: "test.txt", PreferredName: "test.txt",
		Hunks: []unidiff.Hunk{
			{ID: "9e8e97e972ee7e13b80776480da86335e0c8635d675fb446a216c1aa40ece79e", Patch: "diff --git \"a/test.txt\" \"b/test.txt\"\nindex 0edb856..9389e12 100644\n--- \"a/test.txt\"\n+++ \"b/test.txt\"\n@@ -1,7 +1,6 @@\n a\n b\n c\n-d\n e\n f\n g\n"},
			{ID: "7b6e4538c0b1c2ffe0a38164ee1be3b6e547c7cacef149efaca9be8241f0b60c", Patch: "diff --git \"a/test.txt\" \"b/test.txt\"\nindex 0edb856..9389e12 100644\n--- \"a/test.txt\"\n+++ \"b/test.txt\"\n@@ -17,7 +16,7 @@ p\n q\n r\n s\n-t\n+ttt\n u\n v\n w\n"},
		}}}
	assert.Equal(t, expectedDiffs, diffs)

	// Undo the second hunk
	_, err = workspaceRootResolver.RemovePatches(authenticatedUserContext, resolvers.RemovePatchesArgs{Input: resolvers.RemovePatchesInput{
		WorkspaceID: graphql.ID(workspaceID),
		HunkIDs:     []string{diffs[0].Hunks[1].ID},
	}})
	assert.NoError(t, err)

	// Get diff
	diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	expectedDiffs = []unidiff.FileDiff{{OrigName: "test.txt", NewName: "test.txt", PreferredName: "test.txt",
		Hunks: []unidiff.Hunk{
			{ID: "00755ee69c4365ed7304f1e1bc515cf5fef3e22cd89a28c15e635c7faae7888c", Patch: "diff --git \"a/test.txt\" \"b/test.txt\"\nindex 0edb856..215f140 100644\n--- \"a/test.txt\"\n+++ \"b/test.txt\"\n@@ -1,7 +1,6 @@\n a\n b\n c\n-d\n e\n f\n g\n"},
		}}}
	assert.Equal(t, expectedDiffs, diffs)

	// Edit the file so that there are 3 hunks
	err = ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte("aaaa\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nlll\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nzzz\n"), 0o666)
	assert.NoError(t, err)

	diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	expectedDiffs = []unidiff.FileDiff{{OrigName: "test.txt", NewName: "test.txt", PreferredName: "test.txt",
		Hunks: []unidiff.Hunk{
			{ID: "7e412eedbb31eb4a13695ee490fd5a4fe39f6f33611fe77a5feabf2ffa4ed8d0", Patch: "diff --git \"a/test.txt\" \"b/test.txt\"\nindex 0edb856..da65dab 100644\n--- \"a/test.txt\"\n+++ \"b/test.txt\"\n@@ -1,4 +1,4 @@\n-a\n+aaaa\n b\n c\n d\n"},
			{ID: "b500f195ce7c53ad9324bcab5065393858f85675ad760b888270a32f2fd82345", Patch: "diff --git \"a/test.txt\" \"b/test.txt\"\nindex 0edb856..da65dab 100644\n--- \"a/test.txt\"\n+++ \"b/test.txt\"\n@@ -9,7 +9,7 @@ h\n i\n j\n k\n-l\n+lll\n m\n n\n o\n"},
			{ID: "98ce5f04a5f07d7faf61d479ee51a8876c43d738d74403753f442b762cf5942d", Patch: "diff --git \"a/test.txt\" \"b/test.txt\"\nindex 0edb856..da65dab 100644\n--- \"a/test.txt\"\n+++ \"b/test.txt\"\n@@ -23,4 +23,4 @@ v\n w\n x\n y\n-z\n+zzz\n"},
		}}}
	assert.Equal(t, expectedDiffs, diffs)

	// Move the file
	err = os.Rename(path.Join(viewPath, "test.txt"), path.Join(viewPath, "test-2.txt"))
	assert.NoError(t, err)

	// Get diff
	diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	expectedDiffs = []unidiff.FileDiff{{OrigName: "test.txt", NewName: "test-2.txt", PreferredName: "test-2.txt", IsMoved: true,
		Hunks: []unidiff.Hunk{
			{ID: "24bf7f7b8adff226351e7e836e057de609b1ed8b8468994e29da7b4ea35f5a9b", Patch: "diff --git \"a/test.txt\" \"b/test-2.txt\"\nsimilarity index 88%\nrename from \"test.txt\"\nrename to \"test-2.txt\"\nindex 0edb856..da65dab 100644\n--- \"a/test.txt\"\n+++ \"b/test-2.txt\"\n@@ -1,4 +1,4 @@\n-a\n+aaaa\n b\n c\n d\n"},
			{ID: "d47efdcc630a8132e06bfb983274e8a2e0be8730cfbc50b7282655c34eb0574c", Patch: "diff --git \"a/test.txt\" \"b/test-2.txt\"\nsimilarity index 88%\nrename from \"test.txt\"\nrename to \"test-2.txt\"\nindex 0edb856..da65dab 100644\n--- \"a/test.txt\"\n+++ \"b/test-2.txt\"\n@@ -9,7 +9,7 @@ h\n i\n j\n k\n-l\n+lll\n m\n n\n o\n"},
			{ID: "14cd042cd4c5de65164c85c8865cdddd18c7ce9afd3e93ae2ac6f50f7647a782", Patch: "diff --git \"a/test.txt\" \"b/test-2.txt\"\nsimilarity index 88%\nrename from \"test.txt\"\nrename to \"test-2.txt\"\nindex 0edb856..da65dab 100644\n--- \"a/test.txt\"\n+++ \"b/test-2.txt\"\n@@ -23,4 +23,4 @@ v\n w\n x\n y\n-z\n+zzz\n"},
		}}}
	assert.Equal(t, expectedDiffs, diffs)

	// Apply and land
	_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(workspaceID),
	}})
	assert.NoError(t, err)

	newWorkspace, err = viewResolver.Workspace(authenticatedUserContext)
	assert.NoError(t, err)
	workspaceID = string(newWorkspace.ID())

	// Move file without edits
	err = os.Rename(path.Join(viewPath, "test-2.txt"), path.Join(viewPath, "test-3.txt"))
	assert.NoError(t, err)

	diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	expectedDiffs = []unidiff.FileDiff{{OrigName: "test-2.txt", NewName: "test-3.txt", PreferredName: "test-3.txt", IsMoved: true,
		Hunks: []unidiff.Hunk{
			{ID: "b477342d7b12a211ec83fbdc9bf9fb259903046c1ed76050683b503eb36ae69d", Patch: "diff --git \"a/test-2.txt\" \"b/test-3.txt\"\nsimilarity index 100%\nrename from \"test-2.txt\"\nrename to \"test-3.txt\"\n"},
		}}}
	assert.Equal(t, expectedDiffs, diffs)

	// Apply and land
	_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(workspaceID),
	}})
	assert.NoError(t, err)

	newWorkspace, err = viewResolver.Workspace(authenticatedUserContext)
	assert.NoError(t, err)
	workspaceID = string(newWorkspace.ID())

	// Make changes with conflicts, attempting to land should fail gracefully
	// Create a workspace
	var secondWorkspaceRes workspaces.Workspace
	request(t, createUser.ID, createWorkspaceRoute, routes_v3_workspace.CreateRequest{
		CodebaseID: codebaseRes.ID,
	}, &secondWorkspaceRes)
	assert.Len(t, secondWorkspaceRes.ID, 36)

	// Make a change in the first workspace (it's still checked out)
	err = ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte("aaaaa\n"), 0o666)
	assert.NoError(t, err)

	// Get diff
	diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
	assert.NoError(t, err)
	assert.Len(t, diffs, 1)
	assert.Len(t, diffs[0].Hunks, 1)

	// Apply and land
	_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(workspaceID),
	}})
	assert.NoError(t, err)

	newWorkspace, err = viewResolver.Workspace(authenticatedUserContext)
	assert.NoError(t, err)
	workspaceID = string(newWorkspace.ID())

	// Checkout the new workspace
	_, err = viewRootResolver.OpenWorkspaceOnView(authenticatedUserContext, resolvers.OpenViewArgs{Input: resolvers.OpenWorkspaceOnViewInput{
		WorkspaceID: graphql.ID(secondWorkspaceRes.ID),
		ViewID:      graphql.ID(viewRes.ID),
	}})
	assert.NoError(t, err)

	// make changes in the second workspace and try to land it
	err = ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte("bbbbb\n"), 0o666)
	assert.NoError(t, err)

	// Get diff
	diffs, _, err = workspaceService.Diffs(context.Background(), secondWorkspaceRes.ID)
	assert.NoError(t, err)
	assert.Len(t, diffs, 1)
	assert.Len(t, diffs[0].Hunks, 1)

	// Apply and land, this should fail!
	_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(secondWorkspaceRes.ID),
	}})
	assert.Error(t, err)

	// The diffs should not have changed (no change should have been created)
	diffsAfterFailedLand, _, err := workspaceService.Diffs(context.Background(), secondWorkspaceRes.ID)
	assert.NoError(t, err)
	assert.Equal(t, diffs, diffsAfterFailedLand)

	// Switch to the original workspace
	_, err = viewRootResolver.OpenWorkspaceOnView(authenticatedUserContext, resolvers.OpenViewArgs{Input: resolvers.OpenWorkspaceOnViewInput{
		WorkspaceID: graphql.ID(workspaceID),
		ViewID:      graphql.ID(viewRes.ID),
	}})
	assert.NoError(t, err, errors.Unwrap(err))

	var contents = []string{
		"this\nis\na\nfile\naaaaaa",
		"this\nis\na\nfile\naaaaaa\n",
		"this\nis\na\nfile\naaaaaa",
		"this\nis\na\nfile\naaaaaa\n",
		"this\nis\na\nfile\naaaaaa",
		"this\r\nis\r\na\r\nfile\r\naaaaaa\r\n",
		"this\r\nis\r\na\r\nfile\r\naaaaaa",
		"this\r\nis\r\na\r\nfile\r\naaaaaa\r\n",
	}

	for _, cont := range contents {
		// Remove the trailing newline in test.txt
		err = ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte(cont), 0o666)
		assert.NoError(t, err)

		// Get diff
		diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
		assert.NoError(t, err)
		t.Logf("diffs=%+v", diffs)
		assert.Len(t, diffs, 1)
		assert.Len(t, diffs[0].Hunks, 1)

		// Apply and land
		_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
			WorkspaceID: graphql.ID(workspaceID),
		}})
		assert.NoError(t, err)

		newWorkspace, err = viewResolver.Workspace(authenticatedUserContext)
		assert.NoError(t, err)
		workspaceID = string(newWorkspace.ID())
	}

	// List views for user
	userResolver, err := userRootResolver.User(authenticatedUserContext)
	assert.NoError(t, err)
	allUserViews, err := userResolver.Views()
	assert.NoError(t, err)
	assert.Len(t, allUserViews, 1)

	// Make a comment on live changes in a workspace
	err = ioutil.WriteFile(path.Join(viewPath, "file-a.txt"), []byte("Hello World\n"), 0o666)
	assert.NoError(t, err)

	createdCommentResolver, err := commentsRootResolver.CreateComment(authenticatedUserContext, resolvers.CreateCommentArgs{Input: resolvers.CreateCommentInput{
		Message:     "Comment!",
		InReplyTo:   nil,
		Path:        str("file-a.txt"),
		LineStart:   i(1),
		LineEnd:     i(1),
		LineIsNew:   b(true),
		ChangeID:    nil,
		WorkspaceID: gid(graphql.ID(workspaceID)),
		ViewID:      gid(graphql.ID(viewRes.ID)),
	}})
	if assert.NoError(t, err, errors.Unwrap(err)) {
		assert.Equal(t, "Comment!", createdCommentResolver.Message())
	}

	// Get comment from workspace
	{
		getWorkspaceResolver, err := workspaceRootResolver.Workspace(authenticatedUserContext, resolvers.WorkspaceArgs{ID: graphql.ID(workspaceID)})
		assert.NoError(t, err)
		topComments, err := getWorkspaceResolver.Comments()
		assert.NoError(t, err)
		if assert.Len(t, topComments, 1) {
			assert.Equal(t, "Comment!", topComments[0].Message())
			codeContext := topComments[0].CodeContext()
			assert.Equal(t, int32(1), codeContext.LineStart())
			assert.Equal(t, int32(1), codeContext.LineEnd())
			assert.Equal(t, true, codeContext.LineIsNew())
			assert.Equal(t, "file-a.txt", codeContext.Path())
		}
	}

	// Move the file with the comment in it
	err = os.Rename(
		path.Join(viewPath, "file-a.txt"),
		path.Join(viewPath, "file-a-renamed.txt"),
	)
	assert.NoError(t, err)

	// Get comments again
	{
		getWorkspaceResolver, err := workspaceRootResolver.Workspace(authenticatedUserContext, resolvers.WorkspaceArgs{ID: graphql.ID(workspaceID)})
		assert.NoError(t, err)
		topComments, err := getWorkspaceResolver.Comments()
		assert.NoError(t, err)
		if assert.Len(t, topComments, 1) {
			assert.Equal(t, "Comment!", topComments[0].Message())
			codeContext := topComments[0].CodeContext()
			assert.Equal(t, int32(-1), codeContext.LineStart())
			assert.Equal(t, int32(-1), codeContext.LineEnd())
			assert.Equal(t, true, codeContext.LineIsNew())
			// TODO: it would be cool to support detecting that the file has been renamed to "file-a-renamed.txt"
			assert.Equal(t, "file-a.txt", codeContext.Path())
		}
	}

	{
		// Trigger GC
		err := gcService.WorkWithOptions(context.Background(), logger, codebaseRes.ID, 0, 0)
		assert.NoError(t, err)

		// make another change (after gc)
		// Remove the trailing newline in test.txt
		err = ioutil.WriteFile(path.Join(viewPath, "test-after-gc.txt"), []byte("THIS IS EPIC!"), 0o666)
		assert.NoError(t, err)

		// Get diff
		diffs, _, err = workspaceService.Diffs(context.Background(), workspaceID)
		assert.NoError(t, err)
		t.Logf("diffs=%+v", diffs)
		assert.Len(t, diffs, 2)
		assert.Len(t, diffs[0].Hunks, 1)

		// Apply and land
		_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
			WorkspaceID: graphql.ID(workspaceID),
		}})
		assert.NoError(t, err, errors.Unwrap(err))
	}
}

func TestLandEmpty(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		dig.In
		UserRepo              db_user.Repository
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		LandRootResolver      resolvers.LandRootResovler
		CodebaseService       *service_codebase.Service
		WorkspaceService      *service_workspace.Service
		RepoProvider          provider.RepoProvider

		// Dependencies of Gin Routes
		CodebaseUserRepo db_codebases.CodebaseUserRepository
		WorkspaceRepo    db_workspaces.Repository
		ViewRepo         db_view.Repository
		ExecutorProvider executor.Provider
		ViewService      *service_view.Service

		Logger           *zap.Logger
		AnalyticsService *service_analytics.Service

		SnapshotsQueue workers_snapshots.Queue
		CIQueue        *workers_ci.BuildQueue
	}

	var d deps
	if !assert.NoError(t, di.Init(testModule(t)).To(&d)) {
		t.FailNow()
	}

	go func() {
		assert.NoError(t, d.SnapshotsQueue.Start(context.TODO()))
	}()
	go func() {
		assert.NoError(t, d.CIQueue.Start(context.TODO()))
	}()

	createCodebaseRoute := routes_v3_codebase.Create(d.Logger, d.CodebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(d.Logger, d.WorkspaceService, d.CodebaseUserRepo)
	createViewRoute := routes_v3_view.Create(d.Logger, d.ViewRepo, d.CodebaseUserRepo, d.AnalyticsService, d.WorkspaceRepo, d.ExecutorProvider, d.ViewService)

	createUser := users.User{ID: users.ID(uuid.New().String()), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
	assert.NoError(t, d.UserRepo.Create(&createUser))

	authenticatedUserContext := gqldataloader.NewContext(auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: createUser.ID.String()}))

	// Create a codebase
	var codebaseRes codebases.Codebase
	request(t, createUser.ID, createCodebaseRoute, routes_v3_codebase.CreateRequest{Name: "testrepo"}, &codebaseRes)
	assert.Len(t, codebaseRes.ID, 36)
	assert.Equal(t, "testrepo", codebaseRes.Name)
	assert.True(t, codebaseRes.IsReady, "codebase is ready")

	// Create a workspace
	var workspaceResult workspaces.Workspace
	request(t, createUser.ID, createWorkspaceRoute, routes_v3_workspace.CreateRequest{
		CodebaseID: codebaseRes.ID,
	}, &workspaceResult)
	assert.Len(t, workspaceResult.ID, 36)
	workspaceID := workspaceResult.ID

	// Create a view
	var viewRes view.View
	request(t, createUser.ID, createViewRoute, routes_v3_view.CreateRequest{
		CodebaseID:    codebaseRes.ID,
		WorkspaceID:   workspaceID,
		MountPath:     "~/testing",
		MountHostname: "testing.ftw",
	}, &viewRes)
	assert.Len(t, viewRes.ID, 36)
	assert.True(t, viewRes.CreatedAt.After(time.Now().Add(time.Second*-5)))

	// Get diff
	{
		diffs, _, err := d.WorkspaceService.Diffs(context.Background(), workspaceID)
		assert.NoError(t, err)
		assert.Empty(t, diffs)
	}

	// Set workspace draft description
	_, err := d.WorkspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
		ID:               graphql.ID(workspaceID),
		DraftDescription: str("This is my first change"),
	}})
	assert.NoError(t, err)

	// Apply and land
	_, err = d.LandRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(workspaceID),
	}})
	assert.ErrorIs(t, err, gqlerror.ErrInternalServer)

	repo, err := d.RepoProvider.ViewRepo(codebaseRes.ID, viewRes.ID)
	assert.NoError(t, err)

	{
		hb, err := repo.HeadBranch()
		assert.NoError(t, err)
		assert.Equal(t, hb, workspaceID)
	}

	// Make some change, should be able to land it now
	viewPath := d.RepoProvider.ViewPath(codebaseRes.ID, viewRes.ID)
	err = ioutil.WriteFile(path.Join(viewPath, "hello-world.txt"), []byte("hello\n"), 0o666)
	assert.NoError(t, err)

	// Get diff
	{
		diffs, _, err := d.WorkspaceService.Diffs(context.Background(), workspaceID)
		assert.NoError(t, err)
		assert.NotEmpty(t, diffs)
	}

	// Apply and land
	_, err = d.LandRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: graphql.ID(workspaceID),
	}})
	assert.NoError(t, err)

	// Should be on a new branch now
	{
		hb, err := repo.HeadBranch()
		assert.NoError(t, err)
		assert.NotEqual(t, hb, workspaceID)
	}
}

func i(n int32) *int32 {
	return &n
}

func b(n bool) *bool {
	return &n
}

func gid(in graphql.ID) *graphql.ID {
	return &in
}

func TestLargeFiles(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		dig.In
		UserRepo              db_user.Repository
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		ViewRootResolver      resolvers.ViewRootResolver
		LandRootResolver      resolvers.LandRootResovler
		CodebaseService       *service_codebase.Service
		WorkspaceService      *service_workspace.Service
		GitSnapshotter        *service_snapshots.Service
		RepoProvider          provider.RepoProvider

		// Dependencies of Gin Routes
		CodebaseUserRepo db_codebases.CodebaseUserRepository
		WorkspaceRepo    db_workspaces.Repository
		ViewRepo         db_view.Repository
		SnapshotRepo     db_snapshots.Repository
		ExecutorProvider executor.Provider
		EventsSender     *eventsv2.Publisher
		WorkspaceWriter  db_workspaces.WorkspaceWriter
		ViewService      *service_view.Service

		Logger           *zap.Logger
		AnalyticsService *service_analytics.Service
	}

	var d deps
	if !assert.NoError(t, di.Init(testModule(t)).To(&d)) {
		t.FailNow()
	}

	userRepo := d.UserRepo
	repoProvider := d.RepoProvider
	workspaceService := d.WorkspaceService
	landRootResolver := d.LandRootResolver
	createCodebaseRoute := routes_v3_codebase.Create(d.Logger, d.CodebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(d.Logger, d.WorkspaceService, d.CodebaseUserRepo)
	createViewRoute := routes_v3_view.Create(d.Logger, d.ViewRepo, d.CodebaseUserRepo, d.AnalyticsService, d.WorkspaceRepo, d.ExecutorProvider, d.ViewService)

	testCases := []struct {
		name          string
		opts          []vcsvcs.DiffOption
		gitMaxSize    int
		largeFileName string
	}{
		{
			name:          "default",
			largeFileName: "large-img.jpg",
		},
		{
			name:          "low_max_size", // By default, files larger than 50MB have special treatment (are always treated as binary files), lower this to 500kb to make it easier to test
			opts:          []vcsvcs.DiffOption{vcsvcs.WithGitMaxSize(500_000)},
			gitMaxSize:    500_000,
			largeFileName: "large-img-low-max.jpg",
		},
		{
			name:          "low_max_spaces",
			opts:          []vcsvcs.DiffOption{vcsvcs.WithGitMaxSize(500_000)},
			gitMaxSize:    500_000,
			largeFileName: "with space.jpg",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			createUser := users.User{ID: users.ID(uuid.New().String()), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
			assert.NoError(t, userRepo.Create(&createUser))

			authenticatedUserContext := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: createUser.ID.String()})

			// Create a codebase
			var codebaseRes codebases.Codebase
			request(t, createUser.ID, createCodebaseRoute, routes_v3_codebase.CreateRequest{Name: "testrepo"}, &codebaseRes)
			assert.Len(t, codebaseRes.ID, 36)
			assert.Equal(t, "testrepo", codebaseRes.Name)
			assert.True(t, codebaseRes.IsReady, "codebase is ready")

			// Create a workspace
			var workspaceResult workspaces.Workspace
			request(t, createUser.ID, createWorkspaceRoute, routes_v3_workspace.CreateRequest{
				CodebaseID: codebaseRes.ID,
			}, &workspaceResult)
			assert.Len(t, workspaceResult.ID, 36)

			workspaceID := workspaceResult.ID

			t.Log("workspaceID", workspaceID)

			// Create a view
			var viewRes view.View
			request(t, createUser.ID, createViewRoute, routes_v3_view.CreateRequest{
				CodebaseID:    codebaseRes.ID,
				WorkspaceID:   workspaceID,
				MountPath:     "~/testing",
				MountHostname: "testing.ftw",
			}, &viewRes)
			assert.Len(t, viewRes.ID, 36)
			assert.True(t, viewRes.CreatedAt.After(time.Now().Add(time.Second*-5)))

			viewPath := repoProvider.ViewPath(codebaseRes.ID, viewRes.ID)
			t.Logf("viewPath=%s", viewPath)

			gitViewRepo, err := repoProvider.ViewRepo(codebaseRes.ID, viewRes.ID)
			assert.NoError(t, err)

			// Test large files
			{
				copy(t, "testdata/large-img.jpg", path.Join(viewPath, tc.largeFileName))

				// Get diff and apply
				diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceID, service_workspace.WithVCSDiffOptions(tc.opts...))
				assert.NoError(t, err)
				assert.Len(t, diffs, 1)
				assert.Len(t, diffs[0].Hunks, 1)
				t.Logf("diff: %+v", diffs[0])
				_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
					WorkspaceID: graphql.ID(workspaceID),
					DiffMaxSize: tc.gitMaxSize,
				}})
				assert.NoError(t, err)
			}

			getWorkspaceID := func() string {
				viewResolver, err := d.ViewRootResolver.View(authenticatedUserContext, resolvers.ViewArgs{ID: graphql.ID(viewRes.ID)})
				assert.NoError(t, err)

				newWorkspace, err := viewResolver.Workspace(authenticatedUserContext)
				assert.NoError(t, err)
				return string(newWorkspace.ID())
			}

			workspaceID = getWorkspaceID()

			t.Log("workspaceID", workspaceID)

			// Original file should be in the checkout, not the LFS pointer
			{
				stat, err := os.Stat(path.Join(viewPath, tc.largeFileName))
				assert.NoError(t, err)
				assert.True(t, stat.Size() > 1_000_000, "size=%d", stat.Size())
			}

			// LFS pointer should be in the latest commit
			headCommit, err := gitViewRepo.HeadCommit()
			assert.NoError(t, err)
			ptrContents, err := gitViewRepo.FileContentsAtCommit(headCommit.Id().String(), tc.largeFileName)
			assert.NoError(t, err)
			assert.True(t, len(ptrContents) < 500, "len=%d", len(ptrContents))

			// Create file with space in the name
			{
				nameWithSpaces := path.Join(viewPath, "dir", "dir with space", "Aspen 0.1.6.dmg")
				copy(t, "testdata/large-img.jpg", nameWithSpaces)

				// Get diff and apply
				diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceID, service_workspace.WithVCSDiffOptions(tc.opts...))
				assert.NoError(t, err)
				assert.Len(t, diffs, 1)
				assert.Len(t, diffs[0].Hunks, 1)
				t.Logf("diff: %+v", diffs[0])

				_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
					WorkspaceID: graphql.ID(workspaceID),
					DiffMaxSize: tc.gitMaxSize,
				}})
				assert.NoError(t, err)

				// Verify that file was shared
				fp, err := os.Open(nameWithSpaces)
				assert.NoError(t, err)
				finfo, err := fp.Stat()
				assert.NoError(t, err)
				assert.True(t, finfo.Size() > 1_000_000, "size=%d", finfo.Size())
			}

			workspaceID = getWorkspaceID()

			t.Log("workspaceID", workspaceID)

			// Update the large file
			{
				copy(t, "testdata/large-img-2.jpg", path.Join(viewPath, tc.largeFileName))
				// Get diff and apply
				diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceID, service_workspace.WithVCSDiffOptions(tc.opts...))
				assert.NoError(t, err)
				assert.Len(t, diffs, 1)
				assert.Len(t, diffs[0].Hunks, 1)
				t.Logf("diffs: %+v", diffs)

				_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
					WorkspaceID: graphql.ID(workspaceID),
					DiffMaxSize: tc.gitMaxSize,
				}})
				assert.NoError(t, err, errors.Unwrap(err))
			}

			// return
			workspaceID = getWorkspaceID()

			// LFS pointer should be updated
			{
				headCommit, err := gitViewRepo.HeadCommit()
				assert.NoError(t, err)
				ptrContents2, err := gitViewRepo.FileContentsAtCommit(headCommit.Id().String(), tc.largeFileName)
				assert.NoError(t, err)
				assert.True(t, len(ptrContents2) < 500, "len=%d", len(ptrContents2))
				assert.NotEqual(t, string(ptrContents2), string(ptrContents))
			}

			// Delete the large file
			{
				err = os.Remove(path.Join(viewPath, tc.largeFileName))
				assert.NoError(t, err)

				// Get diff and apply
				diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceID, service_workspace.WithVCSDiffOptions(tc.opts...))
				assert.NoError(t, err)
				assert.Len(t, diffs, 1)
				assert.Len(t, diffs[0].Hunks, 1)
				_, err = landRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
					WorkspaceID: graphql.ID(workspaceID),
					DiffMaxSize: tc.gitMaxSize,
				}})
				assert.NoError(t, err)

				// LFS pointer should not exist
				headCommit, err = gitViewRepo.HeadCommit()
				assert.NoError(t, err)
				_, err = gitViewRepo.FileContentsAtCommit(headCommit.Id().String(), tc.largeFileName)
				assert.Error(t, err)
			}
		})
	}
}

func copy(t *testing.T, src string, dst string) {
	err := os.MkdirAll(path.Dir(dst), 0777)
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(src)
	assert.NoError(t, err)
	err = ioutil.WriteFile(dst, data, 0644)
	assert.NoError(t, err)
}

func request(t *testing.T, userID users.ID, route func(*gin.Context), request, response any) {
	requestWithParams(t, userID, route, request, response, nil)
}

func requestWithParams(t *testing.T, userID users.ID, route func(*gin.Context), request, response any, params []gin.Param) {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Params = params

	data, err := json.Marshal(request)
	assert.NoError(t, err)

	c.Request, err = http.NewRequest("POST", "/", bytes.NewReader(data))
	c.Request = c.Request.WithContext(auth.NewContext(context.Background(), &auth.Subject{ID: userID.String(), Type: auth.SubjectUser}))
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

func str(s string) *string {
	return &s
}

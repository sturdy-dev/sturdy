package pkg_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"go.uber.org/dig"

	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebase"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	routes_v3_codebase "getsturdy.com/api/pkg/codebase/routes"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/di"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	routes_v3_view "getsturdy.com/api/pkg/view/routes"
	service_view "getsturdy.com/api/pkg/view/service"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	routes_v3_workspace "getsturdy.com/api/pkg/workspaces/routes"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRevertChangeFromSnapshot(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		dig.In
		UserRepo              db_user.Repository
		CodebaseRootResolver  resolvers.CodebaseRootResolver
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		ViewRootResolver      resolvers.ViewRootResolver

		CodebaseService  *service_codebase.Service
		WorkspaceService service_workspace.Service
		GitSnapshotter   snapshotter.Snapshotter
		RepoProvider     provider.RepoProvider

		CodebaseUserRepo db_codebase.CodebaseUserRepository
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
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	}

	repoProvider := d.RepoProvider
	userRepo := d.UserRepo
	workspaceService := d.WorkspaceService
	workspaceRepo := d.WorkspaceRepo

	createCodebaseRoute := routes_v3_codebase.Create(d.Logger, d.CodebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(d.Logger, d.WorkspaceService, d.CodebaseUserRepo)
	createViewRoute := routes_v3_view.Create(d.Logger, d.ViewRepo, d.CodebaseUserRepo, d.AnalyticsService, d.WorkspaceRepo, d.ExecutorProvider, d.ViewService)

	workspaceRootResolver := d.WorkspaceRootResolver
	codebaseRootResolver := d.CodebaseRootResolver

	createUser := users.User{ID: users.ID(uuid.New().String()), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
	assert.NoError(t, userRepo.Create(&createUser))

	authenticatedUserContext := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: createUser.ID.String()})

	// Create a codebase
	var codebaseRes codebase.Codebase
	request(t, createUser.ID, createCodebaseRoute, routes_v3_codebase.CreateRequest{Name: "testrepo"}, &codebaseRes)
	assert.Len(t, codebaseRes.ID, 36)
	assert.Equal(t, "testrepo", codebaseRes.Name)
	assert.True(t, codebaseRes.IsReady, "codebase is ready")

	// Create a workspace
	var workspaceRes workspaces.Workspace
	request(t, createUser.ID, createWorkspaceRoute, routes_v3_workspace.CreateRequest{
		CodebaseID: codebaseRes.ID,
	}, &workspaceRes)
	assert.Len(t, workspaceRes.ID, 36)

	// Create a view
	var viewRes view.View
	request(t, createUser.ID, createViewRoute, routes_v3_view.CreateRequest{
		CodebaseID:    codebaseRes.ID,
		WorkspaceID:   workspaceRes.ID,
		MountPath:     "~/testing",
		MountHostname: "testing.ftw",
	}, &viewRes)
	assert.Len(t, viewRes.ID, 36)
	assert.True(t, viewRes.CreatedAt.After(time.Now().Add(time.Second*-5)))

	viewPath := repoProvider.ViewPath(codebaseRes.ID, viewRes.ID)

	getWorkspaceID := func() string {
		viewResolver, err := d.ViewRootResolver.View(authenticatedUserContext, resolvers.ViewArgs{ID: graphql.ID(viewRes.ID)})
		assert.NoError(t, err)

		wsResolver, err := viewResolver.Workspace(authenticatedUserContext)
		assert.NoError(t, err)

		return string(wsResolver.ID())
	}

	t.Logf("viewPath=%s", viewPath)

	changes := []struct {
		file     string
		contents string
	}{
		{"hello.txt", "hello-first-change\n"},
		{"hello.txt", "hello-second-change\n"},
		{"wat.txt", "wat\n"},
	}

	for _, ch := range changes {
		// Make changes in the view
		err := ioutil.WriteFile(path.Join(viewPath, ch.file), []byte(ch.contents), 0o666)
		assert.NoError(t, err)

		workspaceID := getWorkspaceID()

		// Get diff
		diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceID)
		assert.NoError(t, err)

		// Set workspace draft description
		_, err = workspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
			ID:               graphql.ID(workspaceID),
			DraftDescription: &ch.file,
		}})
		assert.NoError(t, err)
		// Apply and land
		_, err = workspaceRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
			WorkspaceID: graphql.ID(workspaceID),
			PatchIDs:    []string{diffs[0].Hunks[0].ID},
		}})
		assert.NoError(t, err)
	}

	// Get changelog
	cid := graphql.ID(codebaseRes.ID)
	cbResolver, err := codebaseRootResolver.Codebase(authenticatedUserContext, resolvers.CodebaseArgs{ID: &cid})
	assert.NoError(t, err)
	changeResolvers, err := cbResolver.Changes(authenticatedUserContext, nil)
	assert.NoError(t, err)
	assert.Len(t, changeResolvers, 3)

	// Pick the 2nd change, and revert it
	revertID := changeResolvers[1].ID()
	revertedWsResolver, err := workspaceRootResolver.CreateWorkspace(authenticatedUserContext, resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{
		CodebaseID:              cid,
		OnTopOfChangeWithRevert: &revertID,
	}})
	assert.NoError(t, err)
	assert.Equal(t, "Revert hello.txt", revertedWsResolver.Name())

	// Check the file on trunk (pre land)
	fileOrDirResolver, err := cbResolver.File(authenticatedUserContext, resolvers.CodebaseFileArgs{Path: "hello.txt"})
	fileResolver, ok := fileOrDirResolver.ToFile()
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, "hello-second-change\n", fileResolver.Contents())

	revertedWs, err := workspaceRepo.Get(string(revertedWsResolver.ID()))
	assert.NoError(t, err)

	// Land the change
	// Get diff
	diffs, _, err := workspaceService.Diffs(authenticatedUserContext, revertedWs.ID)
	assert.NoError(t, err)

	// Set workspace draft description
	_, err = workspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
		ID:               revertedWsResolver.ID(),
		DraftDescription: str("Revert"),
	}})
	assert.NoError(t, err)
	// Apply and land
	_, err = workspaceRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: revertedWsResolver.ID(),
		PatchIDs:    []string{diffs[0].Hunks[0].ID},
	}})
	assert.NoError(t, err)
	//nolint:errorlint
	if gerr, ok := err.(*gqlerrors.SturdyGraphqlError); ok {
		t.Logf("err=%+v", gerr.OriginalError())
	}

	// Check file on trunk
	fileOrDirResolver, err = cbResolver.File(authenticatedUserContext, resolvers.CodebaseFileArgs{Path: "hello.txt"})
	fileResolver, ok = fileOrDirResolver.ToFile()
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, "hello-first-change\n", fileResolver.Contents())

	// Pick the 2nd change, and create new workspace from it
	changeID := changeResolvers[1].ID()
	newWs, err := workspaceRootResolver.CreateWorkspace(authenticatedUserContext, resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{
		CodebaseID:    cid,
		OnTopOfChange: &changeID,
	}})
	assert.NoError(t, err)
	assert.Equal(t, "On hello.txt", newWs.Name())

	isUpToDateWithTrunk, err := newWs.UpToDateWithTrunk(authenticatedUserContext)
	assert.NoError(t, err)
	assert.False(t, isUpToDateWithTrunk)
}

func TestRevertChangeFromView(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		dig.In
		UserRepo              db_user.Repository
		CodebaseRootResolver  resolvers.CodebaseRootResolver
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		ViewRootResolver      resolvers.ViewRootResolver

		CodebaseService  *service_codebase.Service
		WorkspaceService service_workspace.Service
		GitSnapshotter   snapshotter.Snapshotter
		RepoProvider     provider.RepoProvider

		CodebaseUserRepo db_codebase.CodebaseUserRepository
		WorkspaceRepo    db_workspaces.Repository
		ViewRepo         db_view.Repository
		SnapshotRepo     db_snapshots.Repository
		ExecutorProvider executor.Provider
		EventsSender     *eventsv2.Publisher
		ViewService      *service_view.Service

		Logger           *zap.Logger
		AnalyticsSerivce *service_analytics.Service
	}

	var d deps
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	}

	repoProvider := d.RepoProvider
	userRepo := d.UserRepo
	workspaceService := d.WorkspaceService
	viewRootResolver := d.ViewRootResolver

	createCodebaseRoute := routes_v3_codebase.Create(d.Logger, d.CodebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(d.Logger, d.WorkspaceService, d.CodebaseUserRepo)
	createViewRoute := routes_v3_view.Create(d.Logger, d.ViewRepo, d.CodebaseUserRepo, d.AnalyticsSerivce, d.WorkspaceRepo, d.ExecutorProvider, d.ViewService)

	workspaceRootResolver := d.WorkspaceRootResolver
	codebaseRootResolver := d.CodebaseRootResolver

	createUser := users.User{ID: users.ID(uuid.New().String()), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
	assert.NoError(t, userRepo.Create(&createUser))

	authenticatedUserContext := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: createUser.ID.String()})

	// Create a codebase
	var codebaseRes codebase.Codebase
	request(t, createUser.ID, createCodebaseRoute, routes_v3_codebase.CreateRequest{Name: "testrepo"}, &codebaseRes)
	assert.Len(t, codebaseRes.ID, 36)
	assert.Equal(t, "testrepo", codebaseRes.Name)
	assert.True(t, codebaseRes.IsReady, "codebase is ready")

	// Create a workspace
	var workspaceRes workspaces.Workspace
	request(t, createUser.ID, createWorkspaceRoute, routes_v3_workspace.CreateRequest{
		CodebaseID: codebaseRes.ID,
	}, &workspaceRes)
	assert.Len(t, workspaceRes.ID, 36)

	// Create a view
	var viewRes view.View
	request(t, createUser.ID, createViewRoute, routes_v3_view.CreateRequest{
		CodebaseID:    codebaseRes.ID,
		WorkspaceID:   workspaceRes.ID,
		MountPath:     "~/testing",
		MountHostname: "testing.ftw",
	}, &viewRes)
	assert.Len(t, viewRes.ID, 36)
	assert.True(t, viewRes.CreatedAt.After(time.Now().Add(time.Second*-5)))

	viewPath := repoProvider.ViewPath(codebaseRes.ID, viewRes.ID)

	t.Logf("viewPath=%s", viewPath)

	getWorkspaceID := func() string {
		viewResolver, err := d.ViewRootResolver.View(authenticatedUserContext, resolvers.ViewArgs{ID: graphql.ID(viewRes.ID)})
		assert.NoError(t, err)

		wsResolver, err := viewResolver.Workspace(authenticatedUserContext)
		assert.NoError(t, err)

		return string(wsResolver.ID())
	}

	changes := []struct {
		file     string
		contents string
	}{
		{"hello.txt", "hello-first-change\n"},
		{"hello.txt", "hello-second-change\n"},
		{"wat.txt", "wat\n"},
	}

	for _, ch := range changes {
		// Make changes in the view
		err := ioutil.WriteFile(path.Join(viewPath, ch.file), []byte(ch.contents), 0o666)
		assert.NoError(t, err)

		workspaceID := getWorkspaceID()

		// Get diff
		diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceID)
		assert.NoError(t, err)

		// Set workspace draft description
		_, err = workspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
			ID:               graphql.ID(workspaceID),
			DraftDescription: &ch.file,
		}})
		assert.NoError(t, err)
		// Apply and land
		_, err = workspaceRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
			WorkspaceID: graphql.ID(workspaceID),
			PatchIDs:    []string{diffs[0].Hunks[0].ID},
		}})
		assert.NoError(t, err)
	}

	// Get changelog
	cid := graphql.ID(codebaseRes.ID)
	cbResolver, err := codebaseRootResolver.Codebase(authenticatedUserContext, resolvers.CodebaseArgs{ID: &cid})
	assert.NoError(t, err)
	changeResolvers, err := cbResolver.Changes(authenticatedUserContext, nil)
	assert.NoError(t, err)
	assert.Len(t, changeResolvers, 3)

	// Pick the 2nd change, and revert it
	revertID := changeResolvers[1].ID()
	revertedWs, err := workspaceRootResolver.CreateWorkspace(authenticatedUserContext, resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{
		CodebaseID:              cid,
		OnTopOfChangeWithRevert: &revertID,
	}})
	assert.NoError(t, err)
	assert.Equal(t, "Revert hello.txt", revertedWs.Name())

	// Open this workspace
	_, err = viewRootResolver.OpenWorkspaceOnView(authenticatedUserContext, resolvers.OpenViewArgs{Input: resolvers.OpenWorkspaceOnViewInput{
		WorkspaceID: revertedWs.ID(),
		ViewID:      graphql.ID(viewRes.ID),
	}})
	if !assert.NoError(t, err) {
		//nolint:errorlint
		t.Logf("err=%+v", err.(*gqlerrors.SturdyGraphqlError).OriginalError())
	}

	// Check the file on trunk (pre land)
	fileOrDirResolver, err := cbResolver.File(authenticatedUserContext, resolvers.CodebaseFileArgs{Path: "hello.txt"})
	fileResolver, ok := fileOrDirResolver.ToFile()
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, "hello-second-change\n", fileResolver.Contents())

	// Land the change
	// Get diff
	diffs, _, err := workspaceService.Diffs(authenticatedUserContext, string(revertedWs.ID()))
	assert.NoError(t, err)

	// Set workspace draft description
	_, err = workspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
		ID:               revertedWs.ID(),
		DraftDescription: str("Revert"),
	}})
	assert.NoError(t, err)
	// Apply and land
	_, err = workspaceRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
		WorkspaceID: revertedWs.ID(),
		PatchIDs:    []string{diffs[0].Hunks[0].ID},
	}})
	assert.NoError(t, err)

	// Check file on trunk
	fileOrDirResolver, err = cbResolver.File(authenticatedUserContext, resolvers.CodebaseFileArgs{Path: "hello.txt"})
	fileResolver, ok = fileOrDirResolver.ToFile()
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, "hello-first-change\n", fileResolver.Contents())

	// Pick the 2nd change, and create new workspace from it
	changeID := changeResolvers[1].ID()
	newWs, err := workspaceRootResolver.CreateWorkspace(authenticatedUserContext, resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{
		CodebaseID:    cid,
		OnTopOfChange: &changeID,
	}})
	assert.NoError(t, err)
	assert.Equal(t, "On hello.txt", newWs.Name())

	isUpToDateWithTrunk, err := newWs.UpToDateWithTrunk(authenticatedUserContext)
	assert.NoError(t, err)
	assert.False(t, isUpToDateWithTrunk)
}

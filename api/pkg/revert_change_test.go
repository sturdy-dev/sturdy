package pkg_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"getsturdy.com/api/pkg/analytics/disabled"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_change "getsturdy.com/api/pkg/change/db"
	graphql_change "getsturdy.com/api/pkg/change/graphql"
	service_change "getsturdy.com/api/pkg/change/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebase"
	db_acl "getsturdy.com/api/pkg/codebase/acl/db"
	provider_acl "getsturdy.com/api/pkg/codebase/acl/provider"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	graphql_codebase "getsturdy.com/api/pkg/codebase/graphql"
	routes_v3_codebase "getsturdy.com/api/pkg/codebase/routes"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	db_comments "getsturdy.com/api/pkg/comments/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/events"
	graphql_file "getsturdy.com/api/pkg/file/graphql"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/internal/sturdytest"
	"getsturdy.com/api/pkg/notification/sender"
	"getsturdy.com/api/pkg/queue"
	db_review "getsturdy.com/api/pkg/review/db"
	graphql_review "getsturdy.com/api/pkg/review/graphql"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	db_statuses "getsturdy.com/api/pkg/statuses/db"
	graphql_statuses "getsturdy.com/api/pkg/statuses/graphql"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	db_suggestion "getsturdy.com/api/pkg/suggestions/db"
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	graphql_view "getsturdy.com/api/pkg/view/graphql"
	routes_v3_view "getsturdy.com/api/pkg/view/routes"
	"getsturdy.com/api/pkg/workspace"
	db_activity "getsturdy.com/api/pkg/workspace/activity/db"
	activity_sender "getsturdy.com/api/pkg/workspace/activity/sender"
	service_activity "getsturdy.com/api/pkg/workspace/activity/service"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	graphql_workspace "getsturdy.com/api/pkg/workspace/graphql"
	ws_meta "getsturdy.com/api/pkg/workspace/meta"
	routes_v3_workspace "getsturdy.com/api/pkg/workspace/routes"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	db_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/db"
	graphql_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/graphql"
	service_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/service"
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

	reposBasePath := os.TempDir()
	repoProvider := provider.New(reposBasePath, "localhost:8888")

	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	if err != nil {
		panic(err)
	}

	logger := zap.NewNop()
	postHogClient := disabled.NewClient()

	userRepo := db_user.NewRepo(d)
	codebaseRepo := db_codebase.NewRepo(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	workspaceRepo := db_workspace.NewRepo(d)
	viewRepo := db_view.NewRepo(d)
	changeRepo := db_change.NewRepo(d)
	changeCommitRepo := db_change.NewCommitRepository(d)
	snapshotRepo := db_snapshots.NewRepo(d)
	commentRepo := db_comments.NewRepo(d)
	workspaceActivityRepo := db_activity.NewActivityRepo(d)
	workspaceActivityReadsRepo := db_activity.NewActivityReadsRepo(d)
	viewEvents := events.NewInMemory()
	executorProvider := executor.NewProvider(logger, repoProvider)
	reviewRepo := db_review.NewReviewRepository(d)
	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, viewEvents)
	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotRepo, workspaceRepo, workspaceRepo, viewRepo, eventsSender, executorProvider, logger)
	snapshotPublisher := worker_snapshots.NewSync(gitSnapshotter)
	activityService := service_activity.New(workspaceActivityReadsRepo, eventsSender)
	activitySender := activity_sender.NewActivitySender(codebaseUserRepo, workspaceActivityRepo, activityService, eventsSender)
	suggestionRepo := db_suggestion.New(d)
	notificationSender := sender.NewNoopNotificationSender()
	commentsService := service_comments.New(commentRepo)
	aclRepo := db_acl.NewACLRepository(d)

	queue := queue.NewNoop()
	buildQueue := workers_ci.New(zap.NewNop(), queue, nil)
	userService := service_user.New(zap.NewNop(), userRepo, postHogClient)

	aclProvider := provider_acl.New(aclRepo, codebaseUserRepo, userRepo)

	workspaceWriter := ws_meta.NewWriterWithEvents(logger, workspaceRepo, eventsSender)
	changeService := service_change.New(aclProvider, userRepo, changeRepo, changeCommitRepo)
	workspaceService := service_workspace.New(
		logger,
		postHogClient,

		workspaceWriter,
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

	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, workspaceService, nil, logger, executorProvider, postHogClient, eventsSender)

	suggestionsService := service_suggestion.New(
		logger,
		suggestionRepo,
		workspaceService,
		executorProvider,
		gitSnapshotter,
		postHogClient,
		notificationSender,
		eventsSender,
	)

	authService := service_auth.New(codebaseService, userService, workspaceService, aclProvider, nil /*organizationService*/)

	createCodebaseRoute := routes_v3_codebase.Create(logger, codebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(logger, workspaceService, codebaseUserRepo)
	createViewRoute := routes_v3_view.Create(logger, viewRepo, codebaseUserRepo, postHogClient, workspaceRepo, gitSnapshotter, snapshotRepo, workspaceWriter, executorProvider, eventsSender)

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
		nil,
		workspaceWatchersService,
	)

	statusesRepo := db_statuses.New(d)
	statusesServcie := service_statuses.New(logger, statusesRepo, eventsSender)
	statusesRootResolver := new(resolvers.StatusesRootResolver)

	changeRootResolver := graphql_change.NewResolver(
		changeService,
		changeRepo,
		changeCommitRepo,
		commentRepo,
		authService,
		nil, // commentResolver
		nil, // authorResolver
		statusesRootResolver,
		nil, // downloadsResovler
		executorProvider,
		logger,
	)

	fileRootResolver := graphql_file.NewFileRootResolver(executorProvider, authService)

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
		changeRootResolver,
		nil, // workspaceActivityResolver
		reviewRootResolver,
		nil, // presenseRootResolver
		nil, // suggestitonsRootResolver
		*statusesRootResolver,
		*workspaceWatchersRootResolver,
		suggestionsService,
		workspaceService,
		authService,
		logger,
		viewEvents,
		workspaceWriter,
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

	codebaseRootResolver := graphql_codebase.NewCodebaseRootResolver(
		codebaseRepo,
		codebaseUserRepo,
		viewRepo,
		workspaceRepo,
		userRepo,
		changeRepo,
		changeCommitRepo,

		nil,
		nil,
		nil,
		nil,
		changeRootResolver,
		fileRootResolver,
		nil, // instantIntegrationRootResolver
		nil,
		nil,

		logger,
		nil,
		nil,
		postHogClient,
		executorProvider,

		authService,
		codebaseService,
		nil,
	)

	*statusesRootResolver = graphql_statuses.New(
		logger,
		statusesServcie,
		changeService,
		workspaceService,
		authService,
		changeRootResolver,
		nil, // github pr resolver
		viewEvents,
	)

	createUser := users.User{ID: uuid.New().String(), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
	assert.NoError(t, userRepo.Create(&createUser))

	authenticatedUserContext := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: createUser.ID})

	// Create a codebase
	var codebaseRes codebase.Codebase
	request(t, createUser.ID, createCodebaseRoute, routes_v3_codebase.CreateRequest{Name: "testrepo"}, &codebaseRes)
	assert.Len(t, codebaseRes.ID, 36)
	assert.Equal(t, "testrepo", codebaseRes.Name)
	assert.True(t, codebaseRes.IsReady, "codebase is ready")

	// Create a workspace
	var workspaceRes workspace.Workspace
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
		err = ioutil.WriteFile(path.Join(viewPath, ch.file), []byte(ch.contents), 0o666)
		assert.NoError(t, err)

		// Get diff
		diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceRes.ID)
		assert.NoError(t, err)

		// Set workspace draft description
		_, err = workspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
			ID:               graphql.ID(workspaceRes.ID),
			DraftDescription: &ch.file,
		}})
		assert.NoError(t, err)
		// Apply and land
		_, err = workspaceRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
			WorkspaceID: graphql.ID(workspaceRes.ID),
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

	isUpToDateWithTrunk, err := newWs.UpToDateWithTrunk()
	assert.NoError(t, err)
	assert.False(t, isUpToDateWithTrunk)
}

func TestRevertChangeFromView(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	reposBasePath := os.TempDir()
	repoProvider := provider.New(reposBasePath, "localhost:8888")

	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	if err != nil {
		panic(err)
	}

	logger := zap.NewNop()
	postHogClient := disabled.NewClient()

	userRepo := db_user.NewRepo(d)
	codebaseRepo := db_codebase.NewRepo(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	workspaceRepo := db_workspace.NewRepo(d)
	viewRepo := db_view.NewRepo(d)
	changeRepo := db_change.NewRepo(d)
	changeCommitRepo := db_change.NewCommitRepository(d)
	snapshotRepo := db_snapshots.NewRepo(d)
	commentRepo := db_comments.NewRepo(d)
	workspaceActivityRepo := db_activity.NewActivityRepo(d)
	workspaceActivityReadsRepo := db_activity.NewActivityReadsRepo(d)
	viewEvents := events.NewInMemory()
	executorProvider := executor.NewProvider(logger, repoProvider)
	reviewRepo := db_review.NewReviewRepository(d)
	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, viewEvents)
	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotRepo, workspaceRepo, workspaceRepo, viewRepo, eventsSender, executorProvider, logger)
	snapshotPublisher := worker_snapshots.NewSync(gitSnapshotter)
	activityService := service_activity.New(workspaceActivityReadsRepo, eventsSender)
	activitySender := activity_sender.NewActivitySender(codebaseUserRepo, workspaceActivityRepo, activityService, eventsSender)
	suggestionRepo := db_suggestion.New(d)
	notificationSender := sender.NewNoopNotificationSender()
	commentsService := service_comments.New(commentRepo)
	aclRepo := db_acl.NewACLRepository(d)

	queue := queue.NewNoop()
	buildQueue := workers_ci.New(zap.NewNop(), queue, nil)
	userService := service_user.New(zap.NewNop(), userRepo, postHogClient)

	workspaceWriter := ws_meta.NewWriterWithEvents(logger, workspaceRepo, eventsSender)
	changeService := service_change.New(nil, userRepo, changeRepo, changeCommitRepo)

	workspaceService := service_workspace.New(
		logger,
		postHogClient,

		workspaceWriter,
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

	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, workspaceService, nil, logger, executorProvider, postHogClient, eventsSender)

	suggestionsService := service_suggestion.New(
		logger,
		suggestionRepo,
		workspaceService,
		executorProvider,
		gitSnapshotter,
		postHogClient,
		notificationSender,
		eventsSender,
	)

	aclProvider := provider_acl.New(aclRepo, codebaseUserRepo, userRepo)

	authService := service_auth.New(codebaseService, userService, workspaceService, aclProvider, nil /*organizationService*/)

	createCodebaseRoute := routes_v3_codebase.Create(logger, codebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(logger, workspaceService, codebaseUserRepo)
	createViewRoute := routes_v3_view.Create(logger, viewRepo, codebaseUserRepo, postHogClient, workspaceRepo, gitSnapshotter, snapshotRepo, workspaceWriter, executorProvider, eventsSender)

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
		nil,
		workspaceWatchersService,
	)

	statusesRepo := db_statuses.New(d)
	statusesServcie := service_statuses.New(logger, statusesRepo, eventsSender)
	statusesRootResolver := new(resolvers.StatusesRootResolver)

	changeRootResolver := graphql_change.NewResolver(
		changeService,
		changeRepo,
		changeCommitRepo,
		commentRepo,
		authService,
		nil, // commentResolver
		nil, // authorResolver
		statusesRootResolver,
		nil, // downloadsResolver
		executorProvider,
		logger,
	)

	fileRootResolver := graphql_file.NewFileRootResolver(executorProvider, authService)

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
		changeRootResolver,
		nil, // workspaceActivityResolver
		reviewRootResolver,
		nil, // presenseRootResolver
		nil, // suggestitonsRootResolver
		*statusesRootResolver,
		*workspaceWatchersRootResolver,
		suggestionsService,
		workspaceService,
		authService,
		logger,
		viewEvents,
		workspaceWriter,
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
		snapshotRepo,
		nil,
		nil,
		workspaceWriter,
		viewEvents,
		eventsSender,
		executorProvider,
		logger,
		nil,
		workspaceWatchersService,
		postHogClient,
		nil,
		authService,
	)

	codebaseRootResolver := graphql_codebase.NewCodebaseRootResolver(
		codebaseRepo,
		codebaseUserRepo,
		viewRepo,
		workspaceRepo,
		userRepo,
		changeRepo,
		changeCommitRepo,

		nil,
		nil,
		nil,
		nil,
		changeRootResolver,
		fileRootResolver,
		nil, // instantIntegrationRootResolver
		nil,
		nil,

		logger,
		nil,
		nil,
		postHogClient,
		executorProvider,

		authService,
		codebaseService,
		nil,
	)

	*statusesRootResolver = graphql_statuses.New(
		logger,
		statusesServcie,
		changeService,
		workspaceService,
		authService,
		changeRootResolver,
		nil, // github pr resolver
		viewEvents,
	)

	createUser := users.User{ID: uuid.New().String(), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
	assert.NoError(t, userRepo.Create(&createUser))

	authenticatedUserContext := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: createUser.ID})

	// Create a codebase
	var codebaseRes codebase.Codebase
	request(t, createUser.ID, createCodebaseRoute, routes_v3_codebase.CreateRequest{Name: "testrepo"}, &codebaseRes)
	assert.Len(t, codebaseRes.ID, 36)
	assert.Equal(t, "testrepo", codebaseRes.Name)
	assert.True(t, codebaseRes.IsReady, "codebase is ready")

	// Create a workspace
	var workspaceRes workspace.Workspace
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
		err = ioutil.WriteFile(path.Join(viewPath, ch.file), []byte(ch.contents), 0o666)
		assert.NoError(t, err)

		// Get diff
		diffs, _, err := workspaceService.Diffs(authenticatedUserContext, workspaceRes.ID)
		assert.NoError(t, err)

		// Set workspace draft description
		_, err = workspaceRootResolver.UpdateWorkspace(authenticatedUserContext, resolvers.UpdateWorkspaceArgs{Input: resolvers.UpdateWorkspaceInput{
			ID:               graphql.ID(workspaceRes.ID),
			DraftDescription: &ch.file,
		}})
		assert.NoError(t, err)
		// Apply and land
		_, err = workspaceRootResolver.LandWorkspaceChange(authenticatedUserContext, resolvers.LandWorkspaceArgs{Input: resolvers.LandWorkspaceInput{
			WorkspaceID: graphql.ID(workspaceRes.ID),
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

	isUpToDateWithTrunk, err := newWs.UpToDateWithTrunk()
	assert.NoError(t, err)
	assert.False(t, isUpToDateWithTrunk)
}

package pkg_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"mash/db"
	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	db_change "mash/pkg/change/db"
	graphql_change "mash/pkg/change/graphql"
	service_change "mash/pkg/change/service"
	workers_ci "mash/pkg/ci/workers"
	"mash/pkg/codebase"
	db_codebase "mash/pkg/codebase/db"
	graphql_codebase "mash/pkg/codebase/graphql"
	routes_v3_codebase "mash/pkg/codebase/routes"
	service_codebase "mash/pkg/codebase/service"
	db_comments "mash/pkg/comments/db"
	service_comments "mash/pkg/comments/service"
	graphql_file "mash/pkg/file/graphql"
	"mash/pkg/github/config"
	db_github "mash/pkg/github/db"
	workers_github "mash/pkg/github/enterprise/workers"
	service_github "mash/pkg/github/service"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/internal/sturdytest"
	"mash/pkg/notification/sender"
	"mash/pkg/posthog"
	"mash/pkg/queue"
	db_review "mash/pkg/review/db"
	graphql_review "mash/pkg/review/graphql"
	db_snapshots "mash/pkg/snapshots/db"
	"mash/pkg/snapshots/snapshotter"
	worker_snapshots "mash/pkg/snapshots/worker"
	db_statuses "mash/pkg/statuses/db"
	graphql_statuses "mash/pkg/statuses/graphql"
	service_statuses "mash/pkg/statuses/service"
	db_suggestion "mash/pkg/suggestions/db"
	service_suggestion "mash/pkg/suggestions/service"
	"mash/pkg/user"
	db_user "mash/pkg/user/db"
	service_user "mash/pkg/user/service"
	"mash/pkg/view"
	db_view "mash/pkg/view/db"
	"mash/pkg/view/events"
	graphql_view "mash/pkg/view/graphql"
	routes_v3_view "mash/pkg/view/routes"
	"mash/pkg/view/view_workspace_snapshot"
	"mash/pkg/workspace"
	db_activity "mash/pkg/workspace/activity/db"
	activity_sender "mash/pkg/workspace/activity/sender"
	service_activity "mash/pkg/workspace/activity/service"
	db_workspace "mash/pkg/workspace/db"
	graphql_workspace "mash/pkg/workspace/graphql"
	ws_meta "mash/pkg/workspace/meta"
	routes_v3_workspace "mash/pkg/workspace/routes"
	service_workspace "mash/pkg/workspace/service"
	db_workspace_watchers "mash/pkg/workspace/watchers/db"
	graphql_workspace_watchers "mash/pkg/workspace/watchers/graphql"
	service_workspace_watchers "mash/pkg/workspace/watchers/service"
	"mash/vcs/executor"
	"mash/vcs/provider"

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
		true,
		"file://../db/migrations",
	)
	if err != nil {
		panic(err)
	}

	logger := zap.NewNop()
	postHogClient := posthog.NewFakeClient()

	userRepo := db_user.NewRepo(d)
	codebaseRepo := db_codebase.NewRepo(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	workspaceRepo := db_workspace.NewRepo(d)
	viewRepo := db_view.NewRepo(d)
	changeRepo := db_change.NewRepo(d)
	changeCommitRepo := db_change.NewCommitRepository(d)
	snapshotRepo := db_snapshots.NewRepo(d)
	gitHubRepositoryRepo := db_github.NewGitHubRepositoryRepo(d)
	gitHubInstallationRepo := db_github.NewGitHubInstallationRepo(d)
	gitHubPRRepo := db_github.NewGitHubPRRepo(d)
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

	queue := queue.NewNoop()
	buildQueue := workers_ci.New(zap.NewNop(), queue, nil)
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo)
	userService := service_user.New(zap.NewNop(), userRepo, nil /*jwtService*/, nil /*onetime*/, nil /*emailsender*/, postHogClient)

	workspaceWriter := ws_meta.NewWriterWithEvents(logger, workspaceRepo, eventsSender)
	changeService := service_change.New(executorProvider, nil, nil, userRepo, changeRepo, changeCommitRepo, nil)
	importer := service_github.ImporterQueue(workers_github.NopImporter())
	cloner := service_github.ClonerQueue(workers_github.NopCloner())
	gitHubService := service_github.New(
		logger,
		gitHubRepositoryRepo,
		gitHubInstallationRepo,
		nil, // gitHubUserRepo
		nil, // gitHubPullRequestRepo
		config.GitHubAppConfig{},
		nil, // gitHubClientProvider
		nil,
		&importer,
		&cloner,
		workspaceWriter,
		workspaceRepo,
		codebaseUserRepo,
		nil,
		executorProvider,
		gitSnapshotter,
		nil, // postHogClient
		nil, // notificationSender
		nil, // eventsSender
		userService,
	)
	workspaceService := service_workspace.New(
		logger,
		postHogClient,

		workspaceWriter,
		workspaceRepo,

		userRepo,
		reviewRepo,

		commentsService,
		changeService,
		gitHubService,

		activitySender,
		executorProvider,
		eventsSender,
		snapshotPublisher,
		gitSnapshotter,
		buildQueue,
	)

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

	authService := service_auth.New(codebaseService, userService, workspaceService, nil /*aclProvider*/, nil /*organizationService*/)

	createCodebaseRoute := routes_v3_codebase.Create(logger, codebaseRepo, codebaseUserRepo, executorProvider, postHogClient, eventsSender, workspaceService)
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
		executorProvider,
		logger,
	)

	fileRootResolver := graphql_file.NewFileRootResolver(executorProvider)

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
		&changeRootResolver,
		nil, // workspaceActivityResolver
		&reviewRootResolver,
		nil, // presenseRootResolver
		nil, // suggestitonsRootResolver
		statusesRootResolver,
		workspaceWatchersRootResolver,
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
		&changeRootResolver,
		&fileRootResolver,
		nil,
		nil, // instantIntegrationRootResolver

		logger,
		nil,
		nil,
		postHogClient,
		executorProvider,

		authService,
	)

	*statusesRootResolver = graphql_statuses.New(
		logger,
		statusesServcie,
		changeService,
		workspaceService,
		authService,
		gitHubPRRepo,
		&changeRootResolver,
		nil, // github pr resolver
		viewEvents,
	)

	createUser := user.User{ID: uuid.New().String(), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
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
		true,
		"file://../db/migrations",
	)
	if err != nil {
		panic(err)
	}

	logger := zap.NewNop()
	postHogClient := posthog.NewFakeClient()

	userRepo := db_user.NewRepo(d)
	codebaseRepo := db_codebase.NewRepo(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	workspaceRepo := db_workspace.NewRepo(d)
	viewRepo := db_view.NewRepo(d)
	changeRepo := db_change.NewRepo(d)
	changeCommitRepo := db_change.NewCommitRepository(d)
	snapshotRepo := db_snapshots.NewRepo(d)
	gitHubRepositoryRepo := db_github.NewGitHubRepositoryRepo(d)
	gitHubInstallationRepo := db_github.NewGitHubInstallationRepo(d)
	gitHubPRRepo := db_github.NewGitHubPRRepo(d)
	commentRepo := db_comments.NewRepo(d)
	workspaceActivityRepo := db_activity.NewActivityRepo(d)
	workspaceActivityReadsRepo := db_activity.NewActivityReadsRepo(d)
	viewEvents := events.NewInMemory()
	executorProvider := executor.NewProvider(logger, repoProvider)
	viewWorkspaceSnapshotsRepo := view_workspace_snapshot.NewRepo(d)
	reviewRepo := db_review.NewReviewRepository(d)
	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, viewEvents)
	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotRepo, workspaceRepo, workspaceRepo, viewRepo, eventsSender, executorProvider, logger)
	snapshotPublisher := worker_snapshots.NewSync(gitSnapshotter)
	activityService := service_activity.New(workspaceActivityReadsRepo, eventsSender)
	activitySender := activity_sender.NewActivitySender(codebaseUserRepo, workspaceActivityRepo, activityService, eventsSender)
	suggestionRepo := db_suggestion.New(d)
	notificationSender := sender.NewNoopNotificationSender()
	commentsService := service_comments.New(commentRepo)

	queue := queue.NewNoop()
	buildQueue := workers_ci.New(zap.NewNop(), queue, nil)
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo)
	userService := service_user.New(zap.NewNop(), userRepo, nil /*jwtService*/, nil /*onetime*/, nil /*emailsender*/, postHogClient)

	workspaceWriter := ws_meta.NewWriterWithEvents(logger, workspaceRepo, eventsSender)
	changeService := service_change.New(executorProvider, nil, nil, userRepo, changeRepo, changeCommitRepo, nil)
	importer := service_github.ImporterQueue(workers_github.NopImporter())
	cloner := service_github.ClonerQueue(workers_github.NopCloner())
	gitHubService := service_github.New(
		logger,
		gitHubRepositoryRepo,
		gitHubInstallationRepo,
		nil, // gitHubUserRepo
		nil, // gitHubPullRequestRepo
		config.GitHubAppConfig{},
		nil, // gitHubClientProvider
		nil,
		&importer,
		&cloner,
		workspaceWriter,
		workspaceRepo,
		codebaseUserRepo,
		nil,
		executorProvider,
		gitSnapshotter,
		nil, // postHogClient
		nil, // notificationSender
		nil, // eventsSender
		nil,
	)
	workspaceService := service_workspace.New(
		logger,
		postHogClient,

		workspaceWriter,
		workspaceRepo,

		userRepo,
		reviewRepo,

		commentsService,
		changeService,
		gitHubService,

		activitySender,
		executorProvider,
		eventsSender,
		snapshotPublisher,
		gitSnapshotter,
		buildQueue,
	)

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

	authService := service_auth.New(codebaseService, userService, workspaceService, nil /*aclProvider*/, nil /*organizationService*/)

	createCodebaseRoute := routes_v3_codebase.Create(logger, codebaseRepo, codebaseUserRepo, executorProvider, postHogClient, eventsSender, workspaceService)
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
		executorProvider,
		logger,
	)

	fileRootResolver := graphql_file.NewFileRootResolver(executorProvider)

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
		&changeRootResolver,
		nil, // workspaceActivityResolver
		&reviewRootResolver,
		nil, // presenseRootResolver
		nil, // suggestitonsRootResolver
		statusesRootResolver,
		workspaceWatchersRootResolver,
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
		viewWorkspaceSnapshotsRepo,
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
		&changeRootResolver,
		&fileRootResolver,
		nil,
		nil, // instantIntegrationRootResolver

		logger,
		nil,
		nil,
		postHogClient,
		executorProvider,

		authService,
	)

	*statusesRootResolver = graphql_statuses.New(
		logger,
		statusesServcie,
		changeService,
		workspaceService,
		authService,
		gitHubPRRepo,
		&changeRootResolver,
		nil, // github pr resolver
		viewEvents,
	)

	createUser := user.User{ID: uuid.New().String(), Name: "Test", Email: uuid.New().String() + "@getsturdy.com"}
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

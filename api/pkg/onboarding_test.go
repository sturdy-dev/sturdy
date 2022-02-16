	"go.uber.org/dig"
	"getsturdy.com/api/pkg/analytics"
	module_api "getsturdy.com/api/pkg/api/module"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	module_github "getsturdy.com/api/pkg/github/module"
	gqldataloader "getsturdy.com/api/pkg/graphql/dataloader"
	module_snapshots "getsturdy.com/api/pkg/snapshots/module"
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

	type deps struct {
		dig.In
		UserRepo              db_user.Repository
		CodebaseRootResolver  resolvers.CodebaseRootResolver
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		UserRootResolver      resolvers.UserRootResolver
		CommentsRootResolver  resolvers.CommentRootResolver
		ViewRootResolver      resolvers.ViewRootResolver
		GcService             *service_gc.Service
		CodebaseService       *service_codebase.Service
		WorkspaceService      service_workspace.Service
		GitSnapshotter        snapshotter.Snapshotter
		RepoProvider          provider.RepoProvider

		// Dependencies of Gin Routes
		CodebaseUserRepo db_codebase.CodebaseUserRepository
		WorkspaceRepo    db_workspace.Repository
		ViewRepo         db_view.Repository
		SnapshotRepo     db_snapshots.Repository
		ExecutorProvider executor.Provider
		EventsSender     events.EventSender

		Logger          *zap.Logger
		AnalyticsClient analytics.Client
	var d deps
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	}
	userRepo := d.UserRepo
	codebaseRootResolver := d.CodebaseRootResolver
	workspaceRootResolver := d.WorkspaceRootResolver
	userRootResolver := d.UserRootResolver
	commentsRootResolver := d.CommentsRootResolver
	viewRootResolver := d.ViewRootResolver
	gcService := d.GcService
	codebaseService := d.CodebaseService
	workspaceService := d.WorkspaceService
	gitSnapshotter := d.GitSnapshotter
	repoProvider := d.RepoProvider

	logger := d.Logger
	analyticsClient := d.AnalyticsClient
	codebaseUserRepo := d.CodebaseUserRepo
	workspaceRepo := d.WorkspaceRepo
	viewRepo := d.ViewRepo
	snapshotRepo := d.SnapshotRepo
	executorProvider := d.ExecutorProvider
	eventsSender := d.EventsSender
	createViewRoute := routes_v3_view.Create(logger, viewRepo, codebaseUserRepo, analyticsClient, workspaceRepo, gitSnapshotter, snapshotRepo, workspaceRepo, executorProvider, eventsSender)
	err := ioutil.WriteFile(path.Join(viewPath, "test.txt"), []byte("hello\n"), 0o666)
		err := gcService.WorkWithOptions(context.Background(), logger, codebaseRes.ID, 0, 0)
	type deps struct {
		dig.In
		UserRepo              db_user.Repository
		WorkspaceRootResolver resolvers.WorkspaceRootResolver
		CodebaseService       *service_codebase.Service
		WorkspaceService      service_workspace.Service
		GitSnapshotter        snapshotter.Snapshotter
		RepoProvider          provider.RepoProvider

		// Dependencies of Gin Routes
		CodebaseUserRepo db_codebase.CodebaseUserRepository
		WorkspaceRepo    db_workspace.Repository
		ViewRepo         db_view.Repository
		SnapshotRepo     db_snapshots.Repository
		ExecutorProvider executor.Provider
		EventsSender     events.EventSender
		WorkspaceWriter  db_workspace.WorkspaceWriter

		Logger          *zap.Logger
		AnalyticsClient analytics.Client
	var d deps
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	userRepo := d.UserRepo
	repoProvider := d.RepoProvider
	workspaceService := d.WorkspaceService
	workspaceRootResolver := d.WorkspaceRootResolver
	createCodebaseRoute := routes_v3_codebase.Create(d.Logger, d.CodebaseService)
	createWorkspaceRoute := routes_v3_workspace.Create(d.Logger, d.WorkspaceService, d.CodebaseUserRepo)
	createViewRoute := routes_v3_view.Create(d.Logger, d.ViewRepo, d.CodebaseUserRepo, d.AnalyticsClient, d.WorkspaceRepo, d.GitSnapshotter, d.SnapshotRepo, d.WorkspaceWriter, d.ExecutorProvider, d.EventsSender)
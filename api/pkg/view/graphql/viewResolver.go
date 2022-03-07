package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	vcs2 "getsturdy.com/api/pkg/snapshots/vcs"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/view/ignore"
	"getsturdy.com/api/pkg/view/open"
	view_vcs "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace_watchers "getsturdy.com/api/pkg/workspaces/watchers/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	concurrentUpdatedViewConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "sturdy_graphql_concurrent_subscriptions",
		ConstLabels: prometheus.Labels{"subscription": "updatedView"},
	})
)

type ViewRootResolver struct {
	viewRepo                 db_view.Repository
	workspaceReader          db_workspaces.WorkspaceReader
	snapshotter              snapshotter.Snapshotter
	snapshotRepo             db_snapshots.Repository
	authorResolver           resolvers.AuthorRootResolver
	workspaceResolver        *resolvers.WorkspaceRootResolver
	workspaceWriter          db_workspaces.WorkspaceWriter
	viewEvents               events.EventReader
	eventSender              events.EventSender
	eventSenderV2            *eventsv2.Publisher
	executorProvider         executor.Provider
	logger                   *zap.Logger
	viewStatusRootResolver   resolvers.ViewStatusRootResolver
	workspaceWatchersService *service_workspace_watchers.Service
	codebaseResolver         resolvers.CodebaseRootResolver
	authService              *service_auth.Service
	analyticsService         *service_analytics.Service
	eventsSubscriber         *eventsv2.Subscriber
}

func NewResolver(
	viewRepo db_view.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	snapshotter snapshotter.Snapshotter,
	snapshotRepo db_snapshots.Repository,
	authorResolver resolvers.AuthorRootResolver,
	workspaceResolver *resolvers.WorkspaceRootResolver,
	workspaceWriter db_workspaces.WorkspaceWriter,
	viewEvents events.EventReader,
	eventSender events.EventSender,
	eventSenderV2 *eventsv2.Publisher,
	executorProvider executor.Provider,
	logger *zap.Logger,
	viewStatusRootResolver resolvers.ViewStatusRootResolver,
	workspaceWatchersService *service_workspace_watchers.Service,
	analyticsService *service_analytics.Service,
	codebaseResolver resolvers.CodebaseRootResolver,
	authService *service_auth.Service,
	eventsSubscriber *eventsv2.Subscriber,
) resolvers.ViewRootResolver {
	return &ViewRootResolver{
		viewRepo:                 viewRepo,
		workspaceReader:          workspaceReader,
		snapshotter:              snapshotter,
		snapshotRepo:             snapshotRepo,
		authorResolver:           authorResolver,
		workspaceResolver:        workspaceResolver,
		workspaceWriter:          workspaceWriter,
		viewEvents:               viewEvents,
		eventSender:              eventSender,
		eventSenderV2:            eventSenderV2,
		executorProvider:         executorProvider,
		logger:                   logger,
		viewStatusRootResolver:   viewStatusRootResolver,
		workspaceWatchersService: workspaceWatchersService,
		analyticsService:         analyticsService,
		codebaseResolver:         codebaseResolver,
		authService:              authService,
		eventsSubscriber:         eventsSubscriber,
	}
}

func (r *ViewRootResolver) View(ctx context.Context, args resolvers.ViewArgs) (resolvers.ViewResolver, error) {
	return r.resolveView(ctx, args.ID)
}

func (r *ViewRootResolver) InternalViewsByUser(userID users.ID) ([]resolvers.ViewResolver, error) {
	views, err := r.viewRepo.ListByUser(userID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	res := make([]resolvers.ViewResolver, 0, len(views))
	for _, vw := range views {
		res = append(res, &Resolver{v: vw, root: r})
	}
	return res, nil
}

func (r *ViewRootResolver) InternalLastUsedViewByUser(ctx context.Context, codebaseID string, userID users.ID) (resolvers.ViewResolver, error) {
	vw, err := r.viewRepo.LastUsedByCodebaseAndUser(ctx, codebaseID, userID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	if vw == nil {
		return nil, gqlerrors.ErrNotFound
	}
	return &Resolver{v: vw, root: r}, nil
}

func (r *ViewRootResolver) resolveView(ctx context.Context, id graphql.ID) (resolvers.ViewResolver, error) {
	v, err := r.viewRepo.Get(string(id))
	if err != nil {
		return nil, gqlerrors.Error(err, "description", "view not found")
	}

	if err := r.authService.CanRead(ctx, v); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &Resolver{v, r}, nil
}

func (r *ViewRootResolver) resolveViewObj(ctx context.Context, v *view.View) (resolvers.ViewResolver, error) {
	if err := r.authService.CanRead(ctx, v); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &Resolver{v, r}, nil
}

func (r *ViewRootResolver) UpdatedViews(ctx context.Context) (chan resolvers.ViewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make(chan resolvers.ViewResolver, 100)

	callback := func(_ context.Context, eventView *view.View) error {
		resolver, err := r.resolveViewObj(ctx, eventView)
		if err != nil {
			return err
		}
		select {
		case res <- resolver:
		default:
			r.logger.Error("dropped updatedView event")
		}
		return nil
	}

	r.eventsSubscriber.User(ctx, userID).OnViewUpdated(ctx, callback)
	r.eventsSubscriber.User(ctx, userID).OnViewStatusUpdated(ctx, callback)

	go func() {
		<-ctx.Done()
		close(res)
	}()

	return res, nil
}

func (r *ViewRootResolver) UpdatedView(ctx context.Context, args resolvers.UpdatedViewArgs) (chan resolvers.ViewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	v, err := r.viewRepo.Get(string(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(err, "description", "view not found")
	}

	if err := r.authService.CanWrite(ctx, v); err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make(chan resolvers.ViewResolver, 100)

	concurrentUpdatedViewConnections.Inc()

	callback := func(_ context.Context, eventView *view.View) error {
		if v.ID != eventView.ID {
			return nil
		}
		resolver, err := r.resolveViewObj(ctx, eventView)
		if err != nil {
			return err
		}
		select {
		case res <- resolver:
		default:
			r.logger.Error("dropped updatedView event")
		}

		return nil
	}

	r.eventsSubscriber.User(ctx, userID).OnViewUpdated(ctx, callback)
	r.eventsSubscriber.User(ctx, userID).OnViewStatusUpdated(ctx, callback)

	go func() {
		<-ctx.Done()
		close(res)
		concurrentUpdatedViewConnections.Dec()
	}()

	return res, nil
}

func (r *ViewRootResolver) RepairView(ctx context.Context, args struct{ ID graphql.ID }) (resolvers.ViewResolver, error) {
	vw, err := r.viewRepo.Get(string(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, vw); err != nil {
		return nil, gqlerrors.Error(err)
	}

	ws, err := r.workspaceReader.Get(vw.WorkspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	err = r.executorProvider.New().Schedule(func(repoProvider provider.RepoProvider) error {
		var restoreWs *workspaces.Workspace
		// This view is the authoritative view of a workspace, restore the workspace afterwards
		if ws.ViewID != nil && *ws.ViewID == vw.ID {
			restoreWs = ws
		}

		if err := recreateView(repoProvider, vw, restoreWs, r.logger, r.snapshotter); err != nil {
			return err
		}
		return nil
	}).ExecView(ws.CodebaseID, vw.ID, "repairView")
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.resolveView(ctx, args.ID)
}

func (r *ViewRootResolver) CreateView(ctx context.Context, args resolvers.CreateViewArgs) (resolvers.ViewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	var mountPath *string
	if args.Input.MountPath != "" {
		mountPath = &args.Input.MountPath
	}

	var mountHostname *string
	if args.Input.MountHostname != "" {
		mountHostname = &args.Input.MountHostname
	}

	t := time.Now()
	e := view.View{
		ID:            uuid.New().String(),
		UserID:        userID,
		CodebaseID:    ws.CodebaseID,
		WorkspaceID:   ws.ID,
		MountPath:     mountPath,     // It's optional
		MountHostname: mountHostname, // It's optional
		CreatedAt:     &t,
	}

	if err := r.viewRepo.Create(e); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err = r.executorProvider.New().
		AllowRebasingState(). // allowed because the view does not exist yet
		Schedule(func(repoProvider provider.RepoProvider) error {
			return view_vcs.Create(repoProvider, ws.CodebaseID, ws.ID, e.ID)
		}).ExecView(ws.CodebaseID, e.ID, "createView"); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Use workspace on view
	if err := open.OpenWorkspaceOnView(ctx, r.logger, &e, ws, r.viewRepo, r.workspaceReader, r.snapshotter, r.snapshotRepo, r.workspaceWriter, r.executorProvider, r.eventSenderV2); errors.Is(err, open.ErrRebasing) {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "View is currently in rebasing state. Please resolve all the conflicts and try again.")
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	}

	r.analyticsService.Capture(ctx, "create view",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
		analytics.Property("view_id", e.ID),
		analytics.Property("mount_path", e.MountPath),
		analytics.Property("mount_hostname", e.MountHostname),
	)

	return r.resolveView(ctx, graphql.ID(e.ID))
}

func recreateView(repoProvider provider.RepoProvider, vw *view.View, ws *workspaces.Workspace, logger *zap.Logger, gitSnapshotter snapshotter.Snapshotter) error {
	trunkPath := repoProvider.TrunkPath(vw.CodebaseID)
	newView := vw.ID + "-recreate-" + uuid.NewString()
	newViewPath := repoProvider.ViewPath(vw.CodebaseID, newView)
	backupPath := repoProvider.ViewPath(vw.CodebaseID, vw.ID+"-replaced-"+uuid.NewString())

	decoratedLogger := logger.Named("recreateView").With(zap.String("new_view_path", newViewPath), zap.String("backup_path", backupPath))

	decoratedLogger.Info("recreating view")

	if _, err := vcs.CloneRepo(trunkPath, newViewPath); err != nil {
		return err
	}

	if ws != nil {
		if err := view_vcs.SetWorkspace(repoProvider, vw.CodebaseID, newView, ws.ID); err != nil {
			return err
		}

		// Attempt to make a snapshot of the existing view
		snapshot, err := gitSnapshotter.Snapshot(vw.CodebaseID, vw.WorkspaceID, snapshots.ActionPreCheckoutOtherView, snapshotter.WithOnView(vw.ID))
		if err != nil {
			return err
		}

		// TODO: What if we can't make a snapshot because the view is FUBAR?
		if err := vcs2.Restore(logger, repoProvider, ws.CodebaseID, ws.ID, newView, snapshot.ID, snapshot.CommitID); err != nil {
			return fmt.Errorf("failed to restore snapshot: %w", err)
		}
		decoratedLogger.Info("restored from snapshot", zap.String("commit_id", snapshot.CodebaseID))
	}

	// Swap replacement
	if err := os.Rename(repoProvider.ViewPath(vw.CodebaseID, vw.ID), backupPath); err != nil {
		return err
	}

	if err := os.Rename(newViewPath, repoProvider.ViewPath(vw.CodebaseID, vw.ID)); err != nil {
		return err
	}

	decoratedLogger.Info("repair view completed")

	return nil
}

type Resolver struct {
	v    *view.View
	root *ViewRootResolver
}

func (r *Resolver) ID() graphql.ID {
	return graphql.ID(r.v.ID)
}

func (r *Resolver) MountPath() string {
	if r.v.MountPath == nil {
		return ""
	}
	return *r.v.MountPath
}

var homeDirPattern = regexp.MustCompile("^/([Uu]sers|home)/[^/]+")

func (r *Resolver) ShortMountPath() string {
	return homeDirPattern.ReplaceAllString(r.MountPath(), "~")
}

func (r *Resolver) MountHostname() string {
	if r.v.MountHostname == nil {
		return ""
	}
	return *r.v.MountHostname
}

func (r *Resolver) LastUsedAt() int32 {
	if r.v.LastUsedAt == nil {
		return 0
	}
	return int32(r.v.LastUsedAt.Unix())
}

func (r *Resolver) CreatedAt() int32 {
	if r.v.CreatedAt == nil {
		return 0
	}
	return int32(r.v.CreatedAt.Unix())
}

func (r *Resolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	return r.root.authorResolver.Author(ctx, graphql.ID(r.v.UserID))
}

func (r *Resolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	ws, err := r.root.workspaceReader.GetByViewID(r.v.ID, false)
	switch {
	case err == nil:
		return (*r.root.workspaceResolver).Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *Resolver) Status(ctx context.Context) (resolvers.ViewStatusResolver, error) {
	res, err := r.root.viewStatusRootResolver.InternalViewStatus(ctx, r.v.ID)
	switch {
	case err == nil:
		return res, nil
	case errors.Is(err, gqlerrors.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (r *Resolver) Codebase(ctx context.Context) (resolvers.CodebaseResolver, error) {
	id := graphql.ID(r.v.CodebaseID)
	cb, err := r.root.codebaseResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return cb, nil
}

func (r *Resolver) IgnoredPaths(ctx context.Context) ([]string, error) {
	var res []string

	err := r.root.executorProvider.New().
		AllowRebasingState(). // allowed to parse .gitignore even if rebasing
		Read(func(repo vcs.RepoReader) error {
			var err error
			res, err = ignore.FindIgnore(os.DirFS(repo.Path()))
			if err != nil {
				return err
			}
			return nil
		}).ExecView(r.v.CodebaseID, r.v.ID, "findIgnores")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Resolver) SuggestingWorkspace() resolvers.WorkspaceResolver {
	return nil
}

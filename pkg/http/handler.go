package http

import (
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	authz "mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	db_change "mash/pkg/change/db"
	routes_v3_change "mash/pkg/change/routes"
	service_change "mash/pkg/change/service"
	routes_ci "mash/pkg/ci/routes"
	service_ci "mash/pkg/ci/service"
	workers_ci "mash/pkg/ci/workers"
	provider_acl "mash/pkg/codebase/acl/provider"
	db_codebase "mash/pkg/codebase/db"
	routes_v3_codebase "mash/pkg/codebase/routes"
	service_codebase "mash/pkg/codebase/service"
	db_comments "mash/pkg/comments/db"
	service_comments "mash/pkg/comments/service"
	worker_gc "mash/pkg/gc/worker"
	"mash/pkg/ginzap"
	ghappclient "mash/pkg/github/client"
	"mash/pkg/github/config"
	db_github "mash/pkg/github/db"
	routes_v3_ghapp "mash/pkg/github/routes"
	service_github "mash/pkg/github/service"
	workers_github "mash/pkg/github/workers"
	sturdygrapql "mash/pkg/graphql"
	service_buildkite "mash/pkg/integrations/buildkite/service"
	service_jwt "mash/pkg/jwt/service"
	"mash/pkg/metrics/ginprometheus"
	db_mutagen "mash/pkg/mutagen/db"
	routes_v3_mutagen "mash/pkg/mutagen/routes"
	db_newsletter "mash/pkg/newsletter/db"
	routes_v3_newsletter "mash/pkg/newsletter/routes"
	db_notification "mash/pkg/notification/db"
	notification_sender "mash/pkg/notification/sender"
	service_notification "mash/pkg/notification/service"
	db_onboarding "mash/pkg/onboarding/db"
	db_pki "mash/pkg/pki/db"
	routes_v3_pki "mash/pkg/pki/routes"
	service_presence "mash/pkg/presence/service"
	db_review "mash/pkg/review/db"
	service_servicetokens "mash/pkg/servicetokens/service"
	db_snapshots "mash/pkg/snapshots/db"
	"mash/pkg/snapshots/snapshotter"
	worker_snapshots "mash/pkg/snapshots/worker"
	service_statuses "mash/pkg/statuses/service"
	"mash/pkg/stream/routes"
	service_suggestion "mash/pkg/suggestions/service"
	routes_v3_sync "mash/pkg/sync/routes"
	service_sync "mash/pkg/sync/service"
	db_user "mash/pkg/user/db"
	routes_v3_user "mash/pkg/user/routes"
	service_user "mash/pkg/user/service"
	view_auth "mash/pkg/view/auth"
	db_view "mash/pkg/view/db"
	"mash/pkg/view/events"
	meta_view "mash/pkg/view/meta"
	routes_v3_view "mash/pkg/view/routes"
	"mash/pkg/view/view_workspace_snapshot"
	"mash/pkg/waitinglist"
	"mash/pkg/waitinglist/acl"
	"mash/pkg/waitinglist/instantintegration"
	db_activity "mash/pkg/workspace/activity/db"
	activity_sender "mash/pkg/workspace/activity/sender"
	service_activity "mash/pkg/workspace/activity/service"
	db_workspace "mash/pkg/workspace/db"
	routes_v3_workspace "mash/pkg/workspace/routes"
	service_workspace "mash/pkg/workspace/service"
	service_workspace_watchers "mash/pkg/workspace/watchers/service"
	"mash/vcs/executor"

	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

func NewHandler(
	logger *zap.Logger,
	userRepo db_user.Repository,
	postHogClient posthog.Client,
	waitingListRepo waitinglist.WaitingListRepo,
	aclInterestRepo acl.ACLInterestRepo,
	instantIntegrationInterestRepo instantintegration.InstantIntegrationInterestRepo,
	codebaseRepo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	viewRepo db_view.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	userPublicKeyRepo db_pki.Repo,
	snapshotterQueue worker_snapshots.Queue,
	snapshotRepo db_snapshots.Repository,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,
	commentRepo db_comments.Repository,
	codebaseViewEvents events.EventReadWriter,
	gcQueue *worker_gc.Queue,
	gitSnapshotter snapshotter.Snapshotter,
	viewWorkspaceSnapshotsRepo view_workspace_snapshot.Repository,
	gitHubInstallationsRepo db_github.GitHubInstallationRepo,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubUserRepo db_github.GitHubUserRepo,
	gitHubPRRepo db_github.GitHubPRRepo,
	gitHubAppConfig config.GitHubAppConfig,
	notificationRepository db_notification.Repository,
	notificationSender notification_sender.NotificationSender,
	gitHubClientProvider ghappclient.ClientProvider,
	gitHubPersonalClientProvider ghappclient.PersonalClientProvider,
	gitHubClonerPublisher *workers_github.ClonerQueue,
	workspaceWriter db_workspace.WorkspaceWriter,
	viewUpdatedFunc meta_view.ViewUpdatedFunc,
	executorProvider executor.Provider,
	viewStatusRepo db_mutagen.ViewStatusRepository,
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository,
	workspaceActivityRepo db_activity.ActivityRepository,
	aclProvider *provider_acl.Provider,
	reviewRepo db_review.ReviewRepository,
	workspaceActivityReadsRepo db_activity.ActivityReadsRepository,
	activitySender activity_sender.ActivitySender,
	eventSender events.EventSender,
	presenceService service_presence.Service,
	suggestionService *service_suggestion.Service,
	workspaceService service_workspace.Service,
	notificationPreferencesService *service_notification.Preferences,
	changeService *service_change.Service,
	userService *service_user.Service,
	ciService *service_ci.Service,
	statusesService *service_statuses.Service,
	completedOnboardingStepsRepo db_onboarding.CompletedOnboardingStepsRepository,
	syncService *service_sync.Service,
	workspaceWatchersService *service_workspace_watchers.Service,
	jwtService *service_jwt.Service,
	activityService *service_activity.Service,
	commentsService *service_comments.Service,
	gitHubService *service_github.Service,
	codebaseService *service_codebase.Service,
	servicetokensService *service_servicetokens.Service,
	buildkiteService *service_buildkite.Service,
	authService *service_auth.Service,
	ciBuildQueue *workers_ci.BuildQueue,
	developmentAllowExtraCorsOrigin string,
) http.Handler {
	logger = logger.With(zap.String("component", "http"))

	allowOrigins := []string{
		// Production
		"https://getsturdy.com",

		// Development
		"http://localhost:8080",

		// Probably unused
		"https://driva.dev",

		// Staging environments for the website
		"https://gustav-staging.driva.dev",
		"https://gustav-staging.getsturdy.com",
	}
	if developmentAllowExtraCorsOrigin != "" {
		logger.Info("adding CORS origin", zap.String("origin", developmentAllowExtraCorsOrigin))
		allowOrigins = append(allowOrigins, developmentAllowExtraCorsOrigin)
	}
	cors := ginCors.New(ginCors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"POST, OPTIONS, GET, PUT, DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, Cookie", "sentry-trace"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	r := gin.New()
	r.Use(accessLogger(logger, time.RFC3339, true))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(cors)

	// Setup Prometheus metrics for Gin
	ginprom := ginprometheus.NewPrometheus("gin", logger)
	ginprom.ReqCntURLLabelMappingFn = metricsMapper
	ginprom.Use(r)

	// GraphQL
	sg := sturdygrapql.New(logger,
		codebaseRepo,
		codebaseUserRepo,
		workspaceReader,
		userRepo,
		viewRepo,
		gitSnapshotter,
		viewWorkspaceSnapshotsRepo,
		snapshotRepo,
		commentRepo,
		codebaseViewEvents,
		changeRepo,
		changeCommitRepo,
		notificationRepository,
		notificationSender,
		gitHubUserRepo,
		gitHubPRRepo,
		gitHubInstallationsRepo,
		gitHubAppConfig,
		gitHubRepositoryRepo,
		gitHubClientProvider,
		gitHubPersonalClientProvider,
		postHogClient,
		workspaceWriter,
		executorProvider,
		viewStatusRepo,
		notificationSettingsRepo,
		workspaceActivityRepo,
		aclProvider,
		reviewRepo,
		workspaceActivityReadsRepo,
		activitySender,
		eventSender,
		gitSnapshotter,
		presenceService,
		suggestionService,
		workspaceService,
		notificationPreferencesService,
		changeService,
		userService,
		ciService,
		statusesService,
		completedOnboardingStepsRepo,
		workspaceWatchersService,
		jwtService,
		activityService,
		userPublicKeyRepo,
		gitHubService,
		codebaseService,
		servicetokensService,
		buildkiteService,
		authService,
	)
	graphql := r.Group("/graphql", sturdygrapql.CorsMiddleware(allowOrigins), authz.GinMiddleware(logger, jwtService))
	graphql.OPTIONS("", func(c *gin.Context) { c.Status(http.StatusOK) })
	graphql.OPTIONS("ws", func(c *gin.Context) { c.Status(http.StatusOK) })
	graphql.POST("", sg.HttpHandler())
	graphql.GET("ws", sg.WebsocketHandler())
	graphql.POST("ws", sg.WebsocketHandler())

	// Public endpoints, no authentication required
	publ := r.Group("")

	// Private endpoints, requires a valid auth cookie
	auth := r.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))

	publ.POST("/v3/auth", routes_v3_user.Auth(logger, userRepo, postHogClient, jwtService))
	publ.POST("/v3/auth/destroy", routes_v3_user.AuthDestroy)
	publ.POST("/v3/auth/magic-link/send", routes_v3_user.SendMagicLink(logger, userService))
	publ.POST("/v3/auth/magic-link/verify", routes_v3_user.VerifyMagicLink(logger, userService, jwtService))
	auth.POST("/v3/auth/client-token", routes_v3_user.ClientToken(userRepo, jwtService))
	auth.POST("/v3/auth/renew-token", routes_v3_user.RenewToken(logger, userRepo, jwtService))

	auth.POST("/v3/users/verify-email", routes_v3_user.SendEmailVerification(logger, userService)) // Used by the web (2021-11-14)
	auth.POST("/v3/user/update-avatar", routes_v3_user.UpdateAvatar(userRepo))                     // Used by the web (2021-10-04)
	auth.GET("/v3/user", routes_v3_user.GetSelf(userRepo, jwtService))                             // Used by the command line client

	auth.POST("/v3/codebases", routes_v3_codebase.Create(logger, codebaseRepo, codebaseUserRepo, executorProvider, postHogClient, eventSender, workspaceService)) // Used by the web (2021-10-04)
	auth.GET("/v3/codebases/:id", routes_v3_codebase.Get(codebaseRepo, codebaseUserRepo, logger, userRepo, executorProvider))                                     // Used by the command line client
	auth.POST("/v3/codebases/:id/invite", routes_v3_codebase.Invite(userRepo, codebaseUserRepo, codebaseService, authService, eventSender, logger))               // Used by the web (2021-10-04)
	publ.GET("/v3/join/get-codebase/:code", routes_v3_codebase.JoinGetCodebase(logger, codebaseRepo))                                                             // Used by the web (2021-10-04)
	auth.POST("/v3/join/codebase/:code", routes_v3_codebase.JoinCodebase(logger, codebaseRepo, codebaseUserRepo, eventSender))                                    // Used by the web (2021-10-04)

	auth.POST("/v3/views", routes_v3_view.Create(logger, viewRepo, codebaseUserRepo, postHogClient, workspaceReader, gitSnapshotter, snapshotRepo, workspaceWriter, executorProvider, eventSender)) // Used by the command line client
	authedViews := auth.Group("/v3/views/:viewID", view_auth.ValidateViewAccessMiddleware(authService, viewRepo))
	authedViews.GET("", routes_v3_view.Get(viewRepo, workspaceReader, userRepo, logger))                                                 // Used by the command line client
	authedViews.POST("/ignore-file", routes_v3_change.IgnoreFile(logger, viewRepo, codebaseUserRepo, executorProvider, viewUpdatedFunc)) // Used by the web (2021-10-04)
	authedViews.GET("/ignores", routes_v3_view.Ignores(logger, executorProvider, viewRepo))                                              // Called from client-side sturdy-cli

	auth.GET("/v3/stream", routes.Stream(logger, viewRepo, codebaseViewEvents, workspaceReader, authService, workspaceService, suggestionService)) // Used by the web (2021-10-04)

	rebase := auth.Group("/v3/rebase/")
	rebase.Use(view_auth.ValidateViewAccessMiddleware(authService, viewRepo))
	rebase.GET(":viewID", routes_v3_sync.Status(viewRepo, executorProvider, logger)) // Used by the web (2021-10-04)
	rebase.POST(":viewID/start", routes_v3_sync.StartV2(logger, syncService))        // Used by the web (2021-10-25)
	rebase.POST(":viewID/resolve", routes_v3_sync.ResolveV2(logger, syncService))    // Used by the web (2021-10-25)

	auth.POST("/v3/changes/:id/update", routes_v3_change.Update(logger, codebaseUserRepo, postHogClient, changeRepo)) // Used by the web (2021-10-04)

	auth.POST("/v3/workspaces", routes_v3_workspace.Create(logger, workspaceService, codebaseUserRepo)) // Used by the command line client
	// Used by LBS to check for health
	publ.GET("/readyz", func(c *gin.Context) { c.Status(http.StatusOK) })

	publ.POST("/v3/waitinglist", waitinglist.Insert(logger, postHogClient, waitingListRepo))                               // Used by the web (2021-10-04)
	publ.POST("/v3/acl-request-enterprise", acl.Insert(logger, postHogClient, aclInterestRepo))                            // Used by the web (2021-10-04)
	publ.POST("/v3/instant-integration", instantintegration.Insert(logger, postHogClient, instantIntegrationInterestRepo)) // Used by the web (2021-10-27)

	auth.POST("/v3/pki/add-public-key", routes_v3_pki.AddPublicKey(userPublicKeyRepo)) // Used by the command line client
	publ.POST("/v3/pki/verify", routes_v3_pki.Verify(userPublicKeyRepo))               // Used by the command line client

	publ.POST("/v3/mutagen/validate-view", routes_v3_mutagen.ValidateView(logger, viewRepo, postHogClient, eventSender))                                                                                      // Called from server-side mutagen
	publ.POST("/v3/mutagen/sync-transitions", routes_v3_mutagen.SyncTransitions(logger, snapshotterQueue, viewRepo, gcQueue, presenceService, snapshotRepo, workspaceReader, suggestionService, eventSender)) // Called from server-side mutagen
	publ.GET("/v3/mutagen/views/:id/allows", routes_v3_mutagen.ListAllows(logger, viewRepo, authService))                                                                                                     // Called form server-side mutagen
	publ.POST("/v3/mutagen/update-status", routes_v3_mutagen.UpdateStatus(logger, viewStatusRepo, viewRepo, eventSender))                                                                                     // Called from client-side mutagen
	auth.GET("/v3/mutagen/get-view/:id", routes_v3_mutagen.GetView(logger, viewRepo, codebaseUserRepo, codebaseRepo))                                                                                         // Called from client-side sturdy-cli

	publ.POST("/v3/github/webhook", routes_v3_ghapp.Webhook(logger, gitHubAppConfig, postHogClient, gitHubInstallationsRepo, gitHubRepositoryRepo, codebaseRepo, executorProvider, gitHubClientProvider, gitHubUserRepo, codebaseUserRepo, gitHubClonerPublisher, gitHubPRRepo, workspaceReader, workspaceWriter, workspaceService, syncService, changeRepo, changeCommitRepo, reviewRepo, eventSender, activitySender, statusesService, commentsService, gitHubService, ciBuildQueue))
	auth.POST("/v3/github/oauth", routes_v3_ghapp.Oauth(logger, gitHubAppConfig, userRepo, gitHubUserRepo, gitHubService))

	publ.POST("/v3/unsubscribe", routes_v3_newsletter.Unsubscribe(logger, userRepo, notificationSettingsRepo))

	publ.POST("/v3/statuses/webhook", routes_ci.WebhookHandler(logger, codebaseRepo, statusesService, ciService, servicetokensService, buildkiteService))

	return r
}

// accessLogger returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
//
// This code has been copied (and modified) from github.com/gin-contrib/zap.Ginzap
func accessLogger(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		// Log errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				if strings.Contains(e, "write: connection reset by peer") ||
					strings.Contains(e, "write: broken pipe") ||
					strings.Contains(e, "unexpected EOF") {
					logger.Warn(e)
				} else {
					logger.Error(e)
				}
			}
		}

		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(timeFormat)),
			zap.Duration("latency", latency),

			// Add Sturdy specific info
			zap.String("x-client-name", c.GetHeader("x-client-name")),
			zap.String("x-client-version", c.GetHeader("x-client-version")),
		)
	}
}

func metricsMapper(c *gin.Context) string {
	url := c.Request.URL.Path
	for _, p := range c.Params {
		url = strings.Replace(url, p.Value, ":"+p.Key, 1)
	}
	return url
}

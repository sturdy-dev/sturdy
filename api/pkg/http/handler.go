package http

import (
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"getsturdy.com/api/pkg/analytics"
	authz "getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_change "getsturdy.com/api/pkg/change/db"
	routes_v3_change "getsturdy.com/api/pkg/change/routes"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	routes_v3_codebase "getsturdy.com/api/pkg/codebase/routes"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/events"
	worker_gc "getsturdy.com/api/pkg/gc/worker"
	"getsturdy.com/api/pkg/ginzap"
	sturdygrapql "getsturdy.com/api/pkg/graphql"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/metrics/ginprometheus"
	db_mutagen "getsturdy.com/api/pkg/mutagen/db"
	routes_v3_mutagen "getsturdy.com/api/pkg/mutagen/routes"
	db_newsletter "getsturdy.com/api/pkg/newsletter/db"
	routes_v3_newsletter "getsturdy.com/api/pkg/newsletter/routes"
	db_pki "getsturdy.com/api/pkg/pki/db"
	routes_v3_pki "getsturdy.com/api/pkg/pki/routes"
	service_presence "getsturdy.com/api/pkg/presence/service"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	"getsturdy.com/api/pkg/stream/routes"
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	routes_v3_sync "getsturdy.com/api/pkg/sync/routes"
	service_sync "getsturdy.com/api/pkg/sync/service"
	db_user "getsturdy.com/api/pkg/users/db"
	routes_v3_user "getsturdy.com/api/pkg/users/routes"
	service_user "getsturdy.com/api/pkg/users/service"
	view_auth "getsturdy.com/api/pkg/view/auth"
	db_view "getsturdy.com/api/pkg/view/db"
	meta_view "getsturdy.com/api/pkg/view/meta"
	routes_v3_view "getsturdy.com/api/pkg/view/routes"
	"getsturdy.com/api/pkg/waitinglist"
	"getsturdy.com/api/pkg/waitinglist/acl"
	"getsturdy.com/api/pkg/waitinglist/instantintegration"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	routes_v3_workspace "getsturdy.com/api/pkg/workspace/routes"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/vcs/executor"

	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DevelopmentAllowExtraCorsOrigin string

type Engine gin.Engine

func ProvideHandler(
	logger *zap.Logger,
	userRepo db_user.Repository,
	analyticsClient analytics.Client,
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
	codebaseViewEvents events.EventReadWriter,
	gcQueue *worker_gc.Queue,
	gitSnapshotter snapshotter.Snapshotter,
	workspaceWriter db_workspace.WorkspaceWriter,
	viewUpdatedFunc meta_view.ViewUpdatedFunc,
	executorProvider executor.Provider,
	viewStatusRepo db_mutagen.ViewStatusRepository,
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository,
	eventSender events.EventSender,
	presenceService service_presence.Service,
	suggestionService *service_suggestion.Service,
	workspaceService service_workspace.Service,
	userService service_user.Service,
	syncService *service_sync.Service,
	jwtService *service_jwt.Service,
	codebaseService *service_codebase.Service,
	authService *service_auth.Service,
	developmentAllowExtraCorsOrigin DevelopmentAllowExtraCorsOrigin,
	grapqhlResolver *sturdygrapql.RootResolver,
) *Engine {
	logger = logger.With(zap.String("component", "http"))
	allowOrigins := []string{
		// Production
		"https://getsturdy.com",
		// Development
		"http://localhost:8080",
		// docker oneliner
		"http://localhost:30080",
		// Probably unused
		"https://driva.dev",
		// Staging environments for the website
		"https://gustav-staging.driva.dev",
		"https://gustav-staging.getsturdy.com",
	}
	if developmentAllowExtraCorsOrigin != "" {
		logger.Info("adding CORS origin", zap.String("origin", string(developmentAllowExtraCorsOrigin)))
		allowOrigins = append(allowOrigins, string(developmentAllowExtraCorsOrigin))
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
	graphql := r.Group("/graphql", sturdygrapql.CorsMiddleware(allowOrigins), authz.GinMiddleware(logger, jwtService))
	graphql.OPTIONS("", func(c *gin.Context) { c.Status(http.StatusOK) })
	graphql.OPTIONS("ws", func(c *gin.Context) { c.Status(http.StatusOK) })
	graphql.POST("", grapqhlResolver.HttpHandler())
	graphql.GET("ws", grapqhlResolver.WebsocketHandler())
	graphql.POST("ws", grapqhlResolver.WebsocketHandler())
	// Public endpoints, no authentication required
	publ := r.Group("")
	// Private endpoints, requires a valid auth cookie
	auth := r.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	publ.POST("/v3/auth", routes_v3_user.Login(logger, userRepo, analyticsClient, jwtService))
	publ.POST("/v3/users", routes_v3_user.Signup(logger, userService, jwtService, analyticsClient))
	publ.POST("/v3/auth/destroy", routes_v3_user.AuthDestroy)
	publ.POST("/v3/auth/magic-link/send", routes_v3_user.SendMagicLink(logger, userService))
	publ.POST("/v3/auth/magic-link/verify", routes_v3_user.VerifyMagicLink(logger, userService, jwtService))
	auth.POST("/v3/auth/client-token", routes_v3_user.ClientToken(userRepo, jwtService))
	auth.POST("/v3/auth/renew-token", routes_v3_user.RenewToken(logger, userRepo, jwtService))
	auth.POST("/v3/users/verify-email", routes_v3_user.SendEmailVerification(logger, userService))                                                                                                    // Used by the web (2021-11-14)
	auth.POST("/v3/user/update-avatar", routes_v3_user.UpdateAvatar(userRepo))                                                                                                                        // Used by the web (2021-10-04)
	auth.GET("/v3/user", routes_v3_user.GetSelf(userRepo, jwtService))                                                                                                                                // Used by the command line client
	auth.POST("/v3/codebases", routes_v3_codebase.Create(logger, codebaseService))                                                                                                                    // Used by the web (2021-10-04)
	auth.GET("/v3/codebases/:id", routes_v3_codebase.Get(codebaseRepo, codebaseUserRepo, logger, userRepo, executorProvider))                                                                         // Used by the command line client
	auth.POST("/v3/codebases/:id/invite", routes_v3_codebase.Invite(codebaseService, authService))                                                                                                    // No longer used (after 2022-01-31)
	publ.GET("/v3/join/get-codebase/:code", routes_v3_codebase.JoinGetCodebase(logger, codebaseRepo))                                                                                                 // Used by the web (2021-10-04)
	auth.POST("/v3/join/codebase/:code", routes_v3_codebase.JoinCodebase(logger, codebaseRepo, codebaseUserRepo, eventSender))                                                                        // Used by the web (2021-10-04)
	auth.POST("/v3/views", routes_v3_view.Create(logger, viewRepo, codebaseUserRepo, analyticsClient, workspaceReader, gitSnapshotter, snapshotRepo, workspaceWriter, executorProvider, eventSender)) // Used by the command line client
	authedViews := auth.Group("/v3/views/:viewID", view_auth.ValidateViewAccessMiddleware(authService, viewRepo))
	authedViews.GET("", routes_v3_view.Get(viewRepo, workspaceReader, userRepo, logger))                                                           // Used by the command line client
	authedViews.POST("/ignore-file", routes_v3_change.IgnoreFile(logger, viewRepo, codebaseUserRepo, executorProvider, viewUpdatedFunc))           // Used by the web (2021-10-04)
	authedViews.GET("/ignores", routes_v3_view.Ignores(logger, executorProvider, viewRepo))                                                        // Called from client-side sturdy-cli
	auth.GET("/v3/stream", routes.Stream(logger, viewRepo, codebaseViewEvents, workspaceReader, authService, workspaceService, suggestionService)) // Used by the web (2021-10-04)
	rebase := auth.Group("/v3/rebase/")
	rebase.Use(view_auth.ValidateViewAccessMiddleware(authService, viewRepo))
	rebase.GET(":viewID", routes_v3_sync.Status(viewRepo, executorProvider, logger))                                    // Used by the web (2021-10-04)
	rebase.POST(":viewID/start", routes_v3_sync.StartV2(logger, syncService))                                           // Used by the web (2021-10-25)
	rebase.POST(":viewID/resolve", routes_v3_sync.ResolveV2(logger, syncService))                                       // Used by the web (2021-10-25)
	auth.POST("/v3/changes/:id/update", routes_v3_change.Update(logger, codebaseUserRepo, analyticsClient, changeRepo)) // Used by the web (2021-10-04)
	auth.POST("/v3/workspaces", routes_v3_workspace.Create(logger, workspaceService, codebaseUserRepo))                 // Used by the command line client
	// Used by LBS to check for health
	publ.GET("/readyz", func(c *gin.Context) { c.Status(http.StatusOK) })
	publ.POST("/v3/waitinglist", waitinglist.Insert(logger, analyticsClient, waitingListRepo))                                                                                                                // Used by the web (2021-10-04)
	publ.POST("/v3/acl-request-enterprise", acl.Insert(logger, analyticsClient, aclInterestRepo))                                                                                                             // Used by the web (2021-10-04)
	publ.POST("/v3/instant-integration", instantintegration.Insert(logger, analyticsClient, instantIntegrationInterestRepo))                                                                                  // Used by the web (2021-10-27)
	auth.POST("/v3/pki/add-public-key", routes_v3_pki.AddPublicKey(userPublicKeyRepo))                                                                                                                        // Used by the command line client
	publ.POST("/v3/pki/verify", routes_v3_pki.Verify(userPublicKeyRepo))                                                                                                                                      // Used by the command line client
	publ.POST("/v3/mutagen/validate-view", routes_v3_mutagen.ValidateView(logger, viewRepo, analyticsClient, eventSender))                                                                                    // Called from server-side mutagen
	publ.POST("/v3/mutagen/sync-transitions", routes_v3_mutagen.SyncTransitions(logger, snapshotterQueue, viewRepo, gcQueue, presenceService, snapshotRepo, workspaceReader, suggestionService, eventSender)) // Called from server-side mutagen
	publ.GET("/v3/mutagen/views/:id/allows", routes_v3_mutagen.ListAllows(logger, viewRepo, authService))                                                                                                     // Called form server-side mutagen
	publ.POST("/v3/mutagen/update-status", routes_v3_mutagen.UpdateStatus(logger, viewStatusRepo, viewRepo, eventSender))                                                                                     // Called from client-side mutagen
	auth.GET("/v3/mutagen/get-view/:id", routes_v3_mutagen.GetView(logger, viewRepo, codebaseUserRepo, codebaseRepo))                                                                                         // Called from client-side sturdy-cli
	publ.POST("/v3/unsubscribe", routes_v3_newsletter.Unsubscribe(logger, userRepo, notificationSettingsRepo))
	return (*Engine)(r)
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

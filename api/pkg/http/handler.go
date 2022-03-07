package http

import (
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	service_analytics "getsturdy.com/api/pkg/analytics/service"
	authz "getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	routes_blobs "getsturdy.com/api/pkg/blobs/routes"
	service_blobs "getsturdy.com/api/pkg/blobs/service"
	db_change "getsturdy.com/api/pkg/changes/db"
	routes_v3_change "getsturdy.com/api/pkg/changes/routes"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	routes_v3_codebase "getsturdy.com/api/pkg/codebase/routes"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/configuration/flags"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	worker_gc "getsturdy.com/api/pkg/gc/worker"
	"getsturdy.com/api/pkg/ginzap"
	sturdygrapql "getsturdy.com/api/pkg/graphql"
	"getsturdy.com/api/pkg/ip"
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
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	routes_v3_sync "getsturdy.com/api/pkg/sync/routes"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/users/avatars/uploader"
	db_user "getsturdy.com/api/pkg/users/db"
	routes_v3_user "getsturdy.com/api/pkg/users/routes"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/pkg/version"
	view_auth "getsturdy.com/api/pkg/view/auth"
	db_view "getsturdy.com/api/pkg/view/db"
	meta_view "getsturdy.com/api/pkg/view/meta"
	routes_v3_view "getsturdy.com/api/pkg/view/routes"
	"getsturdy.com/api/pkg/waitinglist"
	"getsturdy.com/api/pkg/waitinglist/acl"
	"getsturdy.com/api/pkg/waitinglist/instantintegration"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	routes_v3_workspace "getsturdy.com/api/pkg/workspaces/routes"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"

	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Configuration struct {
	Addr             flags.Addr `long:"addr" description:"Address to listen on" default:"localhost:3000"`
	AllowCORSOrigins []string   `long:"allow-cors-origin" description:"Additional origin that is allowed to make CORS requests (can be provided multiple times)"`
}

func ginMode() string {
	if version.IsDevelopment() {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}

type Engine gin.Engine

func ProvideHandler(
	logger *zap.Logger,
	config *Configuration,
	userRepo db_user.Repository,
	analyticsService *service_analytics.Service,
	waitingListRepo waitinglist.WaitingListRepo,
	aclInterestRepo acl.ACLInterestRepo,
	instantIntegrationInterestRepo instantintegration.InstantIntegrationInterestRepo,
	codebaseRepo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	viewRepo db_view.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	userPublicKeyRepo db_pki.Repo,
	snapshotterQueue worker_snapshots.Queue,
	snapshotRepo db_snapshots.Repository,
	changeRepo db_change.Repository,
	gcQueue *worker_gc.Queue,
	gitSnapshotter snapshotter.Snapshotter,
	workspaceWriter db_workspaces.WorkspaceWriter,
	viewUpdatedFunc meta_view.ViewUpdatedFunc,
	executorProvider executor.Provider,
	viewStatusRepo db_mutagen.ViewStatusRepository,
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository,
	eventSender events.EventSender,
	eventSenderV2 *eventsv2.Publisher,
	presenceService service_presence.Service,
	suggestionService *service_suggestion.Service,
	workspaceService service_workspace.Service,
	userService service_user.Service,
	syncService *service_sync.Service,
	jwtService *service_jwt.Service,
	codebaseService *service_codebase.Service,
	authService *service_auth.Service,
	grapqhlResolver *sturdygrapql.RootResolver,
	blobsService *service_blobs.Service,
	uploader uploader.Uploader,
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

	allowOrigins = append(allowOrigins, config.AllowCORSOrigins...)
	cors := ginCors.New(ginCors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"POST, OPTIONS, GET, PUT, DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, Cookie", "sentry-trace"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowWebSockets:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	gin.SetMode(ginMode())
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(accessLogger(logger, time.RFC3339, true))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(cors)
	r.Use(setIp)

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
	publ.POST("/v3/auth", routes_v3_user.Login(logger, userRepo, analyticsService, jwtService))
	publ.POST("/v3/users", routes_v3_user.Signup(logger, userService, jwtService, analyticsService))
	publ.POST("/v3/auth/destroy", routes_v3_user.AuthDestroy)
	auth.POST("/v3/auth/client-token", routes_v3_user.ClientToken(userRepo, jwtService))
	auth.POST("/v3/auth/renew-token", routes_v3_user.RenewToken(logger, userRepo, jwtService))
	auth.POST("/v3/user/update-avatar", routes_v3_user.UpdateAvatar(logger, userRepo, uploader))                                                                                                         // Used by the web (2021-10-04)
	auth.GET("/v3/user", routes_v3_user.GetSelf(userRepo, jwtService))                                                                                                                                   // Used by the command line client
	auth.POST("/v3/codebases", routes_v3_codebase.Create(logger, codebaseService))                                                                                                                       // Used by the web (2021-10-04)
	auth.GET("/v3/codebases/:id", routes_v3_codebase.Get(codebaseRepo, codebaseUserRepo, logger, userService))                                                                                           // Used by the command line client
	auth.POST("/v3/codebases/:id/invite", routes_v3_codebase.Invite(codebaseService, authService))                                                                                                       // No longer used (after 2022-01-31)
	publ.GET("/v3/join/get-codebase/:code", routes_v3_codebase.JoinGetCodebase(logger, codebaseRepo))                                                                                                    // Used by the web (2021-10-04)
	auth.POST("/v3/join/codebase/:code", routes_v3_codebase.JoinCodebase(logger, codebaseRepo, codebaseUserRepo, eventSender))                                                                           // Used by the web (2021-10-04)
	auth.POST("/v3/views", routes_v3_view.Create(logger, viewRepo, codebaseUserRepo, analyticsService, workspaceReader, gitSnapshotter, snapshotRepo, workspaceWriter, executorProvider, eventSenderV2)) // Used by the command line client
	authedViews := auth.Group("/v3/views/:viewID", view_auth.ValidateViewAccessMiddleware(authService, viewRepo))
	authedViews.GET("", routes_v3_view.Get(viewRepo, workspaceReader, logger, userService))                                              // Used by the command line client
	authedViews.POST("/ignore-file", routes_v3_change.IgnoreFile(logger, viewRepo, codebaseUserRepo, executorProvider, viewUpdatedFunc)) // Used by the web (2021-10-04)
	authedViews.GET("/ignores", routes_v3_view.Ignores(logger, executorProvider, viewRepo))                                              // Called from client-side sturdy-cli
	rebase := auth.Group("/v3/rebase/")
	rebase.Use(view_auth.ValidateViewAccessMiddleware(authService, viewRepo))
	rebase.GET(":viewID", routes_v3_sync.Status(viewRepo, executorProvider, logger))                                     // Used by the web (2021-10-04)
	rebase.POST(":viewID/start", routes_v3_sync.StartV2(logger, syncService, workspaceService))                          // Used by the web (2021-10-25)
	rebase.POST(":viewID/resolve", routes_v3_sync.ResolveV2(logger, syncService))                                        // Used by the web (2021-10-25)
	auth.POST("/v3/changes/:id/update", routes_v3_change.Update(logger, codebaseUserRepo, analyticsService, changeRepo)) // Used by the web (2021-10-04)
	auth.POST("/v3/workspaces", routes_v3_workspace.Create(logger, workspaceService, codebaseUserRepo))                  // Used by the command line client
	// Used by LBS to check for health
	publ.GET("/readyz", func(c *gin.Context) { c.Status(http.StatusOK) })
	publ.POST("/v3/waitinglist", waitinglist.Insert(logger, analyticsService, waitingListRepo))                                                                                  // Used by the web (2021-10-04)
	publ.POST("/v3/acl-request-enterprise", acl.Insert(logger, analyticsService, aclInterestRepo))                                                                               // Used by the web (2021-10-04)
	publ.POST("/v3/instant-integration", instantintegration.Insert(logger, analyticsService, instantIntegrationInterestRepo))                                                    // Used by the web (2021-10-27)
	auth.POST("/v3/pki/add-public-key", routes_v3_pki.AddPublicKey(userPublicKeyRepo))                                                                                           // Used by the command line client
	publ.POST("/v3/pki/verify", routes_v3_pki.Verify(userPublicKeyRepo))                                                                                                         // Used by the command line client
	publ.POST("/v3/mutagen/validate-view", routes_v3_mutagen.ValidateView(logger, viewRepo, analyticsService, eventSenderV2))                                                    // Called from server-side mutagen
	publ.POST("/v3/mutagen/sync-transitions", routes_v3_mutagen.SyncTransitions(logger, snapshotterQueue, viewRepo, gcQueue, presenceService, suggestionService, eventSenderV2)) // Called from server-side mutagen
	publ.GET("/v3/mutagen/views/:id/allows", routes_v3_mutagen.ListAllows(logger, viewRepo, authService))                                                                        // Called form server-side mutagen
	publ.POST("/v3/mutagen/update-status", routes_v3_mutagen.UpdateStatus(logger, viewStatusRepo, viewRepo, eventSenderV2))                                                      // Called from client-side mutagen
	auth.GET("/v3/mutagen/get-view/:id", routes_v3_mutagen.GetView(logger, viewRepo, codebaseUserRepo, codebaseRepo))                                                            // Called from client-side sturdy-cli
	publ.POST("/v3/unsubscribe", routes_v3_newsletter.Unsubscribe(logger, userRepo, notificationSettingsRepo))

	routes_blobs.Register(publ.Group("/v3/blobs"), logger, blobsService)
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

func setIp(c *gin.Context) {
	// This is not checking if the remote IP is "trusted" or not
	// TODO: Allow configuration for trusted proxies?
	remoteIp, _ := c.RemoteIP()
	c.Request = c.Request.WithContext(ip.NewContext(c.Request.Context(), remoteIp))
	c.Next()
}

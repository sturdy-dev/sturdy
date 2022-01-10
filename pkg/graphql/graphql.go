package graphql

import (
	"context"
	_ "embed"
	goerrors "errors"
	"net/http"
	"time"

	graphql_buildkite "mash/pkg/integrations/buildkite/graphql"
	service_buildkite "mash/pkg/integrations/buildkite/service"
	graphql_ci "mash/pkg/integrations/graphql"
	service_organization "mash/pkg/organization/service"

	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	graphql_author "mash/pkg/author/graphql"
	db_change "mash/pkg/change/db"
	graphql_change "mash/pkg/change/graphql"
	service_change "mash/pkg/change/service"
	service_ci "mash/pkg/ci/service"
	graphql_acl "mash/pkg/codebase/acl/graphql"
	provider_acl "mash/pkg/codebase/acl/provider"
	db_codebase "mash/pkg/codebase/db"
	graphql_codebase "mash/pkg/codebase/graphql"
	service_codebase "mash/pkg/codebase/service"
	db_comments "mash/pkg/comments/db"
	graphql_comments "mash/pkg/comments/graphql"
	"mash/pkg/ctxlog"
	graphql_features "mash/pkg/features/graphql"
	graphql_file "mash/pkg/file/graphql"
	github_client "mash/pkg/github/client"
	"mash/pkg/github/config"
	db_github "mash/pkg/github/db"
	graphql_github "mash/pkg/github/graphql"
	graphql_pr "mash/pkg/github/graphql/pr"
	service_github "mash/pkg/github/service"
	"mash/pkg/graphql/dataloader"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/graphql/schema"
	service_jwt "mash/pkg/jwt/service"
	db_mutagen "mash/pkg/mutagen/db"
	db_newsletter "mash/pkg/newsletter/db"
	db_notification "mash/pkg/notification/db"
	graphql_notification "mash/pkg/notification/graphql"
	notification_sender "mash/pkg/notification/sender"
	servcie_notification "mash/pkg/notification/service"
	db_onboarding "mash/pkg/onboarding/db"
	graphql_onboarding "mash/pkg/onboarding/graphql"
	graphql_organization "mash/pkg/organization/graphql"
	db_pki "mash/pkg/pki/db"
	graphql_pki "mash/pkg/pki/graphql"
	graphql_presence "mash/pkg/presence/graphql"
	service_presence "mash/pkg/presence/service"
	db_review "mash/pkg/review/db"
	graphql_review "mash/pkg/review/graphql"
	graphql_servicetokens "mash/pkg/servicetokens/graphql"
	service_servicetokens "mash/pkg/servicetokens/service"
	db_snapshots "mash/pkg/snapshots/db"
	"mash/pkg/snapshots/snapshotter"
	graphql_statuses "mash/pkg/statuses/graphql"
	service_statuses "mash/pkg/statuses/service"
	graphql_suggestion "mash/pkg/suggestions/graphql"
	service_suggestion "mash/pkg/suggestions/service"
	db_user "mash/pkg/user/db"
	graphql_user "mash/pkg/user/graphql"
	service_user "mash/pkg/user/service"
	db_view "mash/pkg/view/db"
	"mash/pkg/view/events"
	graphql_view "mash/pkg/view/graphql"
	"mash/pkg/view/view_workspace_snapshot"
	db_activity "mash/pkg/workspace/activity/db"
	graphql_workspace_activity "mash/pkg/workspace/activity/graphql"
	activity_sender "mash/pkg/workspace/activity/sender"
	service_activity "mash/pkg/workspace/activity/service"
	db_workspace "mash/pkg/workspace/db"
	graphql_workspace "mash/pkg/workspace/graphql"
	service_workspace "mash/pkg/workspace/service"
	graphql_workspace_watchers "mash/pkg/workspace/watchers/graphql"
	service_workspace_watchers "mash/pkg/workspace/watchers/service"
	"mash/vcs/executor"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-go/trace"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/posthog/posthog-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type RootResolver struct {
	resolvers.ACLRootResolver
	resolvers.AuthorRootResolver
	resolvers.BuildkiteInstantIntegrationRootResolver
	resolvers.ChangeRootResolver
	resolvers.CodebaseGitHubIntegrationRootResolver
	resolvers.CodebaseRootResolver
	resolvers.CommentRootResolver
	resolvers.FeaturesRootResolver
	resolvers.GitHubAppRootResolver
	resolvers.GitHubPullRequestRootResolver
	resolvers.IntegrationRootResolver
	resolvers.NotificationRootResolver
	resolvers.OnboardingRootResolver
	resolvers.OrganizationRootResolver
	resolvers.PKIRootResolver
	resolvers.PresenceRootResolver
	resolvers.ReviewRootResolver
	resolvers.ServiceTokensRootResolver
	resolvers.StatusesRootResolver
	resolvers.SuggestionRootResolver
	resolvers.UserRootResolver
	resolvers.ViewRootResolver
	resolvers.WorkspaceActivityRootResolver
	resolvers.WorkspaceRootResolver
	resolvers.WorkspaceWatcherRootResolver

	schema     *graphql.Schema
	jwtService *service_jwt.Service
	logger     *zap.Logger
}

func New(
	logger *zap.Logger,
	codebaseRepo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	workspaceReader db_workspace.WorkspaceReader,
	userRepo db_user.Repository,
	viewRepo db_view.Repository,
	snapshotter snapshotter.Snapshotter,
	viewWorkspaceSnapshotsRepo view_workspace_snapshot.Repository,
	snapshotRepo db_snapshots.Repository,
	commentsRepo db_comments.Repository,
	viewEvents events.EventReadWriter,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,
	notificationRepository db_notification.Repository,
	notificationSender notification_sender.NotificationSender,
	gitHubUserRepo db_github.GitHubUserRepo,
	gitHubPRRepo db_github.GitHubPRRepo,
	gitHubInstallationRepo db_github.GitHubInstallationRepo,
	gitHubAppConfig config.GitHubAppConfig,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubClientProvider github_client.ClientProvider,
	gitHubPersonalClientProvider github_client.PersonalClientProvider,
	postHogClient posthog.Client,
	workspaceWriter db_workspace.WorkspaceWriter,
	executorProvider executor.Provider,
	viewStatusRepo db_mutagen.ViewStatusRepository,
	notificationSettingsRepo db_newsletter.NotificationSettingsRepository,
	workspaceActivityRepo db_activity.ActivityRepository,
	aclProvider *provider_acl.Provider,
	reviewRepo db_review.ReviewRepository,
	workspaceActivityReadsRepo db_activity.ActivityReadsRepository,
	activitySender activity_sender.ActivitySender,
	eventSender events.EventSender,
	gitSnapshottet snapshotter.Snapshotter,
	presenceService service_presence.Service,
	suggestionService *service_suggestion.Service,
	workspaceService service_workspace.Service,
	notificationPreferencesServcie *servcie_notification.Preferences,
	changeService *service_change.Service,
	userService *service_user.Service,
	ciService *service_ci.Service,
	statusesService *service_statuses.Service,
	completedOnboardingStepsRepo db_onboarding.CompletedOnboardingStepsRepository,
	workspaceWatchersService *service_workspace_watchers.Service,
	jwtService *service_jwt.Service,
	activityService *service_activity.Service,
	pkiRepo db_pki.Repo,
	gitHubService *service_github.Service,
	codebaseService *service_codebase.Service,
	serviceTokensService *service_servicetokens.Service,
	buidkiteService *service_buildkite.Service,
	authService *service_auth.Service,
	organizationService *service_organization.Service,
) *RootResolver {
	// This pointer dance (pointer to interfaces) allows the resolvers to use each other, without cycles
	var aclResovler = new(resolvers.ACLRootResolver)
	var authorResolver = new(resolvers.AuthorRootResolver)
	var buildkiteRootResolver = new(resolvers.BuildkiteInstantIntegrationRootResolver)
	var changeResolver = new(resolvers.ChangeRootResolver)
	var codebaseGitHubIntegrationResolver = new(resolvers.CodebaseGitHubIntegrationRootResolver)
	var codebaseResolver = new(resolvers.CodebaseRootResolver)
	var commentsResolver = new(resolvers.CommentRootResolver)
	var fileDiffRootResolver = new(resolvers.FileDiffRootResolver)
	var fileResolver = new(resolvers.FileRootResolver)
	var instantIntegrationRootResolver = new(resolvers.IntegrationRootResolver)
	var notificationResolver = new(resolvers.NotificationRootResolver)
	var onboardingRootResolver = new(resolvers.OnboardingRootResolver)
	var organizationRootResolver = new(resolvers.OrganizationRootResolver)
	var pkiRootResolver = new(resolvers.PKIRootResolver)
	var prResolver = new(resolvers.GitHubPullRequestRootResolver)
	var presenceRootResolver = new(resolvers.PresenceRootResolver)
	var reviewResolver = new(resolvers.ReviewRootResolver)
	var serviceTokensRootResolver = new(resolvers.ServiceTokensRootResolver)
	var statusRootResolver = new(resolvers.StatusesRootResolver)
	var suggestionResolver = new(resolvers.SuggestionRootResolver)
	var userResolver = new(resolvers.UserRootResolver)
	var viewResolver = new(resolvers.ViewRootResolver)
	var workspaceActivityResolver = new(resolvers.WorkspaceActivityRootResolver)
	var workspaceResolver = new(resolvers.WorkspaceRootResolver)
	var workspaceWatcherRootResolver = new(resolvers.WorkspaceWatcherRootResolver)

	*aclResovler = graphql_acl.NewResolver(aclProvider, userRepo)
	*authorResolver = graphql_author.NewResolver(userRepo, logger)
	*changeResolver = graphql_change.NewResolver(
		changeService,
		changeRepo,
		changeCommitRepo,
		commentsRepo,
		authService,
		commentsResolver,
		authorResolver,
		statusRootResolver,
		executorProvider,
		logger,
	)
	*codebaseResolver = graphql_codebase.NewResolver(
		codebaseRepo,
		codebaseUserRepo,
		viewRepo,
		workspaceReader,
		userRepo,
		changeRepo,
		changeCommitRepo,

		workspaceResolver,
		authorResolver,
		viewResolver,
		codebaseGitHubIntegrationResolver,
		aclResovler,
		changeResolver,
		fileResolver,
		instantIntegrationRootResolver,

		logger,
		viewEvents,
		eventSender,
		postHogClient,
		executorProvider,

		authService,
	)
	*commentsResolver = graphql_comments.NewResolver(
		userRepo,
		commentsRepo,
		snapshotRepo,
		workspaceReader,
		viewRepo,
		codebaseUserRepo,
		changeRepo,
		workspaceWatchersService,
		authService,

		eventSender,
		viewEvents,
		notificationSender,
		activitySender,

		authorResolver,
		workspaceResolver,
		changeResolver,

		logger,
		postHogClient,
		executorProvider,
	)
	*codebaseGitHubIntegrationResolver = graphql_github.NewResolver(
		gitHubRepositoryRepo,
		gitHubInstallationRepo,
		executorProvider,
		logger,
		gitHubAppConfig,
		gitHubClientProvider,
		workspaceReader,
		workspaceWriter,
		snapshotter,
		snapshotRepo,
		authService,
		codebaseService,
		workspaceResolver,
		codebaseResolver,
		gitHubService,
	)
	*fileDiffRootResolver = graphql_change.NewFileDiffRootResolver()
	*suggestionResolver = graphql_suggestion.New(
		logger,
		authService,
		suggestionService,
		workspaceService,
		authorResolver,
		fileDiffRootResolver,
		workspaceResolver,
		viewEvents,
	)
	*notificationResolver = graphql_notification.NewResolver(
		notificationRepository,
		codebaseUserRepo,
		codebaseRepo,
		notificationPreferencesServcie,
		authService,
		commentsResolver,
		codebaseResolver,
		authorResolver,
		workspaceResolver,
		reviewResolver,
		suggestionResolver,
		codebaseGitHubIntegrationResolver,
		viewEvents,
		eventSender,
		logger,
	)
	*userResolver = graphql_user.NewResolver(
		userRepo,
		gitHubUserRepo,
		notificationSettingsRepo,
		userService,
		viewResolver,
		notificationResolver,
		logger,
	)

	viewStatusResolver := graphql_view.NewViewStatusRootResolver(viewStatusRepo, logger)

	*viewResolver = graphql_view.NewResolver(
		viewRepo,
		workspaceReader,
		snapshotter,
		viewWorkspaceSnapshotsRepo,
		snapshotRepo,
		authorResolver,
		workspaceResolver,
		workspaceWriter,
		viewEvents,
		eventSender,
		executorProvider,
		logger,
		viewStatusResolver,
		workspaceWatchersService,
		postHogClient,
		codebaseResolver,
		authService,
	)
	*prResolver = graphql_pr.NewResolver(
		logger,
		codebaseResolver,
		workspaceResolver,
		statusRootResolver,
		userRepo,
		codebaseRepo,
		workspaceReader,
		viewRepo,
		gitHubAppConfig,
		gitHubUserRepo,
		gitHubPRRepo,
		gitHubInstallationRepo,
		gitHubRepositoryRepo,
		gitHubClientProvider,
		gitHubPersonalClientProvider,
		viewEvents,
		postHogClient,
		authService,
		gitHubService,
	)

	*statusRootResolver = graphql_statuses.New(
		logger,
		statusesService,
		changeService,
		workspaceService,
		authService,
		gitHubPRRepo,
		changeResolver,
		prResolver,
		viewEvents,
	)

	*workspaceWatcherRootResolver = graphql_workspace_watchers.NewRootResolver(
		logger,

		workspaceWatchersService,
		workspaceService,

		authService,
		viewEvents,

		userResolver,
		workspaceResolver,
	)

	*workspaceResolver = graphql_workspace.NewResolver(
		workspaceReader,
		codebaseRepo,
		viewRepo,
		commentsRepo,
		snapshotRepo,

		codebaseResolver,
		authorResolver,
		viewResolver,
		commentsResolver,
		prResolver,
		changeResolver,
		workspaceActivityResolver,
		reviewResolver,
		presenceRootResolver,
		suggestionResolver,
		statusRootResolver,
		workspaceWatcherRootResolver,

		suggestionService,
		workspaceService,
		authService,

		logger,
		viewEvents,
		workspaceWriter,
		executorProvider,
		eventSender,
		gitSnapshottet,
	)
	*workspaceActivityResolver = graphql_workspace_activity.New(
		workspaceActivityRepo,
		workspaceActivityReadsRepo,
		authorResolver,
		commentsResolver,
		changeResolver,
		reviewResolver,
		workspaceResolver,
		activityService,
		authService,
		eventSender,
		viewEvents,
		logger,
	)
	*reviewResolver = graphql_review.New(
		logger,
		reviewRepo,
		workspaceReader,
		authService,

		authorResolver,
		workspaceResolver,

		eventSender,
		viewEvents,
		notificationSender,
		activitySender,
		workspaceWatchersService,
	)
	*fileResolver = graphql_file.NewFileRootResolver(
		executorProvider,
	)

	*presenceRootResolver = graphql_presence.NewRootResolver(
		presenceService,
		authorResolver,
		workspaceResolver,
		logger,
		viewEvents,
	)

	*instantIntegrationRootResolver = graphql_ci.NewRootResolver(
		ciService,
		changeRepo,
		authService,
		buildkiteRootResolver,
		statusRootResolver,
	)

	*onboardingRootResolver = graphql_onboarding.NewRootResolver(completedOnboardingStepsRepo, eventSender, viewEvents)

	*pkiRootResolver = graphql_pki.NewResolver(pkiRepo, userResolver)

	*serviceTokensRootResolver = graphql_servicetokens.New(authService, serviceTokensService, codebaseService)

	*buildkiteRootResolver = graphql_buildkite.New(authService, buidkiteService, ciService, instantIntegrationRootResolver)

	*organizationRootResolver = graphql_organization.New(organizationService, authorResolver)

	r := &RootResolver{
		jwtService: jwtService,
		logger:     logger,

		ACLRootResolver:                         *aclResovler,
		AuthorRootResolver:                      *authorResolver,
		BuildkiteInstantIntegrationRootResolver: *buildkiteRootResolver,
		ChangeRootResolver:                      *changeResolver,
		CodebaseGitHubIntegrationRootResolver:   *codebaseGitHubIntegrationResolver,
		CodebaseRootResolver:                    *codebaseResolver,
		CommentRootResolver:                     *commentsResolver,
		FeaturesRootResolver:                    graphql_features.Resolver,
		GitHubAppRootResolver:                   graphql_github.NewGitHubAppRootResolver(gitHubAppConfig),
		GitHubPullRequestRootResolver:           *prResolver,
		IntegrationRootResolver:                 *instantIntegrationRootResolver,
		NotificationRootResolver:                *notificationResolver,
		OnboardingRootResolver:                  *onboardingRootResolver,
		OrganizationRootResolver:                *organizationRootResolver,
		PKIRootResolver:                         *pkiRootResolver,
		PresenceRootResolver:                    *presenceRootResolver,
		ReviewRootResolver:                      *reviewResolver,
		ServiceTokensRootResolver:               *serviceTokensRootResolver,
		StatusesRootResolver:                    *statusRootResolver,
		SuggestionRootResolver:                  *suggestionResolver,
		UserRootResolver:                        *userResolver,
		ViewRootResolver:                        *viewResolver,
		WorkspaceActivityRootResolver:           *workspaceActivityResolver,
		WorkspaceRootResolver:                   *workspaceResolver,
		WorkspaceWatcherRootResolver:            *workspaceWatcherRootResolver,
	}

	logger = logger.Named("graphql")
	tracer := &metricTracer{logger: logger}
	r.schema = parseSchema(r, tracer, logger)

	return r
}

func (r *RootResolver) HttpHandler() gin.HandlerFunc {
	h := &relay.Handler{Schema: r.schema}
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if subject, ok := auth.SubjectFromGinContext(c); ok {
			ctx = auth.NewContext(ctx, subject)
		}

		ctx = dataloader.NewContext(ctx)

		h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}

func (r *RootResolver) UnauthenticatedHttpHandler(logger *zap.Logger) gin.HandlerFunc {
	// Don't run against the real schema
	fakeSchema := parseSchema(&RootResolver{}, nil, logger)
	h := &relay.Handler{Schema: fakeSchema}
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

type websocketContextBuilder struct {
	jwtService *service_jwt.Service
}

func (c *websocketContextBuilder) BuildContext(ctx context.Context, r *http.Request) (context.Context, error) {
	subject, err := auth.SubjectFromRequest(r, c.jwtService)
	if err != nil {
		return nil, err
	}

	ctx = auth.NewContext(ctx, subject)
	ctx = dataloader.NewContext(ctx)

	return ctx, nil
}

func (r *RootResolver) WebsocketHandler() gin.HandlerFunc {
	h := graphqlws.NewHandlerFunc(r.schema, &relay.Handler{
		Schema: r.schema,
	}, graphqlws.WithContextGenerator(&websocketContextBuilder{
		jwtService: r.jwtService,
	}))

	return func(c *gin.Context) {
		defer func() {
			if val := recover(); val != nil {
				r.logger.Error("panic in websocket handler", zap.Any("panic_value", val))
				return
			}
		}()

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func CorsMiddleware(allowOrigins []string) func(*gin.Context) {
	allowedOrigins := map[string]struct{}{}
	for _, o := range allowOrigins {
		allowedOrigins[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, sentry-trace")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Next()
	}
}

// Reads and parses the schema from file.
// Associates root resolver. Panics if can't read.
func parseSchema(resolver interface{}, tracer *metricTracer, logger *zap.Logger) *graphql.Schema {
	parsedSchema, err := graphql.ParseSchema(
		schema.String,
		resolver,
		graphql.MaxDepth(20),
		graphql.Tracer(tracer),
		graphql.MaxParallelism(10), // Maximum number of resolvers per request allowed to run in parallel.
		graphql.PanicHandler(NewPanicHandler(logger)),
	)
	if err != nil {
		panic(err)
	}
	return parsedSchema
}

type sturdyPanicHandler struct {
	logger *zap.Logger
}

func NewPanicHandler(logger *zap.Logger) *sturdyPanicHandler {
	return &sturdyPanicHandler{
		logger: logger,
	}
}

func (s *sturdyPanicHandler) MakePanicError(ctx context.Context, value interface{}) *errors.QueryError {
	s.logger.Error("panic in graphql resolver", zap.Any("panic_value", value))
	return errors.Errorf("internal server error")
}

var (
	graphqlFieldsHistogramCounter = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "sturdy_graphql_query_fields_duration_millis_total",
		Help:    "Duration in milliseconds",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 50, 100, 200, 500, 2000, 5000},
	}, []string{"typeName", "fieldName", "hasError"})
)

// implements trace.Tracer
type metricTracer struct {
	logger *zap.Logger
}

func (m *metricTracer) TraceQuery(ctx context.Context, queryString string, operationName string, variables map[string]interface{}, varTypes map[string]*introspection.Type) (context.Context, trace.TraceQueryFinishFunc) {
	fn := func(errors []*errors.QueryError) {
		if len(errors) == 0 {
			return
		}

		var highestLogLevel logLevel = logLevelNone
		for _, err := range errors {
			if l := logLevelForError(err); l > highestLogLevel {
				highestLogLevel = l
			}
		}

		logger := m.logger.With(
			zap.String("queryString", queryString),
			zap.String("operationName", operationName),
			zap.Any("variables", variables),
			zap.Any("varTypes", varTypes),
		)

		if subject, ok := auth.FromContext(ctx); ok {
			logger = logger.With(zap.Stringer("subject", subject.Type))
			if subject.ID != "" {
				logger = logger.With(zap.String("subjectId", subject.ID))
			}
		}

		switch highestLogLevel {
		case logLevelWarn:
			logger.Warn("query failed")
		case logLevelErr:
			logger.Error("query failed")
		}
	}
	return ctx, fn
}

func (m *metricTracer) TraceField(ctx context.Context, label, typeName, fieldName string, trivial bool, args map[string]interface{}) (context.Context, trace.TraceFieldFinishFunc) {
	t0 := time.Now()
	fn := func(err *errors.QueryError) {
		if err != nil {
			l := m.logger.With(
				zap.String("label", label),
				zap.String("typeName", typeName),
				zap.String("fieldName", fieldName),
				zap.Bool("trivial", trivial),
				zap.Any("args", args),
			)

			if subject, ok := auth.FromContext(ctx); ok {
				l = l.With(zap.Stringer("subject", subject.Type))
				if subject.ID != "" {
					l = l.With(zap.String("subjectId", subject.ID))
				}
			}

			msg := "field errors"
			var fields []zap.Field

			if gerr, ok := err.Err.(*gqlerrors.SturdyGraphqlError); ok {
				fields = []zap.Field{
					zap.NamedError("sturdyError", gerr.OriginalError()),
					zap.Any("sturdyErrorExtensions", gerr.Extensions()),
					zap.Error(err),
				}
			} else {
				fields = []zap.Field{zap.Error(err)}
			}

			switch logLevelForError(err) {
			case logLevelWarn:
				l.Warn(msg, fields...)
			case logLevelErr:
				l.Error(msg, fields...)
			}
		}

		hasError := "false"
		if err != nil {
			hasError = "true"
		}
		graphqlFieldsHistogramCounter.WithLabelValues(typeName, fieldName, hasError).Observe(float64(time.Since(t0).Milliseconds()))
	}
	return ctx, fn
}

type logLevel int

const (
	logLevelNone logLevel = 0
	logLevelWarn logLevel = 1
	logLevelErr  logLevel = 2
)

func logLevelForError(err error) logLevel {
	if goerrors.Is(err, gqlerrors.ErrUnauthenticated) {
		return logLevelNone
	} else if ctxlog.IsError(err) && !gqlerrors.IsClientSideError(err) {
		return logLevelErr
	} else {
		return logLevelWarn
	}
}

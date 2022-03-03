package graphql

import (
	"context"
	_ "embed"
	goerrors "errors"
	"fmt"
	"net/http"
	"time"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/ctxlog"
	"getsturdy.com/api/pkg/graphql/dataloader"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/graphql/schema"
	"getsturdy.com/api/pkg/ip"
	service_jwt "getsturdy.com/api/pkg/jwt/service"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-go/trace"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
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
	resolvers.GitHubRootResolver
	resolvers.IntegrationRootResolver
	resolvers.LicenseRootResolver
	resolvers.NotificationRootResolver
	resolvers.OnboardingRootResolver
	resolvers.OrganizationRootResolver
	resolvers.PKIRootResolver
	resolvers.PresenceRootResolver
	resolvers.ReviewRootResolver
	resolvers.InstallationsRootResolver
	resolvers.ServiceTokensRootResolver
	resolvers.StatusesRootResolver
	resolvers.SuggestionRootResolver
	resolvers.UserRootResolver
	resolvers.ViewRootResolver
	resolvers.ActivityRootResolver
	resolvers.WorkspaceRootResolver
	resolvers.WorkspaceWatcherRootResolver

	schema     *graphql.Schema
	jwtService *service_jwt.Service
	logger     *zap.Logger
}

func NewRootResolver(
	logger *zap.Logger,
	jwtService *service_jwt.Service,

	aclResovler resolvers.ACLRootResolver,
	authorResolver resolvers.AuthorRootResolver,
	buildkiteRootResolver resolvers.BuildkiteInstantIntegrationRootResolver,
	changeResolver resolvers.ChangeRootResolver,
	codebaseGitHubIntegrationResolver resolvers.CodebaseGitHubIntegrationRootResolver,
	codebaseResolver resolvers.CodebaseRootResolver,
	commentsResolver resolvers.CommentRootResolver,
	featuresRootResolver resolvers.FeaturesRootResolver,
	gitHubRootResolver resolvers.GitHubRootResolver,
	githubAppResolver resolvers.GitHubAppRootResolver,
	instantIntegrationRootResolver resolvers.IntegrationRootResolver,
	licenseRootResolver resolvers.LicenseRootResolver,
	notificationResolver resolvers.NotificationRootResolver,
	onboardingRootResolver resolvers.OnboardingRootResolver,
	organizationRootResolver resolvers.OrganizationRootResolver,
	pkiRootResolver resolvers.PKIRootResolver,
	prResolver resolvers.GitHubPullRequestRootResolver,
	presenceRootResolver resolvers.PresenceRootResolver,
	reviewResolver resolvers.ReviewRootResolver,
	serverStatusRootResolver resolvers.InstallationsRootResolver,
	serviceTokensRootResolver resolvers.ServiceTokensRootResolver,
	statusRootResolver resolvers.StatusesRootResolver,
	suggestionResolver resolvers.SuggestionRootResolver,
	userResolver resolvers.UserRootResolver,
	viewResolver resolvers.ViewRootResolver,
	activityResolver resolvers.ActivityRootResolver,
	workspaceResolver resolvers.WorkspaceRootResolver,
	workspaceWatcherRootResolver resolvers.WorkspaceWatcherRootResolver,
) *RootResolver {
	r := &RootResolver{
		jwtService: jwtService,
		logger:     logger,

		ACLRootResolver:                         aclResovler,
		AuthorRootResolver:                      authorResolver,
		BuildkiteInstantIntegrationRootResolver: buildkiteRootResolver,
		ChangeRootResolver:                      changeResolver,
		CodebaseGitHubIntegrationRootResolver:   codebaseGitHubIntegrationResolver,
		CodebaseRootResolver:                    codebaseResolver,
		CommentRootResolver:                     commentsResolver,
		FeaturesRootResolver:                    featuresRootResolver,
		GitHubAppRootResolver:                   githubAppResolver,
		GitHubPullRequestRootResolver:           prResolver,
		GitHubRootResolver:                      gitHubRootResolver,
		IntegrationRootResolver:                 instantIntegrationRootResolver,
		LicenseRootResolver:                     licenseRootResolver,
		NotificationRootResolver:                notificationResolver,
		OnboardingRootResolver:                  onboardingRootResolver,
		OrganizationRootResolver:                organizationRootResolver,
		PKIRootResolver:                         pkiRootResolver,
		PresenceRootResolver:                    presenceRootResolver,
		ReviewRootResolver:                      reviewResolver,
		InstallationsRootResolver:               serverStatusRootResolver,
		ServiceTokensRootResolver:               serviceTokensRootResolver,
		StatusesRootResolver:                    statusRootResolver,
		SuggestionRootResolver:                  suggestionResolver,
		UserRootResolver:                        userResolver,
		ViewRootResolver:                        viewResolver,
		ActivityRootResolver:                    activityResolver,
		WorkspaceRootResolver:                   workspaceResolver,
		WorkspaceWatcherRootResolver:            workspaceWatcherRootResolver,
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

		if remoteIP, _ := c.RemoteIP(); remoteIP != nil {
			ctx = ip.NewContext(ctx, remoteIP)
		} else {
			r.logger.Error("could not find and set remoteIP", zap.String("remote_addr", c.Request.RemoteAddr))
		}

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
	return ctx, func(errors []*errors.QueryError) {
		fields := []zap.Field{
			zap.String("queryString", queryString),
			zap.String("operationName", operationName),
			zap.Any("variables", variables),
			zap.Any("varTypes", varTypes),
		}

		highestLogLevel := logLevelNone
		for _, err := range errors {
			if l := logLevelForError(err); l > highestLogLevel {
				highestLogLevel = l
			}

			var gerr *gqlerrors.SturdyGraphqlError
			if goerrors.As(err, &gerr) {
				fields = append(fields, zap.Any(fmt.Sprint(err.Path...), goerrors.Unwrap(gerr.OriginalError())))
				for k, v := range gerr.Extensions() {
					fields = append(fields, zap.Any(k, v))
				}
			} else {
				fields = append(fields, zap.NamedError(fmt.Sprint(err.Path...), err))
			}
		}

		if subject, ok := auth.FromContext(ctx); ok {
			fields = append(fields, zap.Stringer("subject", subject.Type))
			if subject.ID != "" {
				fields = append(fields, zap.String("subjectId", subject.ID))
			}
		}

		switch highestLogLevel {
		case logLevelWarn:
			m.logger.With(fields...).Warn("query failed")
		case logLevelErr:
			m.logger.With(fields...).Error("query failed")
		}
	}
}

func (m *metricTracer) TraceField(ctx context.Context, label, typeName, fieldName string, trivial bool, args map[string]interface{}) (context.Context, trace.TraceFieldFinishFunc) {
	t0 := time.Now()
	return ctx, func(err *errors.QueryError) {
		hasError := "false"
		if err != nil {
			hasError = "true"
		}
		graphqlFieldsHistogramCounter.WithLabelValues(typeName, fieldName, hasError).Observe(float64(time.Since(t0).Milliseconds()))
	}
}

type logLevel uint

const (
	logLevelNone logLevel = iota
	logLevelWarn
	logLevelErr
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

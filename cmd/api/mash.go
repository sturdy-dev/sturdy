package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"mash/pkg/integrations"
	db_buildkite "mash/pkg/integrations/buildkite/db"
	service_buildkite "mash/pkg/integrations/buildkite/service"
	db_integrations "mash/pkg/integrations/db"
	db_organization "mash/pkg/organization/db"

	"mash/db"
	authz "mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	db_change "mash/pkg/change/db"
	routes_v3_change "mash/pkg/change/routes"
	service_change "mash/pkg/change/service"
	db_ci "mash/pkg/ci/db"
	routes_ci "mash/pkg/ci/routes"
	service_ci "mash/pkg/ci/service"
	workers_ci "mash/pkg/ci/workers"
	db_acl "mash/pkg/codebase/acl/db"
	provider_acl "mash/pkg/codebase/acl/provider"
	db_codebase "mash/pkg/codebase/db"
	routes_v3_codebase "mash/pkg/codebase/routes"
	service_codebase "mash/pkg/codebase/service"
	db_comments "mash/pkg/comments/db"
	service_comments "mash/pkg/comments/service"
	"mash/pkg/emails"
	"mash/pkg/emails/transactional"
	db_gc "mash/pkg/gc/db"
	worker_gc "mash/pkg/gc/worker"
	"mash/pkg/ginzap"
	ghappclient "mash/pkg/github/client"
	"mash/pkg/github/config"
	db_github "mash/pkg/github/db"
	routes_v3_ghapp "mash/pkg/github/routes"
	service_github "mash/pkg/github/service"
	workers_github "mash/pkg/github/workers"
	"mash/pkg/gitserver"
	"mash/pkg/gmail"
	sturdygrapql "mash/pkg/graphql"
	db_keys "mash/pkg/jwt/keys/db"
	service_jwt "mash/pkg/jwt/service"
	"mash/pkg/metrics/ginprometheus"
	"mash/pkg/metrics/zapprometheus"
	db_mutagen "mash/pkg/mutagen/db"
	routes_v3_mutagen "mash/pkg/mutagen/routes"
	db_newsletter "mash/pkg/newsletter/db"
	routes_v3_newsletter "mash/pkg/newsletter/routes"
	db_notification "mash/pkg/notification/db"
	notification_sender "mash/pkg/notification/sender"
	service_notification "mash/pkg/notification/service"
	db_onboarding "mash/pkg/onboarding/db"
	db_onetime "mash/pkg/onetime/db"
	service_onetime "mash/pkg/onetime/service"
	service_organization "mash/pkg/organization/service"
	db_pki "mash/pkg/pki/db"
	routes_v3_pki "mash/pkg/pki/routes"
	ph "mash/pkg/posthog"
	db_presence "mash/pkg/presence/db"
	service_presence "mash/pkg/presence/service"
	"mash/pkg/queue"
	db_review "mash/pkg/review/db"
	db_servicetokens "mash/pkg/servicetokens/db"
	service_servicetokens "mash/pkg/servicetokens/service"
	db_snapshots "mash/pkg/snapshots/db"
	"mash/pkg/snapshots/snapshotter"
	worker_snapshots "mash/pkg/snapshots/worker"
	db_statuses "mash/pkg/statuses/db"
	service_statuses "mash/pkg/statuses/service"
	"mash/pkg/stream/routes"
	db_suggestion "mash/pkg/suggestions/db"
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
	ws_meta "mash/pkg/workspace/meta"
	routes_v3_workspace "mash/pkg/workspace/routes"
	service_workspace "mash/pkg/workspace/service"
	db_workspace_watchers "mash/pkg/workspace/watchers/db"
	service_workspace_watchers "mash/pkg/workspace/watchers/service"
	"mash/vcs/executor"
	"mash/vcs/provider"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	reposBasePath := flag.String("repos-base-path", "tmp/repos", "path on the filesystem to where all repos can be found")
	httpListenAddr := flag.String("http-listen-addr", "127.0.0.1:3000", "")
	httpPprofListenAddr := flag.String("http-pprof-listen-addr", "127.0.0.1:6060", "")
	gitListenAddr := flag.String("git-listen-addr", "127.0.0.1:3002", "")
	metricsListenAddr := flag.String("metrics-listen-addr", "127.0.0.1:2112", "")
	dbSourceAddr := flag.String("db", "postgres://mash:mash@127.0.0.1:5432/mash?sslmode=disable", "")
	productionLogger := flag.Bool("production-logger", false, "")
	sendPostHogEvents := flag.Bool("send-posthog-events", false, "")
	enableSendInvitesWorker := flag.Bool("send-invites-worker", false, "")
	gmailTokenJsonPath := flag.String("gmail-token-json-path", "", "used by the invite email sender")
	gmailCredentialsJsonPath := flag.String("gmail-credentials-json-path", "", "used by the invite email sender")
	gitHubAppID := flag.Int64("github-app-id", 122610, "")
	gitHubAppName := flag.String("github-app-name", "sturdy-gustav-localhost", "")
	gitHubAppClientID := flag.String("github-app-client-id", "", "")
	gitHubAppSecret := flag.String("github-app-secret", "", "")
	gitHubAppPrivateKeyPath := flag.String("github-app-private-key-path", "", "")
	_ = flag.Bool("unauthenticated-graphql-introspection", false, "")
	gitLfsHostname := flag.String("git-lfs-hostname", "localhost:8888", "")
	enableTransactionalEmails := flag.Bool("enable-transactional-emails", false, "")
	exportBucketName := flag.String("export-bucket-name", "", "the S3 bucket to be used for change exports")
	developmentAllowExtraCorsOrigin := flag.String("development-allow-extra-cors-origin", "", "Additional CORS origin to be allowed")
	localQueue := flag.Bool("use-local-queues", false, "If set, local queue will be used instead of SQS")

	publicApiHostname := flag.String("public-api-hostname", "localhost", "api.getsturdy.com in production")
	// publicGitHostname := flag.String("public-git-hostname", "git.getsturdy.com", "")

	defaultHostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}
	hostname := flag.String("hostname", defaultHostname, "")
	queuePrefix := flag.String("queue-prefix", "dev"+defaultHostname, "set to 'production' when running in production")

	flag.Parse()

	ctx := context.Background()

	d, err := db.Setup(*dbSourceAddr, true, "file://db/migrations")
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger
	if *productionLogger {
		logger, _ = zap.NewProduction(zap.Hooks(zapprometheus.Hook))
	} else {
		logger, _ = zap.NewDevelopment(zap.Hooks(zapprometheus.Hook))
	}

	gitHubAppConfig := config.GitHubAppConfig{
		GitHubAppID:             *gitHubAppID,
		GitHubAppName:           *gitHubAppName,
		GitHubAppClientID:       *gitHubAppClientID,
		GitHubAppSecret:         *gitHubAppSecret,
		GitHubAppPrivateKeyPath: *gitHubAppPrivateKeyPath,
	}

	userRepo := db_user.NewRepo(d)
	aclRepo := db_acl.NewACLRepository(d)
	codebaseRepo := db_codebase.NewRepo(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	viewRepo := db_view.NewRepo(d)
	workspaceRepo := db_workspace.NewRepo(d)
	waitingListRepo := waitinglist.NewWaitingListRepo(d)
	aclInterestRepo := acl.NewACLInterestRepo(d)
	instantIntegrationInterestRepo := instantintegration.NewInstantIntegrationInterestRepo(d)
	userPublicKeyRepo := db_pki.NewRepo(d)
	snapshotRepo := db_snapshots.NewRepo(d)
	changeRepo := db_change.NewRepo(d)
	changeCommitRepo := db_change.NewCommitRepository(d)
	commentRepo := db_comments.NewRepo(d)
	suggestionRepo := db_suggestion.New(d)
	gcRepo := db_gc.NewRepository(d)
	viewWorkspaceSnapshotsRepo := view_workspace_snapshot.NewRepo(d)
	gitHubInstallationsRepo := db_github.NewGitHubInstallationRepo(d)
	gitHubRepositoryRepo := db_github.NewGitHubRepositoryRepo(d)
	gitHubUserRepo := db_github.NewGitHubUserRepo(d)
	gitHubPRRepo := db_github.NewGitHubPRRepo(d)
	notificationRepo := db_notification.NewRepository(d)
	viewStatusRepo := db_mutagen.NewRepository(d)
	notificationSettingsRepo := db_newsletter.NewNotificationSettingsRepository(d)
	workspaceActivityRepo := db_activity.NewActivityRepo(d)
	reviewRepo := db_review.NewReviewRepository(d)
	workspaceActivityReadsRepo := db_activity.NewActivityReadsRepo(d)
	notificationPreferencesRepo := db_notification.NewPeferenceRepository(d)
	keysRepo := db_keys.New(d)
	statusesRepo := db_statuses.New(d)
	ciCommitRepo := db_ci.NewCommitRepository(d)
	ciConfigRepo := db_integrations.NewIntegrationDatabase(d)
	workspaceWatchersRepo := db_workspace_watchers.NewDB(d)
	commentsService := service_comments.New(commentRepo)
	organizationRepo := db_organization.New(d)
	organizationMemberRepo := db_organization.NewMember(d)
	organizationService := service_organization.New(organizationRepo, organizationMemberRepo)

	awsSession, err := session.NewSession(
		&aws.Config{
			Region: aws.String("eu-north-1"),
		})
	if err != nil {
		logger.Fatal("failed to go create AWS session", zap.Error(err))
	}

	var postHogClient posthog.Client
	if *sendPostHogEvents {
		postHogClient, err = posthog.NewWithConfig("ZuDRoGX9PgxGAZqY4RF9CCJJLpx14h3szUPzm7XBWSg", posthog.Config{Endpoint: "https://app.posthog.com"})
		if err != nil {
			panic(err)
		}
	} else {
		postHogClient = ph.NewFakeClient()
	}

	codebaseViewEvents := events.NewInMemory()
	executorProvider := executor.NewProvider(logger, provider.New(*reposBasePath, *gitLfsHostname))
	aclProvider := provider_acl.New(aclRepo, codebaseUserRepo, userRepo)
	gitHubClientProvider := ghappclient.NewClient
	gitHubPersonalClientProvider := ghappclient.NewPersonalClient

	var q queue.Queue
	if *localQueue {
		q = queue.NewInMemory(logger)
	} else {
		sqs, err := queue.NewSQS(logger, awsSession, *hostname, *queuePrefix)
		if err != nil {
			logger.Fatal("failed to create sqs queue", zap.Error(err))
		}
		q = sqs
	}

	notificationPreferencesService := service_notification.NewPreferences(notificationPreferencesRepo)
	jwtService := service_jwt.NewService(logger, keysRepo)

	var emailSender emails.Sender
	if *enableTransactionalEmails {
		emailSender = emails.NewSES(awsSession)
	} else {
		emailSender = emails.NewLogs(logger)
	}

	transactionalEmailSender := transactional.New(
		logger,
		emailSender,

		userRepo,
		codebaseUserRepo,
		commentRepo,
		changeRepo,
		codebaseRepo,
		workspaceRepo,
		suggestionRepo,
		reviewRepo,
		notificationSettingsRepo,
		gitHubRepositoryRepo,

		jwtService,

		notificationPreferencesService,
		postHogClient,
	)

	oneTimeTokenDB := db_onetime.New(d)
	oneTimeService := service_onetime.New(oneTimeTokenDB)

	userService := service_user.New(logger, userRepo, jwtService, oneTimeService, transactionalEmailSender, postHogClient)

	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, codebaseViewEvents)
	notificationSender := notification_sender.NewNotificationSender(codebaseUserRepo, notificationRepo, userRepo, eventsSender, transactionalEmailSender)
	activityService := service_activity.New(workspaceActivityReadsRepo, eventsSender)
	activitySender := activity_sender.NewActivitySender(codebaseUserRepo, workspaceActivityRepo, activityService, eventsSender)
	workspaceWriter := ws_meta.NewWriterWithEvents(logger, workspaceRepo, eventsSender)
	changeService := service_change.New(executorProvider, awsSession, aclProvider, userRepo, changeRepo, changeCommitRepo, exportBucketName)

	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotRepo, workspaceRepo, workspaceWriter, viewRepo, eventsSender, executorProvider, logger)

	// pointer dance to solve circular dependency
	var gitHubService = new(service_github.Service)
	githubImporterQueue := workers_github.NewImporterQueue(logger, q, gitHubService)
	githubClonerQueue := workers_github.NewClonerQueue(logger, q, gitHubService)
	*gitHubService = *service_github.New(
		logger,

		gitHubRepositoryRepo,
		gitHubInstallationsRepo,
		gitHubUserRepo,
		gitHubPRRepo,
		gitHubAppConfig,
		gitHubClientProvider,
		gitHubPersonalClientProvider,

		githubImporterQueue,
		githubClonerQueue,

		workspaceWriter,
		workspaceRepo,
		codebaseUserRepo,
		codebaseRepo,

		executorProvider,

		gitSnapshotter,
		postHogClient,
		notificationSender,
		eventsSender,

		userService,
	)

	// Start queues
	snapshotterQueue := worker_snapshots.New(logger, q, gitSnapshotter)

	viewUpdatedFunc := meta_view.NewViewUpdatedFunc(workspaceRepo, workspaceWriter, eventsSender, snapshotterQueue)
	statusesService := service_statuses.New(logger, statusesRepo, eventsSender)
	ciService := service_ci.New(logger, executorProvider, ciConfigRepo, ciCommitRepo, changeRepo, changeCommitRepo, *publicApiHostname, statusesService, jwtService)
	ciBuildQueue := workers_ci.New(logger, q, ciService)

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
		snapshotterQueue,
		gitSnapshotter,
		ciBuildQueue,
	)
	workspaceWatchersService := service_workspace_watchers.New(workspaceWatchersRepo, eventsSender)

	suggestionService := service_suggestion.New(
		logger,
		suggestionRepo,
		workspaceService,
		executorProvider,
		gitSnapshotter,
		postHogClient,
		notificationSender,
		eventsSender,
	)
	gcQueue := worker_gc.New(logger, q, gcRepo, viewRepo, snapshotRepo, workspaceRepo, suggestionService, executorProvider)

	presenceRepo := db_presence.NewRepo(d)
	presenceService := service_presence.New(presenceRepo, eventsSender)

	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo)
	authService := service_auth.New(codebaseService, userService, workspaceService, aclProvider)
	serviceTokensService := service_servicetokens.New(db_servicetokens.NewDatabase(d))

	buildkiteService := service_buildkite.New(db_buildkite.NewDatabase(d))

	// register ci integrations
	integrations.Register(integrations.ProviderTypeBuildkite, buildkiteService)

	syncService := service_sync.New(logger, executorProvider, viewRepo, workspaceRepo, workspaceWriter, gitSnapshotter)

	completedOnboardingStepsRepo := db_onboarding.New(d)

	gitsrv := gitserver.New(logger, serviceTokensService, jwtService, codebaseService, executorProvider)

	wg, ctx := errgroup.WithContext(ctx)

	// github cloner queue
	wg.Go(func() error {
		if err := githubClonerQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github cloner queue: %w", err)
		}
		return nil
	})

	// github importer queue
	wg.Go(func() error {
		if err := githubImporterQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github importer queue: %v", err)
		}
		return nil
	})

	// snapshotter queue
	wg.Go(func() error {
		if err := snapshotterQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start snapshotter queue: %w", err)
		}
		return nil
	})

	// ci build queue
	wg.Go(func() error {
		if err := ciBuildQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start ci build queue: %w", err)
		}
		return nil
	})

	// gc queue
	wg.Go(func() error {
		if err := gcQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start gc queue: %w", err)
		}
		return nil
	})

	// Start the git HTTP server
	wg.Go(func() error {
		if err := gitsrv.Start(ctx, *gitListenAddr); err != nil {
			return fmt.Errorf("failed to start git server: %w", err)
		}
		return nil
	})

	// Send invitation emails
	// TODO: delete?
	if *enableSendInvitesWorker {
		gmailService, err := gmail.GetService(*gmailCredentialsJsonPath, *gmailTokenJsonPath)
		if err != nil {
			logger.Fatal("failed to get gmail credentials", zap.Error(err))
			return
		}
		go waitinglist.Worker(logger, waitingListRepo, gmailService)
	}

	// Pprof server
	wg.Go(func() error {
		if err := http.ListenAndServe(*httpPprofListenAddr, nil); err != http.ErrServerClosed {
			return fmt.Errorf("failed to start http pprof server: %w", err)
		}
		return nil
	})

	// Metrics server
	wg.Go(func() error {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		srv := http.Server{Addr: *metricsListenAddr, Handler: mux}
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		return nil
	})

	// blocking
	wg.Go(func() error {
		if err := webserver(
			logger,
			*httpListenAddr,
			userRepo,
			postHogClient,
			waitingListRepo,
			aclInterestRepo,
			instantIntegrationInterestRepo,
			codebaseRepo,
			codebaseUserRepo,
			viewRepo,
			workspaceRepo,
			userPublicKeyRepo,
			snapshotterQueue,
			snapshotRepo,
			changeRepo,
			changeCommitRepo,
			commentRepo,
			codebaseViewEvents,
			gcQueue,
			gitSnapshotter,
			viewWorkspaceSnapshotsRepo,
			gitHubInstallationsRepo,
			gitHubRepositoryRepo,
			gitHubUserRepo,
			gitHubPRRepo,
			gitHubAppConfig,
			notificationRepo,
			notificationSender,
			gitHubClientProvider,
			gitHubPersonalClientProvider,
			githubClonerQueue,
			workspaceWriter,
			viewUpdatedFunc,
			executorProvider,
			viewStatusRepo,
			notificationSettingsRepo,
			workspaceActivityRepo,
			aclProvider,
			reviewRepo,
			workspaceActivityReadsRepo,
			activitySender,
			eventsSender,
			presenceService,
			suggestionService,
			workspaceService,
			notificationPreferencesService,
			changeService,
			userService,
			ciService,
			statusesService,
			completedOnboardingStepsRepo,
			syncService,
			workspaceWatchersService,
			jwtService,
			activityService,
			commentsService,
			gitHubService,
			codebaseService,
			serviceTokensService,
			buildkiteService,
			authService,
			ciBuildQueue,
			*developmentAllowExtraCorsOrigin,
			organizationService,
		); err != nil {
			return fmt.Errorf("failed to start webserver: %w", err)
		}
		return nil
	})

	if err := wg.Wait(); err != nil {
		logger.Fatal("failed to start", zap.Error(err))
	}
}

func webserver(
	logger *zap.Logger,
	httpListenAddr string,
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
	organizationService *service_organization.Service,
) error {
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
		organizationService,
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

	// This call is blocking
	if err := r.Run(httpListenAddr); err != http.ErrServerClosed {
		return err
	}

	return nil
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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mash/pkg/integrations"
	db_buildkite "mash/pkg/integrations/buildkite/db"
	service_buildkite "mash/pkg/integrations/buildkite/service"
	db_integrations "mash/pkg/integrations/db"
	"net/http"
	_ "net/http/pprof"
	"os"

	"mash/db"
	service_auth "mash/pkg/auth/service"
	db_change "mash/pkg/change/db"
	service_change "mash/pkg/change/service"
	db_ci "mash/pkg/ci/db"
	service_ci "mash/pkg/ci/service"
	workers_ci "mash/pkg/ci/workers"
	db_acl "mash/pkg/codebase/acl/db"
	provider_acl "mash/pkg/codebase/acl/provider"
	db_codebase "mash/pkg/codebase/db"
	service_codebase "mash/pkg/codebase/service"
	db_comments "mash/pkg/comments/db"
	service_comments "mash/pkg/comments/service"
	"mash/pkg/emails"
	"mash/pkg/emails/transactional"
	db_gc "mash/pkg/gc/db"
	worker_gc "mash/pkg/gc/worker"
	ghappclient "mash/pkg/github/client"
	"mash/pkg/github/config"
	db_github "mash/pkg/github/db"
	service_github "mash/pkg/github/service"
	workers_github "mash/pkg/github/workers"
	"mash/pkg/gitserver"
	"mash/pkg/gmail"
	httpx "mash/pkg/http"
	db_keys "mash/pkg/jwt/keys/db"
	service_jwt "mash/pkg/jwt/service"
	"mash/pkg/metrics/zapprometheus"
	db_mutagen "mash/pkg/mutagen/db"
	db_newsletter "mash/pkg/newsletter/db"
	db_notification "mash/pkg/notification/db"
	notification_sender "mash/pkg/notification/sender"
	service_notification "mash/pkg/notification/service"
	db_onboarding "mash/pkg/onboarding/db"
	db_onetime "mash/pkg/onetime/db"
	service_onetime "mash/pkg/onetime/service"
	db_pki "mash/pkg/pki/db"
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
	db_suggestion "mash/pkg/suggestions/db"
	service_suggestion "mash/pkg/suggestions/service"
	service_sync "mash/pkg/sync/service"
	db_user "mash/pkg/user/db"
	service_user "mash/pkg/user/service"
	db_view "mash/pkg/view/db"
	"mash/pkg/view/events"
	meta_view "mash/pkg/view/meta"
	"mash/pkg/view/view_workspace_snapshot"
	"mash/pkg/waitinglist"
	"mash/pkg/waitinglist/acl"
	"mash/pkg/waitinglist/instantintegration"
	db_activity "mash/pkg/workspace/activity/db"
	activity_sender "mash/pkg/workspace/activity/sender"
	service_activity "mash/pkg/workspace/activity/service"
	db_workspace "mash/pkg/workspace/db"
	ws_meta "mash/pkg/workspace/meta"
	service_workspace "mash/pkg/workspace/service"
	db_workspace_watchers "mash/pkg/workspace/watchers/db"
	service_workspace_watchers "mash/pkg/workspace/watchers/service"
	"mash/vcs/executor"
	"mash/vcs/provider"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

	server := httpx.NewServer(
		httpx.NewHandler(
			logger,
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
		),
	)

	wg.Go(func() error {
		if err := server.ListenAndServe(*httpListenAddr); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	if err := wg.Wait(); err != nil {
		logger.Fatal("failed to start", zap.Error(err))
	}
}

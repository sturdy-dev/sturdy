package main

import (
	"context"
	"flag"
	"log"
	_ "net/http/pprof"
	"os"

	"mash/db"
	"mash/pkg/api"
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
	"mash/pkg/di"
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
	"mash/pkg/graphql"
	"mash/pkg/http"
	"mash/pkg/integrations"
	db_buildkite "mash/pkg/integrations/buildkite/db"
	service_buildkite "mash/pkg/integrations/buildkite/service"
	db_integrations "mash/pkg/integrations/db"
	db_keys "mash/pkg/jwt/keys/db"
	service_jwt "mash/pkg/jwt/service"
	client_license "mash/pkg/license/client"
	db_license "mash/pkg/license/db"
	service_license "mash/pkg/license/service"
	validator_license "mash/pkg/license/validator"
	"mash/pkg/metrics/zapprometheus"
	db_mutagen "mash/pkg/mutagen/db"
	db_newsletter "mash/pkg/newsletter/db"
	db_notification "mash/pkg/notification/db"
	notification_sender "mash/pkg/notification/sender"
	service_notification "mash/pkg/notification/service"
	db_onboarding "mash/pkg/onboarding/db"
	db_onetime "mash/pkg/onetime/db"
	service_onetime "mash/pkg/onetime/service"
	db_organization "mash/pkg/organization/db"
	service_organization "mash/pkg/organization/service"
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
	"github.com/jmoiron/sqlx"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
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
	_ = flag.Bool("send-invites-worker", false, "")
	_ = flag.String("gmail-token-json-path", "", "used by the invite email sender")
	_ = flag.String("gmail-credentials-json-path", "", "used by the invite email sender")
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

	var logger *zap.Logger
	if *productionLogger {
		logger, _ = zap.NewProduction(zap.Hooks(zapprometheus.Hook))
	} else {
		logger, _ = zap.NewDevelopment(zap.Hooks(zapprometheus.Hook))
	}

	providers := []interface{}{
		func() *zap.Logger { return logger },
		func() (ghappclient.ClientProvider, ghappclient.PersonalClientProvider) {
			return ghappclient.NewClient, ghappclient.NewPersonalClient
		},
		func() provider.RepoProvider {
			return provider.New(*reposBasePath, *gitLfsHostname)
		},
		func() (*sqlx.DB, error) {
			return db.Setup(*dbSourceAddr, true, "file://db/migrations")
		},
		func() config.GitHubAppConfig {
			return config.GitHubAppConfig{
				GitHubAppID:             *gitHubAppID,
				GitHubAppName:           *gitHubAppName,
				GitHubAppClientID:       *gitHubAppClientID,
				GitHubAppSecret:         *gitHubAppSecret,
				GitHubAppPrivateKeyPath: *gitHubAppPrivateKeyPath,
			}
		},
		func() (*session.Session, error) {
			awsSession, err := session.NewSession(
				&aws.Config{
					Region: aws.String("eu-north-1"),
				})
			return awsSession, err
		},
		func() (posthog.Client, error) {
			if *sendPostHogEvents {
				return posthog.NewWithConfig("ZuDRoGX9PgxGAZqY4RF9CCJJLpx14h3szUPzm7XBWSg", posthog.Config{Endpoint: "https://app.posthog.com"})
			} else {
				return ph.NewFakeClient(), nil
			}
		},
		events.NewInMemory,
		func(e events.EventReadWriter) events.EventReader {
			return e
		},
		executor.NewProvider,
		provider_acl.New,
		func(awsSession *session.Session) (queue.Queue, error) {
			if *localQueue {
				return queue.NewInMemory(logger), nil
			} else {
				return queue.NewSQS(logger, awsSession, *hostname, *queuePrefix)
			}
		},
		service_notification.NewPreferences,
		service_jwt.NewService,
		func(awsSession *session.Session) emails.Sender {
			if *enableTransactionalEmails {
				return emails.NewSES(awsSession)
			}
			return emails.NewLogs(logger)
		},
		func() service_ci.PublicAPIHostname {
			return service_ci.PublicAPIHostname(*publicApiHostname)
		},
		func() service_change.ExportBucketName {
			return service_change.ExportBucketName(exportBucketName)
		},
		func(repo db_workspace.Repository) db_workspace.WorkspaceReader {
			return repo
		},
		db_onetime.New,
		transactional.New,
		service_onetime.New,
		service_user.New,
		events.NewSender,
		notification_sender.NewNotificationSender,
		service_activity.New,
		activity_sender.NewActivitySender,
		ws_meta.NewWriterWithEvents,
		service_change.New,

		func() *service_github.ImporterQueue {
			return new(service_github.ImporterQueue)
		},
		func() *service_github.ClonerQueue {
			return new(service_github.ClonerQueue)
		},
		service_github.New,

		snapshotter.NewGitSnapshotter,
		meta_view.NewViewUpdatedFunc,
		service_statuses.New,
		service_ci.New,
		service_workspace.New,
		service_workspace_watchers.New,
		service_suggestion.New,
		service_presence.New,
		db_presence.NewRepo,
		service_codebase.New,
		service_auth.New,
		db_servicetokens.NewDatabase,
		service_servicetokens.New,
		db_onboarding.New,
		db_user.NewRepo,
		db_acl.NewACLRepository,
		db_codebase.NewRepo,
		db_codebase.NewCodebaseUserRepo,
		db_view.NewRepo,
		db_workspace.NewRepo,
		waitinglist.NewWaitingListRepo,
		acl.NewACLInterestRepo,
		instantintegration.NewInstantIntegrationInterestRepo,
		db_pki.NewRepo,
		db_snapshots.NewRepo,
		db_change.NewRepo,
		db_change.NewCommitRepository,
		db_comments.NewRepo,
		db_suggestion.New,
		db_gc.NewRepository,
		view_workspace_snapshot.NewRepo,
		db_github.NewGitHubInstallationRepo,
		db_github.NewGitHubRepositoryRepo,
		db_github.NewGitHubUserRepo,
		db_github.NewGitHubPRRepo,
		db_notification.NewRepository,
		db_mutagen.NewRepository,
		db_newsletter.NewNotificationSettingsRepository,
		db_activity.NewActivityRepo,
		db_review.NewReviewRepository,
		db_activity.NewActivityReadsRepo,
		db_notification.NewPeferenceRepository,
		db_keys.New,
		db_statuses.New,
		db_ci.NewCommitRepository,
		db_integrations.NewIntegrationDatabase,
		db_workspace_watchers.NewDB,
		service_comments.New,
		db_organization.New,
		db_organization.NewMember,
		service_sync.New,
		service_organization.New,
		func() http.DevelopmentAllowExtraCorsOrigin {
			return http.DevelopmentAllowExtraCorsOrigin(*developmentAllowExtraCorsOrigin)
		},
		db_buildkite.NewDatabase,
		service_buildkite.New,
		worker_gc.New,
		worker_snapshots.New,
		workers_ci.New,
		gitserver.New,
		workers_github.NewClonerQueue,
		workers_github.NewImporterQueue,
		db_license.New,
		db_license.NewValidationRepository,
		service_license.New,
		client_license.New,
		validator_license.New,
	}

	hooks := []di.Hook{
		di.Needs(http.Module),
		di.Needs(graphql.Module),
		di.Needs(api.Module),
	}

	for _, p := range providers {
		hooks = append(hooks, di.Provides(p))
	}

	c := di.NewModule(hooks...)
	c.Invoke(func(buildkiteService *service_buildkite.Service) {
		integrations.Register(integrations.ProviderTypeBuildkite, buildkiteService)
	})

	var apiServer *api.API
	if err := c.Build(&apiServer); err != nil {
		log.Fatalf("%+v", err)
	}

	if err := apiServer.Start(ctx, &api.Config{
		GitListenAddr:       *gitListenAddr,
		HTTPPProfListenAddr: *httpPprofListenAddr,
		MetricsListenAddr:   *metricsListenAddr,
		HTTPAddr:            *httpListenAddr,
	}); err != nil {
		logger.Fatal("failed to start api server", zap.Error(err))
	}
}

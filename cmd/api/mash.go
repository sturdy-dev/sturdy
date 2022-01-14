package main

import (
	"context"
	"flag"
	"log"
	_ "net/http/pprof"
	"os"

	"mash/db"
	"mash/pkg/api"
	module_api "mash/pkg/api/module"
	module_auth "mash/pkg/auth/module"
	module_author "mash/pkg/author/module"
	module_change "mash/pkg/change/module"
	service_change "mash/pkg/change/service"
	module_ci "mash/pkg/ci/module"
	service_ci "mash/pkg/ci/service"
	module_codebase_acl "mash/pkg/codebase/acl/module"
	module_codebase "mash/pkg/codebase/module"
	module_comments "mash/pkg/comments/module"
	"mash/pkg/di"
	"mash/pkg/emails"
	module_transactional "mash/pkg/emails/transactional/module"
	db_gc "mash/pkg/gc/db"
	worker_gc "mash/pkg/gc/worker"
	"mash/pkg/github/config"
	module_github "mash/pkg/github/module"
	"mash/pkg/gitserver"
	"mash/pkg/graphql"
	"mash/pkg/http"
	"mash/pkg/http/oss"
	module_integrations "mash/pkg/integrations/module"
	db_keys "mash/pkg/jwt/keys/db"
	service_jwt "mash/pkg/jwt/service"
	module_license "mash/pkg/license/module"
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
	module_statuses "mash/pkg/statuses/module"
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
	module_workspace "mash/pkg/workspace/module"
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
		func(e events.EventReadWriter) events.EventReader {
			return e
		},
		func(awsSession *session.Session) (queue.Queue, error) {
			if *localQueue {
				return queue.NewInMemory(logger), nil
			} else {
				return queue.NewSQS(logger, awsSession, *hostname, *queuePrefix)
			}
		},
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
		func() oss.DevelopmentAllowExtraCorsOrigin {
			return oss.DevelopmentAllowExtraCorsOrigin(*developmentAllowExtraCorsOrigin)
		},
		events.NewInMemory,
		executor.NewProvider,
		service_notification.NewPreferences,
		service_jwt.NewService,
		db_onetime.New,
		service_onetime.New,
		service_user.New,
		events.NewSender,
		notification_sender.NewNotificationSender,
		service_activity.New,
		activity_sender.NewActivitySender,
		ws_meta.NewWriterWithEvents,
		snapshotter.NewGitSnapshotter,
		meta_view.NewViewUpdatedFunc,
		service_statuses.New,
		service_workspace_watchers.New,
		service_suggestion.New,
		service_presence.New,
		db_presence.NewRepo,
		db_servicetokens.NewDatabase,
		service_servicetokens.New,
		db_onboarding.New,
		db_user.NewRepo,
		db_view.NewRepo,
		db_workspace.NewRepo,
		waitinglist.NewWaitingListRepo,
		acl.NewACLInterestRepo,
		instantintegration.NewInstantIntegrationInterestRepo,
		db_pki.NewRepo,
		db_snapshots.NewRepo,
		db_suggestion.New,
		db_gc.NewRepository,
		view_workspace_snapshot.NewRepo,
		db_notification.NewRepository,
		db_mutagen.NewRepository,
		db_newsletter.NewNotificationSettingsRepository,
		db_activity.NewActivityRepo,
		db_review.NewReviewRepository,
		db_activity.NewActivityReadsRepo,
		db_notification.NewPeferenceRepository,
		db_keys.New,
		db_statuses.New,
		db_workspace_watchers.NewDB,
		db_organization.New,
		db_organization.NewMember,
		service_sync.New,
		service_organization.New,
		worker_gc.New,
		worker_snapshots.New,
		gitserver.New,
	}

	mainModule := func(c *di.Container) {
		for _, provider := range providers {
			c.Register(provider)
		}

		c.Import(module_api.Module)
		c.Import(module_auth.Module)
		c.Import(module_author.Module)
		c.Import(module_change.Module)
		c.Import(module_ci.Module)
		c.Import(module_codebase.Module)
		c.Import(module_codebase_acl.Module)
		c.Import(module_comments.Module)
		c.Import(module_integrations.Module)

		// todo: continue importing here

		c.Import(module_statuses.Module)
		c.Import(module_transactional.Module)
		c.Import(http.Module)
		c.Import(module_github.Module)
		c.Import(module_workspace.Module)
		c.Import(graphql.Module)
		c.Import(module_license.Module)
	}

	var apiServer api.API
	if err := di.Init(&apiServer, mainModule); err != nil {
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

package main

import (
	"context"
	"flag"
	"log"
	_ "net/http/pprof"
	"os"

	module_analytics "mash/pkg/analytics/module"
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
	"mash/pkg/db"
	"mash/pkg/di"
	"mash/pkg/emails"
	module_transactional "mash/pkg/emails/transactional/module"
	module_features "mash/pkg/features/module"
	module_file "mash/pkg/file/module"
	module_gc "mash/pkg/gc/module"
	"mash/pkg/github/config"
	module_github "mash/pkg/github/module"
	module_gitserver "mash/pkg/gitserver"
	module_graphql "mash/pkg/graphql"
	module_http "mash/pkg/http"
	"mash/pkg/http/oss"
	module_installations "mash/pkg/installations/module"
	module_integrations "mash/pkg/integrations/module"
	module_jwt "mash/pkg/jwt/module"
	module_license "mash/pkg/license/module"
	"mash/pkg/metrics/zapprometheus"
	module_mutagen "mash/pkg/mutagen/module"
	module_newsletter "mash/pkg/newsletter/module"
	module_notification "mash/pkg/notification/module"
	module_onboarding "mash/pkg/onboarding/module"
	module_onetime "mash/pkg/onetime/module"
	module_organization "mash/pkg/organization/module"
	module_pki "mash/pkg/pki/module"
	module_presence "mash/pkg/presence/module"
	"mash/pkg/queue"
	module_review "mash/pkg/review/module"
	module_serverstatus "mash/pkg/serverstatus/module"
	module_servicetokens "mash/pkg/servicetokens/module"
	module_snapshots "mash/pkg/snapshots/module"
	db_statuses "mash/pkg/statuses/db"
	module_statuses "mash/pkg/statuses/module"
	service_statuses "mash/pkg/statuses/service"
	module_suggestions "mash/pkg/suggestions/module"
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
	gitHubAppID := flag.Int64("github-app-id", 122610, "")
	gitHubAppName := flag.String("github-app-name", "sturdy-gustav-localhost", "")
	gitHubAppClientID := flag.String("github-app-client-id", "", "")
	gitHubAppSecret := flag.String("github-app-secret", "", "")
	gitHubAppPrivateKeyPath := flag.String("github-app-private-key-path", "", "")
	gitLfsHostname := flag.String("git-lfs-hostname", "localhost:8888", "")
	enableTransactionalEmails := flag.Bool("enable-transactional-emails", false, "")
	exportBucketName := flag.String("export-bucket-name", "", "the S3 bucket to be used for change exports")
	developmentAllowExtraCorsOrigin := flag.String("development-allow-extra-cors-origin", "", "Additional CORS origin to be allowed")
	localQueue := flag.Bool("use-local-queues", false, "If set, local queue will be used instead of SQS")

	// deprecated flags
	_ = flag.Bool("send-posthog-events", false, "")
	_ = flag.Bool("send-invites-worker", false, "")
	_ = flag.String("gmail-token-json-path", "", "used by the invite email sender")
	_ = flag.Bool("unauthenticated-graphql-introspection", false, "")
	_ = flag.String("gmail-credentials-json-path", "", "used by the invite email sender")

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
		func() context.Context {
			return ctx
		},
		func() *zap.Logger { return logger },
		func() provider.RepoProvider {
			return provider.New(*reposBasePath, *gitLfsHostname)
		},
		func() (*sqlx.DB, error) {
			return db.Setup(*dbSourceAddr)
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
		service_user.New,
		events.NewSender,
		service_activity.New,
		activity_sender.NewActivitySender,
		ws_meta.NewWriterWithEvents,
		meta_view.NewViewUpdatedFunc,
		service_statuses.New,
		service_workspace_watchers.New,
		db_user.NewRepo,
		db_view.NewRepo,
		db_workspace.NewRepo,
		waitinglist.NewWaitingListRepo,
		acl.NewACLInterestRepo,
		instantintegration.NewInstantIntegrationInterestRepo,
		view_workspace_snapshot.NewRepo,
		db_activity.NewActivityRepo,
		db_activity.NewActivityReadsRepo,
		db_statuses.New,
		db_workspace_watchers.NewDB,
		service_sync.New,
	}

	mainModule := func(c *di.Container) {
		for _, provider := range providers {
			c.Register(provider)
		}

		c.Import(module_analytics.Module)
		c.Import(module_api.Module)
		c.Import(module_auth.Module)
		c.Import(module_author.Module)
		c.Import(module_change.Module)
		c.Import(module_ci.Module)
		c.Import(module_codebase.Module)
		c.Import(module_codebase_acl.Module)
		c.Import(module_comments.Module)
		c.Import(module_features.Module)
		c.Import(module_file.Module)
		c.Import(module_gc.Module)
		c.Import(module_github.Module)
		c.Import(module_gitserver.Module)
		c.Import(module_graphql.Module)
		c.Import(module_http.Module)
		c.Import(module_installations.Module)
		c.Import(module_integrations.Module)
		c.Import(module_jwt.Module)
		c.Import(module_license.Module)
		c.Import(module_mutagen.Module)
		c.Import(module_newsletter.Module)
		c.Import(module_notification.Module)
		c.Import(module_onboarding.Module)
		c.Import(module_onetime.Module)
		c.Import(module_organization.Module)
		c.Import(module_pki.Module)
		c.Import(module_presence.Module)
		c.Import(module_review.Module)
		c.Import(module_servicetokens.Module)
		c.Import(module_snapshots.Module)
		c.Import(module_statuses.Module)
		c.Import(module_suggestions.Module)

		// todo: continue importing here

		c.Import(module_transactional.Module)
		c.Import(module_workspace.Module)
		c.Import(module_serverstatus.Module)
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

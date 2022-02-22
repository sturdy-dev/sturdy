package module

import (
	module_analytics "getsturdy.com/api/pkg/analytics/module"
	module_auth "getsturdy.com/api/pkg/auth/module"
	module_author "getsturdy.com/api/pkg/author/module"
	module_aws "getsturdy.com/api/pkg/aws/module"
	module_change "getsturdy.com/api/pkg/change/module"
	module_ci "getsturdy.com/api/pkg/ci/module"
	module_codebase_acl "getsturdy.com/api/pkg/codebase/acl/module"
	module_codebase "getsturdy.com/api/pkg/codebase/module"
	module_comments "getsturdy.com/api/pkg/comments/module"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
	module_emails "getsturdy.com/api/pkg/emails/module"
	module_email_transactional "getsturdy.com/api/pkg/emails/transactional/module"
	module_events "getsturdy.com/api/pkg/events"
	module_features "getsturdy.com/api/pkg/features/module"
	module_file "getsturdy.com/api/pkg/file/module"
	module_gc "getsturdy.com/api/pkg/gc/module"
	module_gitserver "getsturdy.com/api/pkg/gitserver"
	module_graphql "getsturdy.com/api/pkg/graphql"
	module_http "getsturdy.com/api/pkg/http/module"
	module_installations "getsturdy.com/api/pkg/installations/module"
	module_installations_statistics "getsturdy.com/api/pkg/installations/statistics/module"
	module_integrations "getsturdy.com/api/pkg/integrations/module"
	module_jwt "getsturdy.com/api/pkg/jwt/module"
	module_license "getsturdy.com/api/pkg/licenses/module"
	module_logger "getsturdy.com/api/pkg/logger/module"
	"getsturdy.com/api/pkg/metrics"
	module_mutagen "getsturdy.com/api/pkg/mutagen/module"
	module_newsletter "getsturdy.com/api/pkg/newsletter/module"
	module_notification "getsturdy.com/api/pkg/notification/module"
	module_onboarding "getsturdy.com/api/pkg/onboarding/module"
	module_onetime "getsturdy.com/api/pkg/onetime/module"
	module_organization "getsturdy.com/api/pkg/organization/module"
	module_pki "getsturdy.com/api/pkg/pki/module"
	"getsturdy.com/api/pkg/pprof"
	module_presence "getsturdy.com/api/pkg/presence/module"
	module_review "getsturdy.com/api/pkg/review/module"
	module_servicetokens "getsturdy.com/api/pkg/servicetokens/module"
	module_statuses "getsturdy.com/api/pkg/statuses/module"
	module_suggestions "getsturdy.com/api/pkg/suggestions/module"
	module_sync "getsturdy.com/api/pkg/sync/module"
	module_user "getsturdy.com/api/pkg/users/module"
	module_view "getsturdy.com/api/pkg/view/module"
	module_waitinglist "getsturdy.com/api/pkg/waitinglist"
	module_workspace_activity "getsturdy.com/api/pkg/workspace/activity/module"
	module_workspace "getsturdy.com/api/pkg/workspace/module"
	module_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/module"
	module_vcs "getsturdy.com/api/vcs/module"
)

func common(c *di.Container) {
	c.Import(db.Module)
	c.Import(metrics.Module)
	c.Import(pprof.Module)

	c.Import(module_aws.Module)
	c.Import(module_analytics.Module)
	c.Import(module_auth.Module)
	c.Import(module_author.Module)
	c.Import(module_change.Module)
	c.Import(module_ci.Module)
	c.Import(module_codebase.Module)
	c.Import(module_codebase_acl.Module)
	c.Import(module_comments.Module)

	c.Import(module_emails.Module)
	c.Import(module_email_transactional.Module)
	c.Import(module_events.Module)
	c.Import(module_features.Module)
	c.Import(module_file.Module)
	c.Import(module_gc.Module)
	c.Import(module_gitserver.Module)
	c.Import(module_graphql.Module)
	c.Import(module_http.Module)
	c.Import(module_installations.Module)
	c.Import(module_installations_statistics.Module)
	c.Import(module_integrations.Module)
	c.Import(module_jwt.Module)
	c.Import(module_logger.Module)
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
	c.Import(module_statuses.Module)
	c.Import(module_suggestions.Module)
	c.Import(module_sync.Module)
	c.Import(module_user.Module)
	c.Import(module_view.Module)
	c.Import(module_waitinglist.Module)
	c.Import(module_workspace.Module)
	c.Import(module_workspace_activity.Module)
	c.Import(module_workspace_watchers.Module)
	c.Import(module_vcs.Module)
}

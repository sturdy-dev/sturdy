package graphql

import (
	graphql_author "mash/pkg/author/graphql"
	graphql_change "mash/pkg/change/graphql"
	graphql_acl "mash/pkg/codebase/acl/graphql"
	graphql_codebase "mash/pkg/codebase/graphql"
	graphql_comments "mash/pkg/comments/graphql"
	"mash/pkg/digutils"
	graphql_file "mash/pkg/file/graphql"
	graphql_github "mash/pkg/github/graphql"
	graphql_pr "mash/pkg/github/graphql/pr"
	"mash/pkg/graphql/resolvers"
	graphql_buildkite "mash/pkg/integrations/buildkite/graphql"
	graphql_ci "mash/pkg/integrations/graphql"
	graphql_license "mash/pkg/license/graphql"
	graphql_notification "mash/pkg/notification/graphql"
	graphql_onboarding "mash/pkg/onboarding/graphql"
	graphql_organization "mash/pkg/organization/graphql"
	graphql_pki "mash/pkg/pki/graphql"
	graphql_presence "mash/pkg/presence/graphql"
	graphql_review "mash/pkg/review/graphql"
	graphql_servicetokens "mash/pkg/servicetokens/graphql"
	graphql_statuses "mash/pkg/statuses/graphql"
	graphql_suggestion "mash/pkg/suggestions/graphql"
	graphql_user "mash/pkg/user/graphql"
	graphql_view "mash/pkg/view/graphql"
	graphql_workspace_activity "mash/pkg/workspace/activity/graphql"
	graphql_workspace "mash/pkg/workspace/graphql"
	graphql_workspace_watchers "mash/pkg/workspace/watchers/graphql"
)

var Invokers = []interface{}{}

var Providers = []interface{}{
	New,
	graphql_acl.NewResolver,
	graphql_author.NewResolver,
	graphql_change.NewResolver,
	graphql_codebase.NewResolver,
	graphql_comments.NewResolver,
	graphql_github.NewResolver,
	graphql_change.NewFileDiffRootResolver,
	graphql_suggestion.New,
	graphql_notification.NewResolver,
	graphql_user.NewResolver,
	graphql_view.NewViewStatusRootResolver,
	graphql_view.NewResolver,
	graphql_pr.NewResolver,
	graphql_statuses.New,
	graphql_workspace_watchers.NewRootResolver,
	graphql_workspace.NewResolver,
	graphql_workspace_activity.New,
	graphql_review.New,
	graphql_file.NewFileRootResolver,
	graphql_presence.NewRootResolver,
	graphql_ci.NewRootResolver,
	graphql_onboarding.NewRootResolver,
	graphql_pki.NewResolver,
	graphql_servicetokens.New,
	graphql_buildkite.New,
	graphql_organization.New,
	graphql_license.New,
	graphql_github.NewGitHubAppRootResolver,
}

func init() {
	cycleDeps := []interface{}{
		new(resolvers.ACLRootResolver),
		new(resolvers.AuthorRootResolver),
		new(resolvers.BuildkiteInstantIntegrationRootResolver),
		new(resolvers.ChangeRootResolver),
		new(resolvers.CodebaseGitHubIntegrationRootResolver),
		new(resolvers.CodebaseRootResolver),
		new(resolvers.CommentRootResolver),
		new(resolvers.FileDiffRootResolver),
		new(resolvers.FileRootResolver),
		new(resolvers.IntegrationRootResolver),
		new(resolvers.NotificationRootResolver),
		new(resolvers.OnboardingRootResolver),
		new(resolvers.OrganizationRootResolver),
		new(resolvers.PKIRootResolver),
		new(resolvers.GitHubPullRequestRootResolver),
		new(resolvers.PresenceRootResolver),
		new(resolvers.ReviewRootResolver),
		new(resolvers.ServiceTokensRootResolver),
		new(resolvers.StatusesRootResolver),
		new(resolvers.SuggestionRootResolver),
		new(resolvers.UserRootResolver),
		new(resolvers.ViewRootResolver),
		new(resolvers.WorkspaceActivityRootResolver),
		new(resolvers.WorkspaceRootResolver),
		new(resolvers.WorkspaceWatcherRootResolver),
	}

	for _, dep := range cycleDeps {
		provider, invoker := digutils.ResolveCycleFor(dep)
		Providers = append(Providers, provider)
		Invokers = append(Invokers, invoker)
	}
}

package graphql

import (
	graphql_author "mash/pkg/author/graphql"
	graphql_change "mash/pkg/change/graphql"
	graphql_acl "mash/pkg/codebase/acl/graphql"
	graphql_codebase "mash/pkg/codebase/graphql"
	graphql_comments "mash/pkg/comments/graphql"
	"mash/pkg/di"
	graphql_features "mash/pkg/features/graphql"
	graphql_file "mash/pkg/file/graphql"
	graphql_buildkite "mash/pkg/integrations/buildkite/graphql"
	graphql_ci "mash/pkg/integrations/graphql"
	graphql_notification "mash/pkg/notification/graphql"
	graphql_onboarding "mash/pkg/onboarding/graphql"
	graphql_organization "mash/pkg/organization/graphql"
	graphql_pki "mash/pkg/pki/graphql"
	graphql_presence "mash/pkg/presence/graphql"
	graphql_review "mash/pkg/review/graphql"
	graphql_servicetokens "mash/pkg/servicetokens/graphql"
	graphql_suggestion "mash/pkg/suggestions/graphql"
	graphql_user "mash/pkg/user/graphql"
	graphql_view "mash/pkg/view/graphql"
	graphql_workspace_activity "mash/pkg/workspace/activity/graphql"
	graphql_workspace "mash/pkg/workspace/graphql"
	graphql_workspace_watchers "mash/pkg/workspace/watchers/graphql"
)

var Module = di.NewModule(
	di.Provides(NewRootResolver),

	di.ProvidesCycle(graphql_acl.NewResolver),
	di.ProvidesCycle(graphql_author.NewResolver),
	di.ProvidesCycle(graphql_change.NewResolver),
	di.ProvidesCycle(graphql_comments.NewResolver),
	di.ProvidesCycle(graphql_change.NewFileDiffRootResolver),
	di.ProvidesCycle(graphql_suggestion.New),
	di.ProvidesCycle(graphql_notification.NewResolver),
	di.ProvidesCycle(graphql_user.NewResolver),
	di.ProvidesCycle(graphql_view.NewViewStatusRootResolver),
	di.ProvidesCycle(graphql_view.NewResolver),
	di.ProvidesCycle(graphql_workspace_watchers.NewRootResolver),
	di.ProvidesCycle(graphql_workspace.NewResolver),
	di.ProvidesCycle(graphql_workspace_activity.New),
	di.ProvidesCycle(graphql_review.New),
	di.ProvidesCycle(graphql_file.NewFileRootResolver),
	di.ProvidesCycle(graphql_presence.NewRootResolver),
	di.ProvidesCycle(graphql_ci.NewRootResolver),
	di.ProvidesCycle(graphql_onboarding.NewRootResolver),
	di.ProvidesCycle(graphql_pki.NewResolver),
	di.ProvidesCycle(graphql_servicetokens.New),
	di.ProvidesCycle(graphql_buildkite.New),
	di.ProvidesCycle(graphql_organization.New),

	di.Needs(graphql_codebase.Module),
	di.Needs(graphql_features.Module),
)

package graphql

import (
	graphql_activity "getsturdy.com/api/pkg/activity/graphql"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	graphql_changes "getsturdy.com/api/pkg/changes/graphql"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	graphql_codebases "getsturdy.com/api/pkg/codebases/graphql"
	db_comments "getsturdy.com/api/pkg/comments/db"
	graphql_comments "getsturdy.com/api/pkg/comments/graphql"
	"getsturdy.com/api/pkg/di"
	graphql_github_pr "getsturdy.com/api/pkg/github/graphql/pr"
	"getsturdy.com/api/pkg/graphql/resolvers"
	graphql_presence "getsturdy.com/api/pkg/presence/graphql"
	graphql_review "getsturdy.com/api/pkg/review/graphql"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	graphql_suggestions "getsturdy.com/api/pkg/suggestions/graphql"
	graphql_rebase "getsturdy.com/api/pkg/sync/graphql"
	db_view "getsturdy.com/api/pkg/view/db"
	graphql_view "getsturdy.com/api/pkg/view/graphql"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	graphql_workspace_watchers "getsturdy.com/api/pkg/workspaces/watchers/graphql"
)

func Module(c *di.Container) {
	c.Import(db_workspaces.Module)
	c.Import(db_codebases.Module)
	c.Import(db_view.Module)
	c.Import(db_comments.Module)
	c.Import(db_snapshots.Module)
	c.Import(graphql_codebases.Module)
	c.Import(graphql_author.Module)
	c.Import(graphql_view.Module)
	c.Import(graphql_comments.Module)
	c.Import(graphql_suggestions.Module)
	c.Import(graphql_github_pr.Module)
	c.Import(graphql_changes.Module)
	c.Import(graphql_activity.Module)
	c.Import(graphql_review.Module)
	c.Import(graphql_presence.Module)
	c.Import(graphql_workspace_watchers.Module)
	c.Import(graphql_rebase.Module)

	c.Register(NewResolver)

	// populate cyclic resolver
	c.Import(resolvers.Module)
	c.Decorate(func(rp *resolvers.WorkspaceRootResolver, rv resolvers.WorkspaceRootResolver) *resolvers.WorkspaceRootResolver {
		*rp = rv
		return &rv
	})
}

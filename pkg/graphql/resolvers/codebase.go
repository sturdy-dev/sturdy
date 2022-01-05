package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type CodebaseArgs struct {
	ID      *graphql.ID
	ShortID *graphql.ID
}

type CodebaseRootResolver interface {
	Codebase(ctx context.Context, args CodebaseArgs) (CodebaseResolver, error)
	Codebases(ctx context.Context) ([]CodebaseResolver, error)

	// Subscriptions
	UpdatedCodebase(ctx context.Context) (<-chan CodebaseResolver, error)

	// Mutations
	UpdateCodebase(ctx context.Context, args UpdateCodebaseArgs) (CodebaseResolver, error)
}

type CodebaseViewsArgs struct {
	IncludeOthers *bool
}

type UpdateCodebaseArgs struct {
	Input UpdateCodebaseInput
}

type UpdateCodebaseInput struct {
	ID                 graphql.ID
	Name               *string
	DisableInviteCode  *bool
	GenerateInviteCode *bool
	Archive            *bool
	IsPublic           *bool
}

type CodebaseResolver interface {
	ID() graphql.ID
	Name() string
	ShortID() graphql.ID
	Description() string
	InviteCode() *string
	CreatedAt() int32
	ArchivedAt() *int32
	LastUpdatedAt() *int32
	Workspaces(ctx context.Context) ([]WorkspaceResolver, error)
	Members(context.Context) ([]AuthorResolver, error)
	Views(ctx context.Context, args CodebaseViewsArgs) ([]ViewResolver, error)
	LastUsedView(ctx context.Context) (ViewResolver, error)
	GitHubIntegration() (CodebaseGitHubIntegrationResolver, error)
	IsReady() bool
	ACL(context.Context) (ACLResolver, error)
	Changes(ctx context.Context, args *CodebaseChangesArgs) ([]ChangeResolver, error)
	Readme(ctx context.Context) (FileResolver, error)
	File(ctx context.Context, args CodebaseFileArgs) (FileOrDirectoryResolver, error)
	Integrations(ctx context.Context, args IntegrationsArgs) ([]IntegrationResolver, error)
	IsPublic() bool
}

type CodebaseChangesArgs struct {
	Input *CodebaseChangesInput
}

type CodebaseChangesInput struct {
	Limit *int32
}

type CodebaseFileArgs struct {
	Path string
}

type IntegrationsArgs struct {
	ID *graphql.ID
}

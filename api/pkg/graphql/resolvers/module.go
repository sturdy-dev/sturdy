package resolvers

import (
	"getsturdy.com/api/pkg/di"
)

// Here we provide pointer to all cyclic graphql resolvers.
func Module(c *di.Container) {
	c.Register(func() *WorkspaceRootResolver { return new(WorkspaceRootResolver) })
	c.Register(func() *AuthorRootResolver { return new(AuthorRootResolver) })
	c.Register(func() *ChangeRootResolver { return new(ChangeRootResolver) })
	c.Register(func() *CommentRootResolver { return new(CommentRootResolver) })
	c.Register(func() *StatusesRootResolver { return new(StatusesRootResolver) })
	c.Register(func() *CodebaseRootResolver { return new(CodebaseRootResolver) })
	c.Register(func() *OrganizationRootResolver { return new(OrganizationRootResolver) })
	c.Register(func() *ViewRootResolver { return new(ViewRootResolver) })
	c.Register(func() *UserRootResolver { return new(UserRootResolver) })
}

package enterprise

import (
	"context"
	"database/sql"
	"errors"

	"mash/pkg/codebase/graphql/oss"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
)

type CodebaseRootResolver struct {
	*oss.CodebaseRootResolver

	codebaseGitHubIntegrationResolver *resolvers.CodebaseGitHubIntegrationRootResolver
}

func NewCodebaseRootResolver(
	ossResolver *oss.CodebaseRootResolver,

	codebaseGitHubIntegrationResolver *resolvers.CodebaseGitHubIntegrationRootResolver,
) *CodebaseRootResolver {
	return &CodebaseRootResolver{
		CodebaseRootResolver: ossResolver,

		codebaseGitHubIntegrationResolver: codebaseGitHubIntegrationResolver,
	}
}

func (r *CodebaseRootResolver) Codebase(ctx context.Context, args resolvers.CodebaseArgs) (resolvers.CodebaseResolver, error) {
	ossCodebase, err := r.CodebaseRootResolver.Codebase(ctx, args)
	if err != nil {
		return ossCodebase, err
	}
	return &codebaseResolver{
		CodebaseResolver: ossCodebase,
		root:             r,
	}, nil
}

func (r *CodebaseRootResolver) Codebases(ctx context.Context) ([]resolvers.CodebaseResolver, error) {
	ossCodebases, err := r.CodebaseRootResolver.Codebases(ctx)
	if err != nil {
		return ossCodebases, err
	}
	enterpriseCodebases := make([]resolvers.CodebaseResolver, 0, len(ossCodebases))
	for _, ossCodebase := range ossCodebases {
		enterpriseCodebases = append(enterpriseCodebases, &codebaseResolver{
			CodebaseResolver: ossCodebase,
			root:             r,
		})
	}

	return enterpriseCodebases, nil
}

func (r *CodebaseRootResolver) UpdatedCodebase(ctx context.Context) (<-chan resolvers.CodebaseResolver, error) {
	responseChan, err := r.CodebaseRootResolver.UpdatedCodebase(ctx)
	if err != nil {
		return responseChan, err
	}

	c := make(chan resolvers.CodebaseResolver, 100)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c)
			case ossCodebase := <-responseChan:
				c <- &codebaseResolver{
					CodebaseResolver: ossCodebase,
					root:             r,
				}
			}
		}
	}()
	return c, nil
}

func (r *CodebaseRootResolver) UpdateCodebase(ctx context.Context, args resolvers.UpdateCodebaseArgs) (resolvers.CodebaseResolver, error) {
	ossCodebase, err := r.CodebaseRootResolver.UpdateCodebase(ctx, args)
	if err != nil {
		return ossCodebase, err
	}
	return &codebaseResolver{
		CodebaseResolver: ossCodebase,
		root:             r,
	}, nil
}

type codebaseResolver struct {
	resolvers.CodebaseResolver

	root *CodebaseRootResolver
}

func (r *codebaseResolver) GitHubIntegration(ctx context.Context) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	resolver, err := (*r.root.codebaseGitHubIntegrationResolver).InternalCodebaseGitHubIntegration(ctx, r.ID())
	switch {
	case err == nil:
		return resolver, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

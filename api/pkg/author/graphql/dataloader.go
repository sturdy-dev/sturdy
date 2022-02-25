package graphql

import (
	"context"

	gqldataloader "getsturdy.com/api/pkg/graphql/dataloader"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/users"

	"github.com/graph-gophers/dataloader/v6"
	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type authorDataloader struct {
	resolver *AuthorRootResolver
	loader   *dataloader.Loader
}

func NewDataloader(resolver *AuthorRootResolver, logger *zap.Logger) resolvers.AuthorRootResolver {
	return &authorDataloader{
		resolver: resolver,
		loader: dataloader.NewBatchedLoader(
			batchFunction(resolver),
			dataloader.WithCache(gqldataloader.NewContextCache(logger)),
		),
	}
}

func (al *authorDataloader) InternalAuthorFromNameAndEmail(ctx context.Context, name, email string) resolvers.AuthorResolver {
	return al.resolver.InternalAuthorFromNameAndEmail(ctx, name, email)
}

func (al *authorDataloader) Author(ctx context.Context, id graphql.ID) (resolvers.AuthorResolver, error) {
	thunk := al.loader.Load(ctx, dataloader.StringKey(id))
	r, err := thunk()
	if err != nil {
		return nil, err
	}
	return r.(resolvers.AuthorResolver), nil
}

func batchFunction(resolver *AuthorRootResolver) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var (
			results = make([]*dataloader.Result, 0, len(keys))
		)

		for _, key := range keys {
			res, err := resolver.Author(users.ID(key.String()))
			if err != nil {
				results = append(results, &dataloader.Result{Error: err})
			} else {
				results = append(results, &dataloader.Result{Data: res})
			}
		}

		return results
	}
}

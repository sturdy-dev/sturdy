package graphql

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	gqldataloader "getsturdy.com/api/pkg/graphql/dataloader"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/dataloader/v6"
	"go.uber.org/zap"
)

type userDataloader struct {
	resolver *userRootResolver
	loader   *dataloader.Loader
}

func NewDataloader(resolver *userRootResolver, logger *zap.Logger) resolvers.UserRootResolver {
	return &userDataloader{
		resolver: resolver,
		loader: dataloader.NewBatchedLoader(
			batchFunction(resolver),
			dataloader.WithCache(gqldataloader.NewContextCache(logger)),
		),
	}
}

func (dl *userDataloader) User(ctx context.Context) (resolvers.UserResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return dl.InternalUser(ctx, userID)
}

func (dl *userDataloader) UpdateUser(ctx context.Context, args resolvers.UpdateUserArgs) (resolvers.UserResolver, error) {
	r, err := dl.resolver.UpdateUser(ctx, args)
	key := dataloader.StringKey(r.ID())
	dl.loader.Clear(ctx, key).Prime(ctx, key, r)
	return r, err
}

func (dl *userDataloader) VerifyEmail(ctx context.Context, args resolvers.VerifyEmailArgs) (resolvers.UserResolver, error) {
	r, err := dl.resolver.VerifyEmail(ctx, args)
	key := dataloader.StringKey(r.ID())
	dl.loader.Clear(ctx, key).Prime(ctx, key, r)
	return r, err
}

func (dl *userDataloader) InternalUser(ctx context.Context, userID string) (resolvers.UserResolver, error) {
	thunk := dl.loader.Load(ctx, dataloader.StringKey(userID))
	u, err := thunk()
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return u.(resolvers.UserResolver), nil
}

func batchFunction(resolver *userRootResolver) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var (
			results = make([]*dataloader.Result, 0, len(keys))
		)

		for _, key := range keys {
			user, err := resolver.InternalUser(key.String())
			if err != nil {
				results = append(results, &dataloader.Result{Error: err})
			} else {
				results = append(results, &dataloader.Result{Data: user})
			}
		}

		return results
	}
}

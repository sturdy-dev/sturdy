package graphql

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_user "getsturdy.com/api/pkg/users/enterprise/cloud/service"
	"getsturdy.com/api/pkg/users/graphql"
)

type userRootResolver struct {
	*graphql.UserDataloader

	userService *service_user.Service
}

func NewResolver(
	userDataloader *graphql.UserDataloader,
	userService *service_user.Service,
) *userRootResolver {
	return &userRootResolver{
		UserDataloader: userDataloader,
		userService:    userService,
	}
}

func (r *userRootResolver) VerifyEmail(ctx context.Context, args resolvers.VerifyEmailArgs) (resolvers.UserResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	user, err := r.userService.VerifyEmail(ctx, userID, args.Input.Token)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return r.InternalUser(ctx, user.ID)
}

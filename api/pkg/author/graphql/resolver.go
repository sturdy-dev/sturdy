package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/users"
	service_users "getsturdy.com/api/pkg/users/service"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type AuthorRootResolver struct {
	userService service_users.Service
}

func NewResolver(userService service_users.Service, logger *zap.Logger) resolvers.AuthorRootResolver {
	userRoot := &AuthorRootResolver{
		userService: userService,
	}
	return NewDataloader(userRoot, logger)
}

func (r *AuthorRootResolver) Author(ctx context.Context, id users.ID) (resolvers.AuthorResolver, error) {
	uu, err := r.userService.GetByID(ctx, id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &AuthorResolver{user: uu}, nil
}

func (r *AuthorRootResolver) InternalAuthorFromNameAndEmail(_ context.Context, name, email string) resolvers.AuthorResolver {
	return &authorNameEmailResolver{name, email}
}

type AuthorResolver struct {
	user *users.User
}

func (r *AuthorResolver) ID() graphql.ID {
	return graphql.ID(r.user.ID)
}

func (r *AuthorResolver) Name() string {
	return r.user.Name
}

func (r *AuthorResolver) AvatarUrl() *string {
	return r.user.AvatarURL
}

func (r *AuthorResolver) Email() string {
	return r.user.Email
}

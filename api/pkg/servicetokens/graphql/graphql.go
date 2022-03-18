package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebases"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/servicetokens"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"

	"github.com/graph-gophers/graphql-go"
)

type rootResolver struct {
	authService          *service_auth.Service
	serviceTokensService *service_servicetokens.Service
	codebaseService      *service_codebase.Service
}

func New(
	authService *service_auth.Service,
	serviceTokensService *service_servicetokens.Service,
	codebaseService *service_codebase.Service,
) resolvers.ServiceTokensRootResolver {
	return &rootResolver{
		authService:          authService,
		serviceTokensService: serviceTokensService,
		codebaseService:      codebaseService,
	}
}

func (r *rootResolver) CreateServiceToken(ctx context.Context, args resolvers.CreateServiceTokenArgs) (resolvers.ServiceTokenResovler, error) {
	codebase, err := r.codebaseService.GetByShortID(ctx, codebases.ShortCodebaseID(args.Input.ShortCodebaseID))
	if err != nil {
		return nil, gqlerror.Error(fmt.Errorf("codebase not found: %w", err))
	}

	if err := r.authService.CanWrite(ctx, codebase); err != nil {
		return nil, gqlerror.Error(err)
	}

	plainTextToken, token, err := r.serviceTokensService.Create(ctx, codebase.ID, args.Input.Name)
	if err != nil {
		return nil, gqlerror.Error(fmt.Errorf("failed to create token: %w", err))
	}

	return &resolver{
		token:          token,
		plainTextToken: &plainTextToken,
	}, nil
}

type resolver struct {
	plainTextToken *string
	token          *servicetokens.Token
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.token.ID)
}

func (r *resolver) Name() string {
	return r.token.Name
}

func (r *resolver) CreatedAt() int32 {
	return int32(r.token.CreatedAt.Unix())
}

func (r *resolver) LastUsedAt() *int32 {
	if r.token.LastUsedAt == nil {
		return nil
	}
	luat := int32(r.token.LastUsedAt.Unix())
	return &luat
}

func (r *resolver) Token() *string {
	return r.plainTextToken
}

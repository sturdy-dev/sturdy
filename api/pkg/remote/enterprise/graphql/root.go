package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebases"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/crypto"
	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/remote/enterprise/service"
	service_user "getsturdy.com/api/pkg/users/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
)

type remoteRootResolver struct {
	service            *service.EnterpriseService
	workspaceService   service_workspace.Service
	authService        *service_auth.Service
	codebaseService    *service_codebase.Service
	userService        service_user.Service
	cryptoRootResolver resolvers.CryptoRootResolver
}

func New(
	service *service.EnterpriseService,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	codebaseService *service_codebase.Service,
	userService service_user.Service,
	cryptoRootResolver resolvers.CryptoRootResolver,
) resolvers.RemoteRootResolver {
	return &remoteRootResolver{
		service:            service,
		workspaceService:   workspaceService,
		authService:        authService,
		codebaseService:    codebaseService,
		userService:        userService,
		cryptoRootResolver: cryptoRootResolver,
	}
}

func (r *remoteRootResolver) InternalRemoteByCodebaseID(ctx context.Context, codebaseID codebases.ID) (resolvers.RemoteResolver, error) {
	cb, err := r.codebaseService.GetByID(ctx, codebaseID)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerror.Error(err)
	}

	rem, err := r.service.Get(ctx, codebaseID)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	return &resolver{remote: rem, root: r}, nil
}

func (r *remoteRootResolver) CreateOrUpdateCodebaseRemote(ctx context.Context, args resolvers.CreateOrUpdateCodebaseRemoteArgsArgs) (resolvers.RemoteResolver, error) {
	codebaseID := codebases.ID(args.Input.CodebaseID)
	cb, err := r.codebaseService.GetByID(ctx, codebaseID)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerror.Error(err)
	}

	var keyPairID *crypto.KeyPairID
	if args.Input.KeyPairID != nil {
		kpi := crypto.KeyPairID(*args.Input.KeyPairID)
		keyPairID = &kpi
	}

	rem, err := r.service.SetRemote(
		ctx,
		codebaseID,
		&service.SetRemoteInput{
			Name:              args.Input.Name,
			URL:               args.Input.Url,
			BasicAuthUsername: args.Input.BasicAuthUsername,
			BasicAuthPassword: args.Input.BasicAuthPassword,
			TrackedBranch:     args.Input.TrackedBranch,
			BrowserLinkRepo:   args.Input.BrowserLinkRepo,
			BrowserLinkBranch: args.Input.BrowserLinkBranch,
			KeyPairID:         keyPairID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add remote: %w", err)
	}

	return &resolver{remote: rem, root: r}, nil
}

package graphql

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/remote"
)

type resolver struct {
	remote *remote.Remote
	root   *remoteRootResolver
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.remote.ID)
}

func (r *resolver) Name() string {
	return r.remote.Name
}

func (r *resolver) URL() string {
	return r.remote.URL
}

func (r *resolver) TrackedBranch() string {
	return r.remote.TrackedBranch
}

func (r *resolver) BasicAuthUsername() *string {
	return r.remote.BasicAuthUsername
}

func (r *resolver) BasicAuthPassword() *string {
	return r.remote.BasicAuthPassword
}

func (r *resolver) BrowserLinkRepo() string {
	return r.remote.BrowserLinkRepo
}

func (r *resolver) BrowserLinkBranch() string {
	return r.remote.BrowserLinkBranch
}

func (r *resolver) KeyPair(ctx context.Context) (resolvers.KeyPairResolver, error) {
	if r.remote.KeyPairID == nil {
		return nil, nil
	}
	kp, err := r.root.cryptoRootResolver.InternalGetByID(ctx, *r.remote.KeyPairID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return kp, nil
}

func (r *resolver) Enabled() bool {
	return r.remote.Enabled
}

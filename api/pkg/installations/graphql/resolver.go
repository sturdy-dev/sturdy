package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"github.com/graph-gophers/graphql-go"
)

type resolver struct {
	root *rootResolver
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.root.installation.ID)
}

func (r *resolver) NeedsFirstTimeSetup(ctx context.Context) (bool, error) {
	hasOrg, err := r.root.service.HasOrganization(ctx)
	if err != nil {
		return false, gqlerrors.Error(err)
	}
	return hasOrg, nil
}

func (r *resolver) Version() string {
	return r.root.installation.Version
}

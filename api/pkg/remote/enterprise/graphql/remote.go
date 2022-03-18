package graphql

import (
	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/remote"
)

type resolver struct {
	remote *remote.Remote
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

func (r *resolver) BasicAuthUsername() string {
	return r.remote.BasicAuthUsername
}

func (r *resolver) BasicAuthPassword() string {
	return r.remote.BasicAuthPassword
}

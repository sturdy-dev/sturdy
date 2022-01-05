package graphql

import (
	"context"

	"mash/pkg/codebase/acl"
	"mash/pkg/codebase/acl/access"
	provider_acl "mash/pkg/codebase/acl/provider"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	db_user "mash/pkg/user/db"

	"github.com/graph-gophers/graphql-go"
	"github.com/tailscale/hujson"
)

type ACLRootResolver struct {
	aclProvider *provider_acl.Provider
	userRepo    db_user.Repository
}

func NewResolver(
	aclProvider *provider_acl.Provider,
	userRepo db_user.Repository,
) resolvers.ACLRootResolver {
	return &ACLRootResolver{
		aclProvider: aclProvider,
		userRepo:    userRepo,
	}
}

func (r *ACLRootResolver) CanI(ctx context.Context, args resolvers.CanIArgs) (bool, error) {
	resource := new(acl.Identity)
	resource.ParseString(args.Resource)

	if !resource.Type.IsValid() {
		return false, gqlerrors.Error(gqlerrors.ErrBadRequest, "resource", "unsupported resource type")
	}

	action := acl.Action(args.Action)
	if !action.IsValid() {
		return false, gqlerrors.Error(gqlerrors.ErrBadRequest, "action", "unsupported type")
	}

	a, err := r.aclProvider.GetByCodebaseID(ctx, string(args.CodebaseID))
	if err != nil {
		return false, gqlerrors.Error(err)
	}

	allowed, err := access.UserCan(ctx, r.userRepo, a.Policy, action, *resource)
	if err != nil {
		return false, gqlerrors.Error(err)
	}

	return allowed, nil
}

func (r *ACLRootResolver) InternalACLByCodebaseID(ctx context.Context, codebaseID graphql.ID) (resolvers.ACLResolver, error) {
	a, err := r.aclProvider.GetByCodebaseID(ctx, string(codebaseID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &aclResolver{a: a, root: r}, nil
}

func (r *ACLRootResolver) UpdateACL(ctx context.Context, args resolvers.UpdateACLArgs) (resolvers.ACLResolver, error) {
	a, err := r.aclProvider.GetByCodebaseID(ctx, string(args.Input.CodebaseID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	allowed, err := access.UserCanWriteACL(ctx, r.userRepo, a.Policy, string(a.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	if !allowed {
		return nil, gqlerrors.Error(gqlerrors.ErrForbidden)
	}

	if args.Input.Policy == nil {
		return &aclResolver{a: a, root: r}, nil
	}

	policy := acl.Policy{}
	if err := hujson.Unmarshal([]byte(*args.Input.Policy), &policy); err != nil {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "policy", "failed to decode as json")
	}

	if errs := policy.Errors(string(a.ID)); len(errs) > 0 {
		msgs := make([]string, 0, len(errs)*2)
		for k, v := range errs {
			msgs = append(msgs, k, v.Error())
		}
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, msgs...)

	}

	a.RawPolicy = *args.Input.Policy

	if err := r.aclProvider.Update(ctx, a); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &aclResolver{a: a, root: r}, nil
}

type aclResolver struct {
	a    acl.ACL
	root *ACLRootResolver
}

func (r *aclResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID)
}

func (r *aclResolver) Policy() (string, error) {
	return r.a.RawPolicy, nil
}

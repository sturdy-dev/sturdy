package access

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebases/acl"
	"getsturdy.com/api/pkg/users"
)

type userRepository interface {
	Get(users.ID) (*users.User, error)
}

type aclProvider interface {
	GetByCodebaseID(context.Context, string) (acl.ACL, error)
}

func UserCan(
	ctx context.Context,
	userRepo userRepository,
	aclPolicy acl.Policy,
	action acl.Action,
	resource acl.Identity,
) (bool, error) {
	// todo: move this into auth package
	userID, err := auth.UserID(ctx)
	if err != nil {
		return false, err
	}

	allowedByID := aclPolicy.Assert(
		acl.Identity{Type: acl.Users, ID: userID.String()},
		action,
		resource,
	)
	if allowedByID {
		return true, nil
	}

	user, err := userRepo.Get(userID)
	if err != nil {
		return false, err
	}

	allowedByEmail := aclPolicy.Assert(
		acl.Identity{Type: acl.Users, ID: user.Email},
		action,
		resource,
	)
	return allowedByEmail, nil

}

func UserCanWriteACL(
	ctx context.Context,
	userRepo userRepository,
	aclPolicy acl.Policy,
	aclID string,
) (bool, error) {
	action := acl.ActionWrite
	resource := acl.Identity{Type: acl.ACLs, ID: aclID}
	return UserCan(ctx, userRepo, aclPolicy, action, resource)
}

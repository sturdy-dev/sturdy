package graphql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"mash/pkg/auth"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/pki"
	"mash/pkg/pki/db"
)

type pkiRootResolver struct {
	repo  db.Repo
	users resolvers.UserRootResolver
}

func NewResolver(repo db.Repo, users resolvers.UserRootResolver) resolvers.PKIRootResolver {
	return &pkiRootResolver{
		repo:  repo,
		users: users,
	}
}

func (p *pkiRootResolver) AddPublicKey(ctx context.Context, args resolvers.AddPublicKeyArgs) (resolvers.UserResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	rows, err := p.repo.GetKeyByUserID(userID)
	switch {
	// If the user doesn't have any public keys, continue
	case errors.Is(err, sql.ErrNoRows):
	case err == nil:
		for _, row := range rows {
			// If the user already has this public key registered, skip adding it again...
			if row.PublicKey == args.PublicKey {
				return p.users.User(ctx)
			}
		}
		// ... otherwise continue

	default:
		// If there was another error, fail
		return nil, gqlerrors.Error(err)
	}

	upk := pki.UserPublicKey{
		PublicKey: args.PublicKey,
		UserID:    userID,
		AddedAt:   time.Now(),
	}

	if err = p.repo.Create(upk); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return p.users.User(ctx)
}

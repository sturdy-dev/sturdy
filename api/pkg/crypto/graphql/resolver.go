package graphql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/crypto"
	db_crypto "getsturdy.com/api/pkg/crypto/db"
	"getsturdy.com/api/pkg/crypto/rsa"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type rootResolver struct {
	repo db_crypto.KeyPairRepository
}

func New(repo db_crypto.KeyPairRepository) resolvers.CryptoRootResolver {
	return &rootResolver{
		repo: repo,
	}
}

func (r *rootResolver) InternalGetByID(ctx context.Context, id crypto.KeyPairID) (resolvers.KeyPairResolver, error) {
	kp, err := r.repo.Get(ctx, id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &keyPairResolver{kp: kp}, nil
}

func (r *rootResolver) GenerateKeyPair(ctx context.Context, args resolvers.GeneratePublicKeyArgs) (resolvers.KeyPairResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var public crypto.PublicKey
	var private crypto.PrivateKey

	switch args.Input.KeyPairType {
	case resolvers.KeyPairType_RSA_4096:
		public, private, err = rsa.GenerateRsaKeypair()
	default:
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "unsupported keyPairType")
	}
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	keyPair := crypto.KeyPair{
		ID:         crypto.KeyPairID(uuid.NewString()),
		PublicKey:  public,
		PrivateKey: private,
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
		CreatedBy:  userID,
	}

	if err := r.repo.Create(ctx, keyPair); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &keyPairResolver{kp: &keyPair}, nil
}

type keyPairResolver struct {
	kp *crypto.KeyPair
}

func (k *keyPairResolver) ID() graphql.ID {
	return graphql.ID(k.kp.ID)
}

func (k *keyPairResolver) PublicKey() string {
	return string(k.kp.PublicKey)
}

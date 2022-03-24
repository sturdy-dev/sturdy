package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/crypto"
)

type CryptoRootResolver interface {
	InternalGetByID(ctx context.Context, id crypto.KeyPairID) (KeyPairResolver, error)

	// Mutations
	GenerateKeyPair(ctx context.Context, args GeneratePublicKeyArgs) (KeyPairResolver, error)
}

type GeneratePublicKeyArgs struct {
	Input GeneratePublicKeyInput
}

type GeneratePublicKeyInput struct {
	KeyPairType KeyPairType
}

type KeyPairType string

const (
	KeyPairType_RSA_4096 KeyPairType = "RSA_4096"
)

type KeyPairResolver interface {
	ID() graphql.ID
	PublicKey() string
}

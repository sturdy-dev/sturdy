package resolvers

import "context"

type PKIRootResolver interface {
	// Mutation
	AddPublicKey(context.Context, AddPublicKeyArgs) (UserResolver, error)
}

type AddPublicKeyArgs struct {
	PublicKey string
}

package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type ServiceTokensRootResolver interface {
	CreateServiceToken(context.Context, CreateServiceTokenArgs) (ServiceTokenResovler, error)
}

type CreateServiceTokenArgs struct {
	Input CreateServiceTokenInput
}

type CreateServiceTokenInput struct {
	Name            string
	ShortCodebaseID string
}

type ServiceTokenResovler interface {
	ID() graphql.ID
	Name() string
	CreatedAt() int32
	LastUsedAt() *int32

	Token() *string
}

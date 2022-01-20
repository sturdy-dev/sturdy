package resolvers

import (
	"context"
	"getsturdy.com/api/pkg/integrations"

	"github.com/graph-gophers/graphql-go"
)

type IntegrationRootResolver interface {
	// mutations
	TriggerInstantIntegration(ctx context.Context, args TriggerInstantIntegrationArgs) ([]StatusResolver, error)
	DeleteIntegration(ctx context.Context, args DeleteIntegrationArgs) (IntegrationResolver, error)

	// internal
	InternalIntegrationProvider(*integrations.Integration) IntegrationResolver
	InternalIntegrationsByCodebaseID(context.Context, string) ([]IntegrationResolver, error)
	InternalIntegrationByID(context.Context, string) (IntegrationResolver, error)
}

type TriggerInstantIntegrationArgs struct {
	Input TriggerInstantIntegrationInput
}

type TriggerInstantIntegrationInput struct {
	ChangeID  graphql.ID
	Providers *[]InstantIntegrationProviderType
}

type DeleteIntegrationArgs struct {
	Input DeleteIntegrationInput
}

type DeleteIntegrationInput struct {
	ID graphql.ID
}

type commonIntegrationResolver interface {
	ID() graphql.ID
	CodebaseID() graphql.ID
	Provider() (InstantIntegrationProviderType, error)
	CreatedAt() int32
	UpdatedAt() *int32
	DeletedAt() *int32
}

type BuildkiteIntegration interface {
	commonIntegrationResolver

	Configuration(context.Context) (BuildkiteConfigurationResolver, error)
}

type IntegrationResolver interface {
	ToBuildkiteIntegration() (BuildkiteIntegration, bool)

	commonIntegrationResolver
}

type InstantIntegrationProviderType string

const (
	InstantIntegrationProviderUndefined InstantIntegrationProviderType = ""
	InstantIntegrationProviderBuildkite InstantIntegrationProviderType = "Buildkite"
)

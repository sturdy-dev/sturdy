package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type BuildkiteInstantIntegrationRootResolver interface {
	// mutations
	CreateOrUpdateBuildkiteIntegration(context.Context, CreateOrUpdateBuildkiteIntegrationArgs) (IntegrationResolver, error)

	// internal
	InternalBuildkiteConfigurationByIntegrationID(context.Context, string) (BuildkiteConfigurationResolver, error)
}

type CreateOrUpdateBuildkiteIntegrationArgs struct {
	Input CreateOrUpdateBuildkiteIntegrationInput
}

type CreateOrUpdateBuildkiteIntegrationInput struct {
	CodebaseID       graphql.ID
	IntegrationID    *graphql.ID
	OrganizationName string
	PipelineName     string
	APIToken         string
	WebhookSecret    string
}

type BuildkiteConfigurationResolver interface {
	ID() graphql.ID
	OrganizationName() string
	PipelineName() string
	APIToken() string
	WebhookSecret() string
}

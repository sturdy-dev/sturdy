package graphql

import (
	"mash/pkg/integrations/buildkite"

	"github.com/graph-gophers/graphql-go"
)

type buildkiteConfigurationResover struct {
	buildkiteConfig *buildkite.Config
}

func (r *buildkiteConfigurationResover) ID() graphql.ID {
	return graphql.ID(r.buildkiteConfig.ID)
}

func (r *buildkiteConfigurationResover) OrganizationName() string {
	return r.buildkiteConfig.OrganizationName
}

func (r *buildkiteConfigurationResover) PipelineName() string {
	return r.buildkiteConfig.PipelineName
}

func (r *buildkiteConfigurationResover) APIToken() string {
	return r.buildkiteConfig.APIToken
}

func (r *buildkiteConfigurationResover) WebhookSecret() string {
	return r.buildkiteConfig.WebhookSecret
}

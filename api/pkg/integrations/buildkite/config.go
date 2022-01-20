package buildkite

import "time"

type Config struct {
	ID         string `db:"id"`          // TODO: Remove?
	CodebaseID string `db:"codebase_id"` // TODO: Remove?

	IntegrationID string `db:"integration_id"`

	OrganizationName string    `db:"organization_name"`
	PipelineName     string    `db:"pipeline_name"`
	APIToken         string    `db:"api_token"`
	WebhookSecret    string    `db:"webhook_secret"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

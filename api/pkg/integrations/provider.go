package integrations

import (
	"context"
)

type Provider interface {
	CreateBuild(ctx context.Context, codebaseID, ciCommitId, title string) (*Build, error)
}

type ProviderType string

const (
	ProviderTypeUndefined ProviderType = ""
	ProviderTypeBuildkite ProviderType = "buildkite"
)

type Build struct {
	Name        string
	Description string
	URL         string
}

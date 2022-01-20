//go:build enterprise || cloud
// +build enterprise cloud

package resolvers

import (
	"context"
)

type CloudRootResolver interface {
	SelfHostedPing(context.Context, SelfHostedPingArgs) SelfHostedInstallationStatus
}

type SelfHostedPingArgs struct {
	Input SelfHostedPingInput
}

type SelfHostedPingInput struct {
	Version string
}

type SelfHostedInstallationStatus interface {
	Ok() bool
}

package posthog

import (
	"flag"
	"fmt"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/analytics/disabled"

	"github.com/posthog/posthog-go"
)

var (
	analyticsEnabled = flag.Bool("analytics.enabled", true, "Enable analytics")
	posthogAPIToken  = flag.String("analytics.posthog.api-token", "", "Posthog API token")
)

func NewClient() (analytics.Client, error) {
	if !*analyticsEnabled {
		return disabled.NewClient(), nil
	}
	if *posthogAPIToken == "" {
		return nil, fmt.Errorf("--analytics.posthog.api-token is required")
	}
	return analytics.New(posthog.New(*posthogAPIToken)), nil
}

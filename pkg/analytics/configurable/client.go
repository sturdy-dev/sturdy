package configurable

import (
	"flag"
	"fmt"

	"mash/pkg/analytics"
	"mash/pkg/analytics/disabled"
	"mash/pkg/analytics/posthog"
)

var (
	analyticsEnabled = flag.Bool("analytics.enabled", true, "Enable analytics")
	analyticsType    = flag.String("analytics.type", "", "Analytics type, must be one of: posthog")

	posthogAPIToken = flag.String("analytics.posthog.api-token", "", "Posthog API token (required if analytics type is posthog)")
)

func NewClient() (analytics.Client, error) {
	if !*analyticsEnabled {
		return disabled.NewClient(), nil
	}

	switch *analyticsType {
	case "posthog":
		if *posthogAPIToken == "" {
			return nil, fmt.Errorf("--analytics.posthog.api-token is required")
		}
		return posthog.NewClient(*posthogAPIToken), nil
	case "":
		return nil, fmt.Errorf("--analytics.type is required")
	default:
		return nil, fmt.Errorf("unknown analytics type: %s", *analyticsType)
	}
}

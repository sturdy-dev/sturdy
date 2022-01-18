package proxy

import (
	"flag"

	"mash/pkg/analytics"
	"mash/pkg/analytics/disabled"
	"mash/pkg/installations"

	"github.com/posthog/posthog-go"
)

var (
	analyticsEnabled = flag.Bool("analytics.enabled", true, "Enable analytics")
)

type client struct {
	posthog.Client

	installation *installations.Installation
}

func (c *client) Enqueue(event analytics.Message) error {
	switch e := event.(type) {
	case *analytics.Capture:
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type)
	case analytics.Capture:
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type)
	case analytics.Identify:
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type)
	case *analytics.Identify:
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type)
	}

	return c.Client.Enqueue(event)
}

func NewClient(installation *installations.Installation) (analytics.Client, error) {
	if !*analyticsEnabled {
		return disabled.NewClient(), nil
	}

	// api token is intentionally not set here, as it is not needed for the proxy client
	posthogClient, err := posthog.NewWithConfig("", posthog.Config{
		Endpoint: "https://api.getsturdy.com/v3/analytics",
	})
	if err != nil {
		return nil, err
	}

	return analytics.New(&client{
		Client:       posthogClient,
		installation: installation,
	}), nil
}

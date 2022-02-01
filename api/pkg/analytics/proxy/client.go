package proxy

import (
	"flag"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/analytics/disabled"
	"getsturdy.com/api/pkg/installations"

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
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type.String())
		e.Properties.Set("version", c.installation.Version)
	case analytics.Capture:
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type.String())
		e.Properties.Set("version", c.installation.Version)
	case analytics.Identify:
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type.String())
		e.Properties.Set("version", c.installation.Version)
	case *analytics.Identify:
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", c.installation.ID)
		e.Properties.Set("installation_type", c.installation.Type.String())
		e.Properties.Set("version", c.installation.Version)
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

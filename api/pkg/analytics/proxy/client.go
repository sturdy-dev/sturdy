package proxy

import (
	"flag"
	"fmt"

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

	installation installations.GetInstallationFunc
}

func (c *client) Enqueue(event analytics.Message) error {
	ins, err := c.installation()
	if err != nil {
		return fmt.Errorf("could not get current installation: %w", err)
	}

	switch e := event.(type) {
	case *analytics.Capture:
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", ins.ID)
		e.Properties.Set("installation_type", ins.Type.String())
		e.Properties.Set("version", ins.Version)
	case analytics.Capture:
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", ins.ID)
		e.Properties.Set("installation_type", ins.Type.String())
		e.Properties.Set("version", ins.Version)
	case analytics.Identify:
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", ins.ID)
		e.Properties.Set("installation_type", ins.Type.String())
		e.Properties.Set("version", ins.Version)
	case *analytics.Identify:
		if e.Properties == nil {
			e.Properties = make(map[string]interface{})
		}
		e.Properties.Set("installation_id", ins.ID)
		e.Properties.Set("installation_type", ins.Type.String())
		e.Properties.Set("version", ins.Version)
	}

	return c.Client.Enqueue(event)
}

func NewClient(installation installations.GetInstallationFunc) (analytics.Client, error) {
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

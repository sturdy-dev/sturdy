package proxy

import (
	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/analytics/disabled"
	"getsturdy.com/api/pkg/installations"

	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

type Configuration struct {
	Disable bool `long:"disable" description:"Disable analytics"`
}

type client struct {
	posthog.Client

	logger       *zap.Logger
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

	c.identifyInstallation()

	return c.Client.Enqueue(event)
}

func (c *client) identifyInstallation() {
	if err := c.Enqueue(analytics.GroupIdentify{
		Type: "installation_id", // this should match other event's property key
		Key:  c.installation.ID,
		Properties: map[string]interface{}{
			"type":            c.installation.Type.String(),
			"version":         c.installation.Version,
			"license_key_set": c.installation.LicenseKey != nil,
		},
	}); err != nil {
		c.logger.Error("failed to identify installation", zap.Error(err))
	}
}

func NewClient(
	cfg *Configuration,
	installation *installations.Installation,
	logger *zap.Logger,
) (analytics.Client, error) {
	if cfg.Disable {
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
		logger:       logger,
		installation: installation,
	}), nil
}

package proxy

import (
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

func (c *client) setProperties(properties map[string]any) map[string]any {
	if properties == nil {
		properties = make(map[string]any, 3)
	}
	properties["installation_id"] = c.installation.ID
	properties["installation_type"] = c.installation.Type.String()
	properties["version"] = c.installation.Version
	return properties
}

func (c *client) setGroups(groups map[string]any) map[string]any {
	if groups == nil {
		groups = make(map[string]any, 1)
	}
	groups["installation"] = c.installation.ID
	return groups
}

func (c *client) Enqueue(event posthog.Message) error {
	switch e := event.(type) {
	case *posthog.Capture:
		e.Properties = c.setProperties(e.Properties)
		e.Groups = c.setGroups(e.Groups)
	case posthog.Capture:
		e.Properties = c.setProperties(e.Properties)
		e.Groups = c.setGroups(e.Groups)
	case posthog.Identify:
		e.Properties = c.setProperties(e.Properties)
	case *posthog.Identify:
		e.Properties = c.setProperties(e.Properties)
	case *posthog.GroupIdentify:
		e.Properties = c.setProperties(e.Properties)
	case posthog.GroupIdentify:
		e.Properties = c.setProperties(e.Properties)
	}

	c.identifyInstallation()

	return c.Client.Enqueue(event)
}

func (c *client) identifyInstallation() {
	if err := c.Client.Enqueue(posthog.GroupIdentify{
		Type: "installation", // this should match other event's property key
		Key:  c.installation.ID,
		Properties: map[string]any{
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
) (posthog.Client, error) {
	if cfg.Disable {
		return disabled.NewClient(logger), nil
	}

	// api token is intentionally not set here, as it is not needed for the proxy client
	posthogClient, err := posthog.NewWithConfig("", posthog.Config{
		Endpoint: "https://api.getsturdy.com/v3/analytics",
	})
	if err != nil {
		return nil, err
	}

	c := &client{
		Client:       posthogClient,
		logger:       logger,
		installation: installation,
	}
	c.identifyInstallation()

	return c, nil
}

package proxy

import (
	"mash/pkg/analytics"

	"github.com/posthog/posthog-go"
)

type Client struct{}

// TODO: set something to identify the source of the events
func NewClient() (analytics.Client, error) {
	// api token is intentionally not set here, as it is not needed for the proxy client
	posthogClient, err := posthog.NewWithConfig("", posthog.Config{
		Endpoint: "https://api.getsturdy.com/v3/analytics",
	})
	if err != nil {
		return nil, err
	}
	return analytics.New(posthogClient), nil
}

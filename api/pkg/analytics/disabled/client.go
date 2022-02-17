package disabled

import (
	"github.com/posthog/posthog-go"
)

type Client struct{}

func NewClient() posthog.Client {
	return &Client{}
}

func (*Client) Enqueue(posthog.Message) error { return nil }

func (*Client) Close() error { return nil }

func (*Client) IsFeatureEnabled(string, string, bool) (bool, error) { return false, nil }

func (*Client) ReloadFeatureFlags() error { return nil }

func (*Client) GetFeatureFlags() ([]posthog.FeatureFlag, error) { return nil, nil }

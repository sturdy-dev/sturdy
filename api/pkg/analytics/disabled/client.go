package disabled

import (
	"getsturdy.com/api/pkg/analytics"
	"github.com/posthog/posthog-go"
)

type Client struct{}

func NewClient() analytics.Client {
	return analytics.New(&Client{})
}

func (*Client) Enqueue(analytics.Message) error { return nil }

func (*Client) Close() error { return nil }

func (*Client) IsFeatureEnabled(string, string, bool) (bool, error) { return false, nil }

func (*Client) ReloadFeatureFlags() error { return nil }

func (*Client) GetFeatureFlags() ([]posthog.FeatureFlag, error) { return nil, nil }

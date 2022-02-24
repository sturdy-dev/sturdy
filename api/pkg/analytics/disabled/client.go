package disabled

import (
	"github.com/posthog/posthog-go"

	"go.uber.org/zap"
)

type Client struct{}

func NewClient(logger *zap.Logger) posthog.Client {
	logger.Info("analytics disabled")
	return &Client{}
}

func (*Client) Enqueue(posthog.Message) error { return nil }

func (*Client) Close() error { return nil }

func (*Client) IsFeatureEnabled(string, string, bool) (bool, error) { return false, nil }

func (*Client) ReloadFeatureFlags() error { return nil }

func (*Client) GetFeatureFlags() ([]posthog.FeatureFlag, error) { return nil, nil }

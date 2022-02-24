package posthog

import (
	"getsturdy.com/api/pkg/analytics/disabled"
	"getsturdy.com/api/pkg/analytics/proxy"

	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

type Configuration struct {
	proxy.Configuration
	Posthog *postHogConfiguration `flags-group:"posthog" namespace:"posthog" required:"true"`
}

type postHogConfiguration struct {
	ApiToken string `long:"api-token" description:"PostHog API token"`
}

func NewClient(logger *zap.Logger, cfg *Configuration) (posthog.Client, error) {
	if cfg.Disable {
		return disabled.NewClient(logger), nil
	}
	return posthog.New(cfg.Posthog.ApiToken), nil
}

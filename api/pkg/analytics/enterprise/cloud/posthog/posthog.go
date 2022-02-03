package posthog

import (
	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/analytics/disabled"
	"getsturdy.com/api/pkg/analytics/proxy"

	"github.com/posthog/posthog-go"
)

type Configuration struct {
	proxy.Configuration
	Posthog *postHogConfiguration `flags-group:"posthog" namespace:"posthog" required:"true"`
}

type postHogConfiguration struct {
	ApiToken string `long:"api-token" description:"PostHog API token"`
}

func NewClient(cfg *Configuration) (analytics.Client, error) {
	if cfg.Disable {
		return disabled.NewClient(), nil
	}
	return analytics.New(posthog.New(cfg.Posthog.ApiToken)), nil
}

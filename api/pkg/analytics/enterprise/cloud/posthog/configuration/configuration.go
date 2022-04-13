package configuration

import (
	proxy "getsturdy.com/api/pkg/analytics/proxy/configuration"
)

type Configuration struct {
	proxy.Configuration
	Posthog *postHogConfiguration `flags-group:"posthog" namespace:"posthog" required:"true"`
}

type postHogConfiguration struct {
	ApiToken string `long:"api-token" description:"PostHog API token"`
}

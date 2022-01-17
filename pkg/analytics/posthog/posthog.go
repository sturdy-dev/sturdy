package posthog

import (
	"mash/pkg/analytics"

	"github.com/posthog/posthog-go"
)

func NewClient(apiToken string) analytics.Client {
	return analytics.New(posthog.New(apiToken))
}

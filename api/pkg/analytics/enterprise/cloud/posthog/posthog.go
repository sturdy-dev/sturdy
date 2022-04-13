package posthog

import (
	"getsturdy.com/api/pkg/analytics/disabled"
	"getsturdy.com/api/pkg/analytics/enterprise/cloud/posthog/configuration"

	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

func NewClient(logger *zap.Logger, cfg *configuration.Configuration) (posthog.Client, error) {
	if cfg.Disable {
		return disabled.NewClient(logger), nil
	}
	return posthog.New(cfg.Posthog.ApiToken), nil
}

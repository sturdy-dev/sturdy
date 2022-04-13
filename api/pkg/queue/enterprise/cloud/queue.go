package cloud

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/enterprise/cloud/configuration"
)

func New(awsSession *session.Session, logger *zap.Logger, cfg *configuration.Configuration) (queue.Queue, error) {
	if cfg.Hostname == "" {
		defaultHostname, err := os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("failed to get hostname: %w", err)
		}
		cfg.Hostname = defaultHostname
	}
	if cfg.Local {
		return queue.NewInMemory(logger), nil
	}
	return NewSQS(logger, awsSession, cfg.Hostname, cfg.Prefix)
}

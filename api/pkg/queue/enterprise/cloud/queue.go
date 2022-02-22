package cloud

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/queue"
)

type Configuration struct {
	Local    bool   `long:"local" description:"Use in-memory queue instead of SQS"`
	Hostname string `long:"hostname" description:"Hostname of the queue"`
	Prefix   string `long:"prefix" description:"Prefix for queue names"`
}

func New(awsSession *session.Session, logger *zap.Logger, cfg *Configuration) (queue.Queue, error) {
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

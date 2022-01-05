package emails

import (
	"context"

	"go.uber.org/zap"
)

var _ Sender = &logsClient{}

type logsClient struct {
	logger *zap.Logger
}

func NewLogs(logger *zap.Logger) *logsClient {
	return &logsClient{
		logger: logger.Named("email sender"),
	}
}

func (s *logsClient) Send(ctx context.Context, msg *Email) error {
	s.logger.Info(
		"would have sent an email",
		zap.Any("message", msg),
	)
	return nil
}

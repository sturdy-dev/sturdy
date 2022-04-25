package cloud

import (
	"context"

	"getsturdy.com/api/pkg/emails"
	"go.uber.org/zap"
)

var _ emails.Sender = &logsClient{}

type logsClient struct {
	logger *zap.Logger
}

func NewLogs(logger *zap.Logger) *logsClient {
	return &logsClient{
		logger: logger.Named("email sender"),
	}
}

func (s *logsClient) Send(ctx context.Context, msg *emails.Email) error {
	s.logger.Warn(
		"would have sent an email",
		zap.Any("message", msg),
	)
	return nil
}

package cloud

import (
	"getsturdy.com/api/pkg/emails"
	"getsturdy.com/api/pkg/emails/enterprise/cloud/configuration"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

func New(cfg *configuration.Configuration, awsSession *session.Session, logger *zap.Logger) emails.Sender {
	// Not enabled, write emails to log
	if !cfg.Enable {
		return NewLogs(logger)
	}

	if cfg.Provider == "postmark" {
		return NewPostmarkClient(cfg.Postmark)
	}

	if cfg.Provider == "ses" || cfg.Provider == "" {
		return NewSES(awsSession)
	}

	logger.Fatal("No valid email configuration found", zap.String("provider", cfg.Provider))
	panic("No valid email configuration found")
}

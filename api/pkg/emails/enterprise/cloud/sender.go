package cloud

import (
	"getsturdy.com/api/pkg/emails"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type Configuration struct {
	Enable   bool                   `long:"enable" description:"Use AWS SES to send emails"`
	Provider string                 `long:"provider" description:"Which email provider to use" default:"ses"`
	Postmark *PostmarkConfiguration `flags-group:"postmark" namespace:"postmark"`
}

func New(cfg *Configuration, awsSession *session.Session, logger *zap.Logger) emails.Sender {
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

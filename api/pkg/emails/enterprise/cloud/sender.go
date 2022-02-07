package cloud

import (
	"getsturdy.com/api/pkg/emails"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type Configuration struct {
	Enable bool `long:"enable" description:"Use AWS SES to send emails"`
}

func New(cfg *Configuration, awsSession *session.Session, logger *zap.Logger) emails.Sender {
	if cfg.Enable {
		return NewSES(awsSession)
	}
	return NewLogs(logger)
}

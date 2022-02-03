package emails

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type Email struct {
	To      string
	Subject string
	Html    string
}

type Sender interface {
	Send(context.Context, *Email) error
}

type Configuration struct {
	Enable bool `long:"enable" description:"Use AWS SES to send emails"`
}

func New(cfg *Configuration, awsSession *session.Session, logger *zap.Logger) Sender {
	if cfg.Enable {
		return NewSES(awsSession)
	}
	return NewLogs(logger)
}

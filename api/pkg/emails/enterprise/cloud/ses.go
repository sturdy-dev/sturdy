package cloud

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/emails"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var _ emails.Sender = &sesClient{}

type sesClient struct {
	sesClient *ses.SES
}

func NewSES(s *session.Session) *sesClient {
	return &sesClient{
		sesClient: ses.New(s),
	}
}

func (s *sesClient) Send(ctx context.Context, msg *emails.Email) error {
	if _, err := s.sesClient.SendEmailWithContext(ctx, &ses.SendEmailInput{
		Destination: &ses.Destination{ToAddresses: []*string{aws.String(msg.To)}},
		Source:      aws.String("Sturdy <no-reply@getsturdy.com>"),
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    &msg.Subject,
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(msg.Html),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to send email to ses: %w", err)
	}
	return nil
}

package cloud

import (
	"context"
	"fmt"

	"github.com/keighl/postmark"

	"getsturdy.com/api/pkg/emails"
	"getsturdy.com/api/pkg/emails/enterprise/cloud/configuration"
)

var _ emails.Sender = &postmarkSender{}

type postmarkSender struct {
	postmarkClient *postmark.Client
}

func NewPostmarkClient(config *configuration.PostmarkConfiguration) *postmarkSender {
	return &postmarkSender{
		postmarkClient: postmark.NewClient(config.ServerToken, ""),
	}
}

func (s *postmarkSender) Send(ctx context.Context, msg *emails.Email) error {
	email := postmark.Email{
		From:     "support@getsturdy.com",
		To:       msg.To,
		Subject:  msg.Subject,
		HtmlBody: msg.Html,
	}

	_, err := s.postmarkClient.SendEmail(email)
	if err != nil {
		return fmt.Errorf("failed to send email via postmark: %w", err)
	}

	return nil
}

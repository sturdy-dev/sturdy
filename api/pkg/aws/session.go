package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func New(cfg *Configuration) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
}

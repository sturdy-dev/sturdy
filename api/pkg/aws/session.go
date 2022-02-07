package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Configuration struct {
	Region string `long:"region" description:"AWS region to use" default:"eu-north-1"`
}

func New(cfg *Configuration) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
}

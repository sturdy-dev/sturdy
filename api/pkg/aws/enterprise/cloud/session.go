package cloud

import (
	"getsturdy.com/api/pkg/aws/enterprise/cloud/configuration"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func New(cfg *configuration.Configuration) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
}

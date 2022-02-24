package uploader

import (
	"context"
	"fmt"
	"io"

	"getsturdy.com/api/pkg/users/avatars"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var _ Uploader = &S3{}

type S3 struct {
	uploader *s3manager.Uploader
}

func NewS3(s *session.Session) *S3 {
	return &S3{
		uploader: s3manager.NewUploader(s),
	}
}

func (s *S3) Upload(ctx context.Context, key string, file io.Reader) (*avatars.Avatar, error) {
	if _, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String("usercontent.getsturdy.com"),
		Key:    aws.String(key),
		Body:   file,
	}); err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}
	return &avatars.Avatar{
		URL: fmt.Sprintf("https://usercontent.getsturdy.com/%s", key),
	}, nil
}

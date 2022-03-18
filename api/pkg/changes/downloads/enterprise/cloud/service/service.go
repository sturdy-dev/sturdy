package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"

	service_change "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	*service_change.Service

	maybeBucketName  *string
	executorProvider executor.Provider
	awsSession       *session.Session
	logger           *zap.Logger
}

type Configuration struct {
	ExportBucketName string `long:"export-bucket-name" description:"The name of the S3 bucket to export change archives to"`
}

func New(
	service *service_change.Service,
	cfg *Configuration,
	executorProvider executor.Provider,
	awsSession *session.Session,
	logger *zap.Logger,
) *Service {
	if cfg == nil {
		cfg = &Configuration{}
	}
	return &Service{
		Service:          service,
		executorProvider: executorProvider,
		awsSession:       awsSession,
		maybeBucketName:  &cfg.ExportBucketName,
		logger:           logger,
	}
}

type ArchiveFormat int

const (
	ArchiveFormatUnknown ArchiveFormat = iota
	ArchiveFormatZip
	ArchiveFormatTarGz
)

func (svc *Service) CreateArchive(ctx context.Context, allower *unidiff.Allower, codebaseID codebases.ID, commitID string, format ArchiveFormat) (string, error) {
	if svc.maybeBucketName == nil || len(*svc.maybeBucketName) == 0 {
		return "", fmt.Errorf("--export-bucket-name is not defined")
	}

	var archiveFileExt string
	switch format {
	case ArchiveFormatTarGz:
		archiveFileExt = ".tar.gz"
	case ArchiveFormatZip:
		archiveFileExt = ".zip"
	default:
		return "", fmt.Errorf("unexpected archive format")
	}

	archiveFilePath := fmt.Sprintf("%s/%s%s", codebaseID, commitID, archiveFileExt)
	if err := svc.executorProvider.New().Write(func(repo vcs.RepoWriter) error {
		if err := repo.CreateNewBranchAt("archive", commitID); err != nil {
			return fmt.Errorf("failed to create archive branch: %w", err)
		}
		if err := repo.CheckoutBranchWithForce("archive"); err != nil {
			return fmt.Errorf("failed to checkout archive: %w", err)
		}
		return repo.LargeFilesPull()
	}).Read(func(repo vcs.RepoReader) error {
		reader, writer := io.Pipe()
		var eg errgroup.Group
		switch format {
		case ArchiveFormatTarGz:
			eg.Go(svc.walkTarGz(writer, repo.Path(), allower))
		case ArchiveFormatZip:
			eg.Go(svc.walkZip(writer, repo.Path(), allower))
		default:
			return fmt.Errorf("unexpected archive format")
		}
		eg.Go(func() error {
			// Write to S3
			uploader := s3manager.NewUploader(svc.awsSession)
			_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
				Body:   reader,
				Bucket: svc.maybeBucketName,
				Key:    aws.String(archiveFilePath),
			})
			if err != nil {
				return fmt.Errorf("failed to upload: %w", err)
			}
			return nil
		})
		return eg.Wait()
	}).ExecTemporaryView(codebaseID, "createArchive"); err != nil {
		return "", fmt.Errorf("executor failed: %w", err)
	}
	// Get a pre signed URL
	client := s3.New(svc.awsSession)
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: svc.maybeBucketName,
		Key:    aws.String(archiveFilePath),
	})
	presignedURL, err := req.Presign(time.Minute * 10)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL: %w", err)
	}

	return presignedURL, nil
}

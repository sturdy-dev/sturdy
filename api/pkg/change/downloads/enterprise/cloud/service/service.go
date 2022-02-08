package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"

	service_change "getsturdy.com/api/pkg/change/service"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
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

func (svc *Service) CreateArchive(ctx context.Context, allower *unidiff.Allower, codebaseID, commitID string, format ArchiveFormat) (string, error) {
	if svc.maybeBucketName == nil || len(*svc.maybeBucketName) == 0 {
		return "", fmt.Errorf("--export-bucket-name is not defined")
	}

	viewID := fmt.Sprintf("create-archive-%s", uuid.NewString())

	var viewPath string

	if err := svc.executorProvider.New().
		AllowRebasingState(). // Allowed because the view does not exist yet
		Schedule(func(repoProvider provider.RepoProvider) error {
			viewPath = repoProvider.ViewPath(codebaseID, viewID)
			trunkPath := repoProvider.TrunkPath(codebaseID)

			if _, err := vcs.CloneRepoShared(trunkPath, viewPath); err != nil {
				return fmt.Errorf("failed to create checkout: %w", err)
			}

			// open with LFS support
			repo, err := repoProvider.ViewRepo(codebaseID, viewID)
			if err != nil {
				return fmt.Errorf("failed to open view: %w", err)
			}

			if err := repo.CreateNewBranchAt("archive", commitID); err != nil {
				return fmt.Errorf("failed to create archive branch: %w", err)
			}

			if err := repo.CheckoutBranchWithForce("archive"); err != nil {
				return fmt.Errorf("failed to checkout archive: %w", err)
			}

			return repo.LargeFilesPull()
		}).ExecView(codebaseID, viewID, "createArchive"); err != nil {
		return "", fmt.Errorf("executor failed: %w", err)
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

	reader, writer := io.Pipe()

	var eg errgroup.Group

	switch format {
	case ArchiveFormatTarGz:
		eg.Go(svc.walkTarGz(writer, viewPath, allower))
	case ArchiveFormatZip:
		eg.Go(svc.walkZip(writer, viewPath, allower))
	default:
		return "", fmt.Errorf("unexpected archive format")
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

	if err := eg.Wait(); err != nil {
		return "", err
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

	// Disk cleanup
	if err := os.RemoveAll(viewPath); err != nil {
		svc.logger.Error("failed to cleanup archive", zap.Error(err))
		// don't fail
	}

	return presignedURL, nil
}

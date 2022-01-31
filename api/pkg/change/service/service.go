package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	acl_provider "getsturdy.com/api/pkg/codebase/acl/provider"
	"getsturdy.com/api/pkg/unidiff"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/workspace"
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
	executorProvider executor.Provider
	awsSession       *session.Session
	aclProvider      *acl_provider.Provider
	userRepo         db_user.Repository
	changeRepo       db_change.Repository
	commitChangeRepo db_change.CommitRepository
	maybeBucketName  *string
}

type ExportBucketName *string

func New(
	executorProvider executor.Provider,
	awsSession *session.Session,
	aclProvider *acl_provider.Provider,
	userRepo db_user.Repository,
	changeRepo db_change.Repository,
	commitChangeRepo db_change.CommitRepository,
	maybeBucketName ExportBucketName,
) *Service {
	return &Service{
		executorProvider: executorProvider,
		awsSession:       awsSession,
		aclProvider:      aclProvider,
		userRepo:         userRepo,
		changeRepo:       changeRepo,
		commitChangeRepo: commitChangeRepo,
		maybeBucketName:  maybeBucketName,
	}
}

func (svc *Service) ListChangeCommits(ctx context.Context, ids ...change.ID) ([]*change.ChangeCommit, error) {
	return svc.commitChangeRepo.ListByChangeIDs(ctx, ids...)
}

func (svc *Service) GetChangeCommitByCommitIDAndCodebaseID(ctx context.Context, commitID, codebaseID string) (*change.ChangeCommit, error) {
	changeCommit, err := svc.commitChangeRepo.GetByCommitID(commitID, codebaseID)
	if err != nil {
		return nil, err
	}
	return &changeCommit, nil
}

func (svc *Service) GetChangeByID(ctx context.Context, id change.ID) (*change.Change, error) {
	ch, err := svc.changeRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (svc *Service) GetChangeCommitOnTrunkByChangeID(ctx context.Context, id change.ID) (*change.ChangeCommit, error) {
	ch, err := svc.commitChangeRepo.GetByChangeIDOnTrunk(id)
	if err != nil {
		return nil, err
	}
	return &ch, nil
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
			repo, err := vcs.CloneRepoShared(trunkPath, viewPath)
			if err != nil {
				return fmt.Errorf("failed to create checkout: %w", err)
			}

			// open with LFS support
			repo, err = repoProvider.ViewRepo(codebaseID, viewID)
			if err != nil {
				return fmt.Errorf("failed to open view: %w", err)
			}

			if err := repo.CreateNewBranchAt("archive", commitID); err != nil {
				return fmt.Errorf("failed to create archive branch: %w", err)
			}

			if err := repo.CheckoutBranchWithForce("archive"); err != nil {
				return fmt.Errorf("failed to checkout archive: %w", err)
			}

			if err := repo.LargeFilesPull(); err != nil {
				// ignore err
			}

			return nil
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

	return presignedURL, nil
}

func (s *Service) Create(ctx context.Context, ws *workspace.Workspace, commitID, msg string) (*change.Change, error) {
	changeID := change.ID(uuid.NewString())
	t := time.Now()
	changeChange := change.Change{
		ID:                 changeID,
		CodebaseID:         ws.CodebaseID,
		Title:              &msg,
		UpdatedDescription: ws.DraftDescription,
		UserID:             &ws.UserID,
		CreatedAt:          &t,
	}
	if err := s.changeRepo.Insert(changeChange); err != nil {
		return nil, fmt.Errorf("failed to insert change: %w", err)
	}

	if err := s.commitChangeRepo.Insert(change.ChangeCommit{
		ChangeID:   changeID,
		CommitID:   commitID,
		CodebaseID: ws.CodebaseID,
		Trunk:      true,
	}); err != nil {
		return nil, fmt.Errorf("failed to insert change commit: %w", err)
	}

	return &changeChange, nil
}

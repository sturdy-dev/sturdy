package service

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
)

func (svc *Service) CreateBuild(ctx context.Context, codebaseID codebases.ID, snapshotCommitSha, branchName string) error {
	return nil
}

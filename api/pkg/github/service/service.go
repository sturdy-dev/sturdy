package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
)

type Service interface {
	CreateBuild(ctx context.Context, codebaseID codebases.ID, snapshotCommitSha, branchName string) error
}

type svc struct{}

func (s svc) CreateBuild(ctx context.Context, codebaseID codebases.ID, snapshotCommitSha, branchName string) error {
	return fmt.Errorf("CreateBuild is not implemented in this version of Sturdy")
}

func New() Service {
	return &svc{}
}

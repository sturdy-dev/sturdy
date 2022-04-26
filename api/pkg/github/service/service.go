package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
)

type Service interface {
	CreateBuild(ctx context.Context, codebaseID codebases.ID, snapshotCommitSha, title, branchName string) (*Build, error)
}

type Build struct {
	Name        string
	Description *string
	URL         string
}

type svc struct{}

func (s svc) CreateBuild(ctx context.Context, codebaseID codebases.ID, snapshotCommitSha, title, branchName string) (*Build, error) {
	return nil, fmt.Errorf("CreateBuild is not implemented in this version of Sturdy")
}

func New() Service {
	return &svc{}
}

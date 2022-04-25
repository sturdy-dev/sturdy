package service

import (
	"context"
	"fmt"
)

type Service interface {
	CreateBuild(ctx context.Context, integrationID, ciCommitId, title string) (*Build, error)
}

type Build struct {
	Name        string
	Description *string
	URL         string
}

type svc struct{}

func (s svc) CreateBuild(ctx context.Context, integrationID, ciCommitId, title string) (*Build, error) {
	return nil, fmt.Errorf("CreateBuild is not implemented in this version of Sturdy")
}

func New() Service {
	return &svc{}
}

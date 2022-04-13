package service

import (
	"context"
	"errors"

	"getsturdy.com/api/pkg/codebases"
)

type Service interface {
	Pull(ctx context.Context, codebaseID codebases.ID) error
	PushTrunk(ctx context.Context, codebaseID codebases.ID) error
}

type service struct{}

func New() Service {
	return &service{}
}

func (*service) Pull(context.Context, codebases.ID) error {
	return errors.New("not available")
}

func (*service) PushTrunk(context.Context, codebases.ID) error {
	return errors.New("not available")
}

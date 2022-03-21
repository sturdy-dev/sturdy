package service

import (
	"context"
	"errors"

	"getsturdy.com/api/pkg/codebases"
	remote_service "getsturdy.com/api/pkg/remote/service"
)

type service struct{}

func New() remote_service.Service {
	return &service{}
}

func (*service) Pull(_ context.Context, _ codebases.ID) error {
	return errors.New("not available")
}

package validator

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"getsturdy.com/api/pkg/graphql/model"
	"getsturdy.com/api/pkg/license/enterprise/client"
	service_user "getsturdy.com/api/pkg/user/service"
)

type Validator struct {
	client      *client.Client
	logger      *zap.Logger
	userService *service_user.Service

	lastStatus *model.LicenseValidation
}

func New(
	client *client.Client,
	logger *zap.Logger,
	userService *service_user.Service,
) *Validator {
	return &Validator{
		client:      client,
		logger:      logger,
		userService: userService,
	}
}

func (v *Validator) Run() error {
	return nil
	// TODO: Enable on self-hosted installations

	ctx := context.Background()
	for {
		if v.run(ctx) == nil {
			time.Sleep(time.Minute * 10)
		} else {
			time.Sleep(time.Minute)
		}
	}
}

func (v *Validator) run(ctx context.Context) error {
	userCount, err := v.userService.UserCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user count: %w", err)
	}

	val, err := v.client.Validate(ctx, "foobarLicenseKey", userCount)
	v.logger.Info("license validation", zap.Any("status", val), zap.Error(err))
	v.lastStatus = val

	if err != nil {
		return fmt.Errorf("failed to validate license: %w", err)
	}

	return err
}

func (v *Validator) Status() *model.LicenseValidation {
	return v.lastStatus
}

package worker

import (
	"context"
	"fmt"
	"time"

	service_installations "getsturdy.com/api/pkg/installations/service"
	validator_license "getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"

	backoff "github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

var (
	runEvery = time.Hour
)

type Worker struct {
	logger               *zap.Logger
	validator            *validator_license.Validator
	installationsService *service_installations.Service
}

func New(
	logger *zap.Logger,
	validator *validator_license.Validator,
	installationsService *service_installations.Service,
) *Worker {
	return &Worker{
		logger:               logger.Named("licenses_worker"),
		validator:            validator,
		installationsService: installationsService,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info("starting")

	if err := w.attemptValidation(ctx); err != nil {
		w.logger.Error("failed to validate license", zap.Error(err))
	}

	ticker := time.NewTicker(runEvery)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := w.attemptValidation(ctx); err != nil {
				w.logger.Error("failed to validate license", zap.Error(err))
			}
		case <-ctx.Done():
			w.logger.Info("stopping")
			return nil
		}
	}
}

func (w *Worker) attemptValidation(ctx context.Context) error {
	exp := backoff.NewExponentialBackOff()
	exp.MaxElapsedTime = runEvery
	return backoff.Retry(func() error {
		w.logger.Info("validating license")
		if err := w.validateLicense(ctx); err != nil {
			w.logger.Warn("failed to validate license", zap.Error(err))
			return err
		}
		return nil
	}, exp)
}

func (w *Worker) validateLicense(ctx context.Context) error {
	installation, err := w.installationsService.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get installation: %w", err)
	}

	if installation.LicenseKey == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	license, err := w.validator.Validate(ctx, *installation.LicenseKey)
	if err != nil {
		return fmt.Errorf("failed to validate license: %w", err)
	}

	if err := w.installationsService.UpdateLicense(ctx, license); err != nil {
		return fmt.Errorf("failed to update license: %w", err)
	}

	return nil
}

package enterprise

import (
	"context"
	"fmt"

	"mash/pkg/api"
	"mash/pkg/api/oss"
	workers_github "mash/pkg/github/enterprise/workers"
	"mash/pkg/license/enterprise/validator"

	"golang.org/x/sync/errgroup"
)

type API struct {
	ossAPI *oss.API

	githubClonerQueue   *workers_github.ClonerQueue
	githubImporterQueue workers_github.ImporterQueue
	licenseValidator    *validator.Validator
}

func ProvideAPI(
	ossAPI *oss.API,

	githubClonerQueue *workers_github.ClonerQueue,
	githubImporterQueue workers_github.ImporterQueue,
	licenseValidator *validator.Validator,
) *API {
	return &API{
		ossAPI:              ossAPI,
		githubClonerQueue:   githubClonerQueue,
		githubImporterQueue: githubImporterQueue,
		licenseValidator:    licenseValidator,
	}
}

func (a *API) Start(ctx context.Context, cfg *api.Config) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return a.ossAPI.Start(ctx, cfg)
	})

	// github cloner queue
	wg.Go(func() error {
		if err := a.githubClonerQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github cloner queue: %v", err)
		}
		return nil
	})

	wg.Go(func() error {
		if err := a.githubImporterQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github importer queue: %v", err)
		}
		return nil
	})

	wg.Go(func() error {
		if err := a.licenseValidator.Run(); err != nil {
			return fmt.Errorf("failed to start license validator: %w", err)
		}
		return nil
	})

	return wg.Wait()
}

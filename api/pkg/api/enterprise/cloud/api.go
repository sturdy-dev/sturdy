package cloud

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/api"
	workers_github "getsturdy.com/api/pkg/github/enterprise/workers"

	"golang.org/x/sync/errgroup"
)

type API struct {
	ossAPI *api.API

	githubClonerQueue   *workers_github.ClonerQueue
	githubImporterQueue workers_github.ImporterQueue
}

func ProvideAPI(
	ossAPI *api.API,

	githubClonerQueue *workers_github.ClonerQueue,
	githubImporterQueue workers_github.ImporterQueue,
) *API {
	return &API{
		ossAPI:              ossAPI,
		githubClonerQueue:   githubClonerQueue,
		githubImporterQueue: githubImporterQueue,
	}
}

func (a *API) Start(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return a.ossAPI.Start(ctx)
	})

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

	return wg.Wait()
}

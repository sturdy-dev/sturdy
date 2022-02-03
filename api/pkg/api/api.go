package api

import (
	"context"
	"fmt"

	workers_ci "getsturdy.com/api/pkg/ci/workers"
	worker_gc "getsturdy.com/api/pkg/gc/worker"
	"getsturdy.com/api/pkg/gitserver"
	httpx "getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/metrics"
	"getsturdy.com/api/pkg/pprof"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"

	"golang.org/x/sync/errgroup"
)

type Starter interface {
	Start(context.Context) error
}

type API struct {
	httpServer       *httpx.Server
	snapshotterQueue worker_snapshots.Queue
	ciBuildQueue     *workers_ci.BuildQueue
	gcQueue          *worker_gc.Queue
	gitsrv           *gitserver.Server
	pprof            *pprof.Server
	metrics          *metrics.Server
}

func ProvideAPI(
	httpServer *httpx.Server,
	snapshotterQueue worker_snapshots.Queue,
	ciBuildQueue *workers_ci.BuildQueue,
	gcQueue *worker_gc.Queue,
	gitsrv *gitserver.Server,
	pprof *pprof.Server,
	metrics *metrics.Server,
) *API {
	return &API{
		httpServer:       httpServer,
		snapshotterQueue: snapshotterQueue,
		ciBuildQueue:     ciBuildQueue,
		gcQueue:          gcQueue,
		gitsrv:           gitsrv,
		pprof:            pprof,
		metrics:          metrics,
	}
}

func (a *API) Start(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)
	// snapshotter queue
	wg.Go(func() error {
		if err := a.snapshotterQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start snapshotter queue: %w", err)
		}
		return nil
	})
	// ci build queue
	wg.Go(func() error {
		if err := a.ciBuildQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start ci build queue: %w", err)
		}
		return nil
	})
	// gc queue
	wg.Go(func() error {
		if err := a.gcQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start gc queue: %w", err)
		}
		return nil
	})
	// Start the git HTTP server
	wg.Go(func() error {
		if err := a.gitsrv.Start(); err != nil {
			return fmt.Errorf("failed to start git server: %w", err)
		}
		return nil
	})
	// Pprof server
	wg.Go(func() error {
		if err := a.pprof.Start(); err != nil {
			return fmt.Errorf("failed to start http pprof server: %w", err)
		}
		return nil
	})
	// Metrics server
	wg.Go(func() error {
		if err := a.metrics.Start(); err != nil {
			return fmt.Errorf("failed to start http metrics server: %w", err)
		}
		return nil
	})
	wg.Go(func() error {
		if err := a.httpServer.Start(); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	return wg.Wait()
}

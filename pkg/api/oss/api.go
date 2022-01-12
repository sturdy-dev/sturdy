package oss

import (
	"context"
	"fmt"
	"net/http"

	"mash/pkg/api"
	workers_ci "mash/pkg/ci/workers"
	worker_gc "mash/pkg/gc/worker"
	"mash/pkg/gitserver"
	httpx "mash/pkg/http"
	worker_snapshots "mash/pkg/snapshots/worker"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

type API struct {
	httpServer       *httpx.Server
	snapshotterQueue worker_snapshots.Queue
	ciBuildQueue     *workers_ci.BuildQueue
	gcQueue          *worker_gc.Queue
	gitsrv           *gitserver.Server
}

func ProvideAPI(
	httpServer *httpx.Server,
	snapshotterQueue worker_snapshots.Queue,
	ciBuildQueue *workers_ci.BuildQueue,
	gcQueue *worker_gc.Queue,
	gitsrv *gitserver.Server,
) *API {
	return &API{
		httpServer:       httpServer,
		snapshotterQueue: snapshotterQueue,
		ciBuildQueue:     ciBuildQueue,
		gcQueue:          gcQueue,
		gitsrv:           gitsrv,
	}
}

func (a *API) Start(ctx context.Context, cfg *api.Config) error {
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
		if err := a.gitsrv.Start(ctx, cfg.GitListenAddr); err != nil {
			return fmt.Errorf("failed to start git server: %w", err)
		}
		return nil
	})
	// Pprof server
	wg.Go(func() error {
		if err := http.ListenAndServe(cfg.HTTPPProfListenAddr, nil); err != http.ErrServerClosed {
			return fmt.Errorf("failed to start http pprof server: %w", err)
		}
		return nil
	})
	// Metrics server
	wg.Go(func() error {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		srv := http.Server{Addr: cfg.MetricsListenAddr, Handler: mux}
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		return nil
	})
	wg.Go(func() error {
		if err := a.httpServer.ListenAndServe(cfg.HTTPAddr); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	return wg.Wait()
}

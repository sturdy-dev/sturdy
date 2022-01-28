package worker

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/publisher"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/service"

	backoff "github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

var (
	runEvery = time.Hour
)

type Worker struct {
	logger    *zap.Logger
	service   *service.Service
	publisher *publisher.Publisher
}

func New(
	logger *zap.Logger,
	service *service.Service,
	publisher *publisher.Publisher,
) *Worker {
	return &Worker{
		logger:    logger.Named("statistics_worker"),
		service:   service,
		publisher: publisher,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info("starting")

	if err := w.attemptSendStatistics(ctx); err != nil {
		w.logger.Error("failed to send statistics", zap.Error(err))
	}

	ticker := time.NewTicker(runEvery)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := w.attemptSendStatistics(ctx); err != nil {
				w.logger.Error("failed to send statistics", zap.Error(err))
			}
		case <-ctx.Done():
			w.logger.Info("stopping")
			return nil
		}
	}
}

func (w *Worker) attemptSendStatistics(ctx context.Context) error {
	exp := backoff.NewExponentialBackOff()
	exp.MaxElapsedTime = runEvery
	return backoff.Retry(func() error {
		w.logger.Info("sending statistics")
		if err := w.sendStatistics(ctx); err != nil {
			w.logger.Warn("failed to send statistics", zap.Error(err))
			return err
		}
		return nil
	}, exp)
}

func (w *Worker) sendStatistics(ctx context.Context) error {
	statistic, err := w.service.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := w.publisher.Publish(ctx, statistic); err != nil {
		return fmt.Errorf("failed to publish statistics: %w", err)
	}

	return nil
}

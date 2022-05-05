package events

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type subscriber struct {
	ctx      context.Context
	callback callback
}

type callback func(context.Context, *event) error

type subscriptionID string

type pubSub struct {
	logger *zap.Logger

	subscribersGuard *sync.RWMutex
	subscribers      map[Topic]map[Type]map[subscriptionID]subscriber
}

func New(logger *zap.Logger) *pubSub {
	return &pubSub{
		logger:           logger.Named("events_pubsub"),
		subscribersGuard: &sync.RWMutex{},
		subscribers:      map[Topic]map[Type]map[subscriptionID]subscriber{},
	}
}

func (r *pubSub) pub(topic Topic, evt *event) {
	r.subscribersGuard.RLock()
	handlers := r.subscribers[topic][evt.Type]
	r.subscribersGuard.RUnlock()

	logger := r.logger.With(zap.Stringer("topic", topic), zap.Stringer("type", evt.Type))

	for _, handler := range handlers {
		handler := handler
		go func() {
			start := time.Now()

			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic in events v2 publisher", zap.Any("recover", rec), zap.Stack("stack"),
						zap.Duration("duration", time.Since(start)),
						zap.String("handler", functionName(handler.callback)),
					)
				}
			}()

			if err := handler.callback(handler.ctx, evt); err != nil {
				logger.Error(
					"failed to handle event",
					zap.Duration("duration", time.Since(start)),
					zap.Error(errors.Unwrap(err)),
					zap.String("handler", functionName(handler.callback)),
				)
			}
		}()
	}
}

func functionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func (r *pubSub) sub(ctx context.Context, fn callback, topic Topic, tt ...Type) {
	id := subscriptionID(uuid.NewString())

	r.subscribersGuard.Lock()
	for _, t := range tt {
		if r.subscribers[topic] == nil {
			r.subscribers[topic] = make(map[Type]map[subscriptionID]subscriber)
		}
		if r.subscribers[topic][t] == nil {
			r.subscribers[topic][t] = make(map[subscriptionID]subscriber)
		}
		r.subscribers[topic][t][id] = subscriber{ctx: ctx, callback: fn}
	}
	r.subscribersGuard.Unlock()

	go func() {
		<-ctx.Done()
		r.subscribersGuard.Lock()
		for _, t := range tt {
			delete(r.subscribers[topic][t], id)
		}
		r.subscribersGuard.Unlock()
	}()
}

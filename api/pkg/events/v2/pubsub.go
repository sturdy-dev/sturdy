package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"getsturdy.com/api/pkg/users"

	"go.uber.org/zap"
)

type subscriber struct {
	ctx      context.Context
	callback callback
}

type callback func(context.Context, *event) error

type Topic string

func (t Topic) String() string {
	return string(t)
}

func user(userID users.ID) Topic {
	return Topic(fmt.Sprintf("user:%s", userID))
}

type subscriptionID string

type PubSub struct {
	logger *zap.Logger

	subscribersGuard *sync.RWMutex
	subscribers      map[Topic]map[Type]map[subscriptionID]subscriber
}

func New(logger *zap.Logger) *PubSub {
	return &PubSub{
		logger:           logger.Named("events_pubsub"),
		subscribersGuard: &sync.RWMutex{},
		subscribers:      map[Topic]map[Type]map[subscriptionID]subscriber{},
	}
}

func (r *PubSub) pub(topic Topic, evt *event) {
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
					)
				}
			}()

			if err := handler.callback(handler.ctx, evt); err != nil {
				logger.Error("failed to handle event", zap.Error(err), zap.Duration("duration", time.Since(start)))
			}
		}()
	}
}

func (r *PubSub) sub(ctx context.Context, fn callback, topic Topic, tt ...Type) {
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

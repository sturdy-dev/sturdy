package events

import (
	"context"
	"fmt"
	"sync"

	"getsturdy.com/api/pkg/users"

	"go.uber.org/zap"
)

type subscriber func(context.Context, *event) error

type Topic string

func User(userID users.ID) Topic {
	return Topic(fmt.Sprintf("user:%s", userID))
}

func Workspace(id string) Topic {
	return Topic(fmt.Sprintf("workspace:%s", id))
}

type PubSub struct {
	logger *zap.Logger

	subscribersGuard *sync.RWMutex
	subscribers      map[Topic]map[Type][]subscriber
}

func New(logger *zap.Logger) *PubSub {
	return &PubSub{
		logger:           logger.Named("events_pubsub"),
		subscribersGuard: &sync.RWMutex{},
		subscribers:      map[Topic]map[Type][]subscriber{},
	}
}

func (r *PubSub) pub(topic Topic, evt *event) {
	r.subscribersGuard.RLock()
	handlers := r.subscribers[topic][evt.Type]
	r.subscribersGuard.RUnlock()

	ctx := context.Background()
	for _, fn := range handlers {
		fn := fn
		go func() {
			if err := fn(ctx, evt); err != nil {
				r.logger.Error("failed to publish event", zap.Error(err))
			}
		}()
	}
}

func (r *PubSub) sub(fn subscriber, topic Topic, tt ...Type) {
	r.subscribersGuard.Lock()
	for _, t := range tt {
		r.subscribers[topic][t] = append(r.subscribers[topic][t], fn)
	}
	r.subscribersGuard.Unlock()
}

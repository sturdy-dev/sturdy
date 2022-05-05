package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"getsturdy.com/api/pkg/queue/names"

	"go.uber.org/zap"
)

var _ Queue = &InMemoryQueue{}

type InMemoryQueue struct {
	logger     *zap.Logger
	sync       bool
	chansGuard sync.RWMutex
	chans      map[names.IncompleteQueueName]chan Message

	bufferSize int
	timeout    time.Duration
}

func NewInMemory(logger *zap.Logger) *InMemoryQueue {
	return &InMemoryQueue{
		logger:     logger.Named("inmemory queue"),
		chans:      make(map[names.IncompleteQueueName]chan Message),
		bufferSize: 10,
		timeout:    time.Second,
	}
}

func (q *InMemoryQueue) Sync() *InMemoryQueue {
	q.sync = true
	return q
}

type inmemorymessage struct {
	marshalledMessage []byte
	ack               chan struct{}
}

func newInmemoryMessage(v any) (*inmemorymessage, error) {
	marshaled, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &inmemorymessage{
		marshalledMessage: marshaled,
		ack:               make(chan struct{}),
	}, nil
}

func (m *inmemorymessage) AwaitAcked() {
	<-m.ack
}

func (m *inmemorymessage) Ack() error {
	close(m.ack)
	return nil
}

func (m *inmemorymessage) As(v any) error {
	return json.Unmarshal(m.marshalledMessage, v)
}

func (q *InMemoryQueue) getChannel(name names.IncompleteQueueName) chan Message {
	q.chansGuard.RLock()
	ch, found := q.chans[name]
	q.chansGuard.RUnlock()
	if found {
		return ch
	}
	q.chansGuard.Lock()
	ch = make(chan Message, q.bufferSize)
	q.chans[name] = ch
	q.chansGuard.Unlock()
	return ch
}

func (q *InMemoryQueue) Publish(ctx context.Context, name names.IncompleteQueueName, msg any) error {
	q.logger.Info("publishing message", zap.String("queue", name.String()))

	ch := q.getChannel(name)
	m, err := newInmemoryMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	select {
	case ch <- m:
		if q.sync {

			fmt.Printf("\nnikitag: %+v\n\n", name.String())

			m.AwaitAcked()
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(q.timeout):
		return fmt.Errorf("failed to publish message: queue is full")
	}
}

func (q *InMemoryQueue) Subscribe(ctx context.Context, name names.IncompleteQueueName, messages chan<- Message) error {
	q.logger.Info("subscribing to queue", zap.String("queue", name.String()))
	ch := q.getChannel(name)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m := <-ch:
			q.logger.Info("received message", zap.String("queue", name.String()))
			messages <- m
		}
	}
}

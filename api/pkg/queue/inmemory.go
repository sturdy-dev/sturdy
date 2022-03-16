package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"getsturdy.com/api/pkg/queue/names"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

var _ Queue = &memoryQueue{}

type memoryQueue struct {
	logger     *zap.Logger
	chansGuard sync.RWMutex
	chans      map[names.IncompleteQueueName][]chan<- Message
}

func NewInMemory(logger *zap.Logger) *memoryQueue {
	return &memoryQueue{
		logger: logger,
		chans:  make(map[names.IncompleteQueueName][]chan<- Message),
	}
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

func (q *memoryQueue) Publish(ctx context.Context, name names.IncompleteQueueName, msg any) error {
	q.chansGuard.RLock()
	defer q.chansGuard.RUnlock()
	chans, ok := q.chans[name]
	if !ok {
		return nil
	}

	wg, _ := errgroup.WithContext(ctx)
	for _, ch := range chans {
		ch := ch
		wg.Go(func() error {
			m, err := newInmemoryMessage(msg)
			if err != nil {
				return fmt.Errorf("failed to create message: %w", err)
			}
			ch <- m
			return nil
		})
	}
	return wg.Wait()
}

func (q *memoryQueue) Subscribe(ctx context.Context, name names.IncompleteQueueName, messages chan<- Message) error {
	q.chansGuard.Lock()
	q.chans[name] = append(q.chans[name], messages)
	q.chansGuard.Unlock()

	<-ctx.Done()
	return nil
}

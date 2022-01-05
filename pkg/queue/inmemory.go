package queue

import (
	"context"
	"encoding/json"
	"mash/pkg/queue/names"
	"sync"

	"go.uber.org/zap"
)

var _ Queue = &memoryQueue{}

type memoryQueue struct {
	logger *zap.Logger
	qGuard *sync.RWMutex
	q      map[names.IncompleteQueueName]chan Message
}

func NewInMemory(logger *zap.Logger) *memoryQueue {
	return &memoryQueue{
		logger: logger,
		qGuard: &sync.RWMutex{},
		q:      make(map[names.IncompleteQueueName]chan Message),
	}
}

type inmemorymessage struct {
	marshalledMessage []byte
	ack               bool
}

func (m *inmemorymessage) Ack() error {
	m.ack = true
	return nil
}

func (m *inmemorymessage) As(v interface{}) error {
	return json.Unmarshal(m.marshalledMessage, v)
}

func (q *memoryQueue) Publish(_ context.Context, name names.IncompleteQueueName, msg interface{}) error {
	q.logger.Info("publishing message", zap.String("queue", string(name)))

	q.qGuard.Lock()
	ch, ok := q.q[name]
	if !ok {
		ch = make(chan Message)
		q.q[name] = ch
	}
	q.qGuard.Unlock()

	marshalled, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	go func() {
		ch <- &inmemorymessage{
			marshalledMessage: marshalled,
		}
	}()
	return nil
}

func (q *memoryQueue) Subscribe(ctx context.Context, name names.IncompleteQueueName, messages chan<- Message) error {
	q.logger.Info("new subscription", zap.String("queue", string(name)))

	q.qGuard.Lock()
	ch, ok := q.q[name]
	if !ok {
		ch = make(chan Message)
		q.q[name] = ch
	}
	q.qGuard.Unlock()

	for {
		select {
		case <-ctx.Done():
			q.logger.Info("stopping subscription", zap.String("queue", string(name)))
			return nil
		case msg := <-ch:
			q.logger.Info("new message", zap.String("queue", string(name)))
			messages <- msg
		}
	}
}

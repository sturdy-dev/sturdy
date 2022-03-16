package queue

import (
	"context"

	"getsturdy.com/api/pkg/queue/names"
)

type noopQueue struct{}

func NewNoop() *noopQueue {
	return &noopQueue{}
}

func (*noopQueue) Publish(context.Context, names.IncompleteQueueName, any) error {
	return nil
}

func (*noopQueue) Subscribe(ctx context.Context, _ names.IncompleteQueueName, _ chan<- Message) error {
	<-ctx.Done()
	return nil
}

package queue

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"getsturdy.com/api/pkg/queue/names"
)

var _ Queue = &Sync{}

type Sync struct {
	chans map[names.IncompleteQueueName][]chan<- Message
}

func NewSync() *Sync {
	return &Sync{
		chans: make(map[names.IncompleteQueueName][]chan<- Message),
	}
}

func (q *Sync) Publish(ctx context.Context, name names.IncompleteQueueName, msg any) error {
	chs, ok := q.chans[name]
	if !ok {
		return nil
	}

	wg, _ := errgroup.WithContext(ctx)
	for _, ch := range chs {
		ch := ch
		wg.Go(func() error {
			m, err := newInmemoryMessage(msg)
			if err != nil {
				return fmt.Errorf("failed to create message: %w", err)
			}
			ch <- m
			m.AwaitAcked()
			return nil
		})
	}
	return wg.Wait()
}

func (q *Sync) Subscribe(ctx context.Context, name names.IncompleteQueueName, mesasges chan<- Message) error {
	q.chans[name] = append(q.chans[name], mesasges)
	<-ctx.Done()
	return nil
}

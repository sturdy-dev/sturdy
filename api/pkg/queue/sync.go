package queue

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/queue/names"
	"golang.org/x/sync/errgroup"
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

func (q *Sync) Publish(ctx context.Context, name names.IncompleteQueueName, msg interface{}) error {
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
				return fmt.Errorf("failed to create message: %v", err)
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

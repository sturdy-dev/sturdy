package queue

import (
	"context"
	"testing"

	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/queue/names"

	"github.com/stretchr/testify/assert"
)

func TestInmemory__once(t *testing.T) {
	q := NewInMemory(logger.NewTest(t))
	name := names.IncompleteQueueName("testing")

	msgs := make(chan Message)
	go func() { // consumer 1
		assert.NoError(t, q.Subscribe(context.TODO(), name, msgs))
	}()
	go func() { // consumer 2
		assert.NoError(t, q.Subscribe(context.TODO(), name, msgs))
	}()
	go func() { // producer
		for i := 0; i < 10; i++ {
			assert.NoError(t, q.Publish(context.TODO(), name, i))
		}
	}()

	got := make([]int, 0, 10)
	for msg := range msgs {
		var i int
		assert.NoError(t, msg.As(&i))

		assert.NotContains(t, got, i) // no duplicates
		got = append(got, i)

		assert.NoError(t, msg.Ack())
		if len(got) == 10 {
			break
		}
	}
}

func TestInmemory__fifo(t *testing.T) {
	q := NewInMemory(logger.NewTest(t))
	name := names.IncompleteQueueName("testing")

	msgs := make(chan Message)
	go func() { // consumer 1
		assert.NoError(t, q.Subscribe(context.TODO(), name, msgs))
	}()
	go func() { // producer
		for i := 0; i < 10; i++ {
			assert.NoError(t, q.Publish(context.TODO(), name, i))
		}
	}()

	got := make([]int, 0, 10)
	for msg := range msgs {
		var i int
		assert.NoError(t, msg.As(&i))

		assert.Equal(t, len(got), i)

		got = append(got, i)

		assert.NoError(t, msg.Ack())
		if len(got) == 10 {
			break
		}
	}
}

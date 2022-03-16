package queue

import (
	"context"

	"getsturdy.com/api/pkg/queue/names"
)

type Queue interface {
	// Publish publishes a message to the queue.
	Publish(context.Context, names.IncompleteQueueName, any) error
	// Subscribe returns a channel that will receive messages from the queue.
	Subscribe(context.Context, names.IncompleteQueueName, chan<- Message) error
}

type Message interface {
	// As unmarshals message into the given interface.
	As(any) error
	// Ack marks the message as acknowleged.
	Ack() error
}

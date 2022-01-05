package queue

import (
	"context"

	"mash/pkg/queue/names"
)

type Queue interface {
	// Publish publishes a message to the queue.
	Publish(context.Context, names.IncompleteQueueName, interface{}) error
	// Subscribe returns a channel that will receive messages from the queue.
	Subscribe(context.Context, names.IncompleteQueueName, chan<- Message) error
}

type Message interface {
	// As unmarshals message into the given interface.
	As(interface{}) error
	// Ack marks the message as acknowleged.
	Ack() error
}

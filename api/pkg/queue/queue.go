package queue

import (
	"context"
	"fmt"
	"os"

	"getsturdy.com/api/pkg/queue/names"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type Configuration struct {
	Local    bool   `long:"local" description:"Use in-memory queue instead of SQS"`
	Hostname string `long:"hostname" description:"Hostname of the queue"`
	Prefix   string `long:"prefix" description:"Prefix for queue names"`
}

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

func New(awsSession *session.Session, logger *zap.Logger, cfg *Configuration) (Queue, error) {
	if cfg.Hostname == "" {
		defaultHostname, err := os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("failed to get hostname: %v", err)
		}
		cfg.Hostname = defaultHostname
	}
	if cfg.Local {
		return NewInMemory(logger), nil
	}
	return NewSQS(logger, awsSession, cfg.Hostname, cfg.Prefix)
}

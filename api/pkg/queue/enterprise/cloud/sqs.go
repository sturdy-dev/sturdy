package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/broadcast"
	"getsturdy.com/api/pkg/queue/names"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.uber.org/zap"
)

var _ queue.Queue = &sqsQueue{}

type sqsQueue struct {
	logger *zap.Logger

	session     *session.Session
	hostname    string
	queuePrefix string

	publishersGuard *sync.RWMutex
	publishers      map[names.IncompleteQueueName]publisher
}

func NewSQS(
	logger *zap.Logger,
	awsSession *session.Session,
	hostname string,
	queuePrefix string,
) (*sqsQueue, error) {
	return &sqsQueue{
		logger:      logger.Named("sqsQueue"),
		session:     awsSession,
		hostname:    hostname,
		queuePrefix: queuePrefix,

		publishersGuard: &sync.RWMutex{},
		publishers:      map[names.IncompleteQueueName]publisher{},
	}, nil
}

func (q *sqsQueue) getPublisher(name names.IncompleteQueueName) (publisher, error) {
	q.publishersGuard.RLock()
	cached, found := q.publishers[name]
	q.publishersGuard.RUnlock()

	if found {
		return cached, nil
	}

	publisher, err := newPublisher(q.logger, q.session, names.BuildQueueName(q.queuePrefix, q.hostname, name))
	if err != nil {
		return nil, fmt.Errorf("failed to create aws publisher: %w", err)
	}

	q.publishersGuard.Lock()
	q.publishers[name] = publisher
	q.publishersGuard.Unlock()

	return publisher, nil
}

func (q *sqsQueue) Publish(_ context.Context, name names.IncompleteQueueName, v any) error {
	q.logger.Info("publishing message", zap.String("queue", string(name)))

	publish, err := q.getPublisher(name)
	if err != nil {
		return fmt.Errorf("failed to create sqs publisher: %w", err)
	}

	if err := publish(v); err != nil {
		return fmt.Errorf("failed to publish message to sqs: %w", err)
	}
	return nil
}

func (q *sqsQueue) Subscribe(ctx context.Context, name names.IncompleteQueueName, messages chan<- queue.Message) error {
	q.logger.Info("new subscription", zap.String("queue", string(name)))

	sqsMessages := make(chan queue.Message)
	if err := Sub(q.session, q.logger, names.BuildQueueName(q.queuePrefix, q.hostname, name), sqsMessages); err != nil {
		return fmt.Errorf("failed to subscribe to queue: %w", err)
	}

	for {
		select {
		case msg := <-sqsMessages:
			q.logger.Info("new message", zap.String("queue", string(name)))
			messages <- msg
		case <-ctx.Done():
			q.logger.Info("stopping subscription", zap.String("queue", string(name)))
			return nil
		}
	}
}

type message struct {
	q             *sqs.SQS
	receiptHandle *string
	queueUrl      *string
	data          *string
}

func (m *message) As(out any) error {
	if err := unmarshal([]byte(*m.data), out); err != nil {
		// try parsing using deprecated message format
		return m.AsJSON(out)
	}
	return nil
}

// deprecated message format
func (m *message) AsJSON(out any) error {
	if err := json.Unmarshal([]byte(*m.data), out); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return nil
}

func (m *message) Ack() error {
	_, err := m.q.DeleteMessage(&sqs.DeleteMessageInput{
		ReceiptHandle: m.receiptHandle,
		QueueUrl:      m.queueUrl,
	})
	if err != nil {
		return err
	}
	return nil
}

func SubToBroadcast(awsSession *session.Session, logger *zap.Logger, queueName names.BroadcastQueueSubscriberName, snsTopicArn string, output chan queue.Message) error {
	return sub(awsSession, logger, string(queueName), snsTopicArn, output)
}

func Sub(awsSession *session.Session, logger *zap.Logger, queueName names.QueueName, output chan queue.Message) error {
	// Subscribe without setting up a SNS subscription
	return sub(awsSession, logger, string(queueName), "", output)
}

func sub(awsSession *session.Session, logger *zap.Logger, queueName string, snsTopicArn string, output chan queue.Message) error {
	q := sqs.New(awsSession)
	stsClient := sts.New(awsSession)

	queueUrl, queueArn, err := getOrCreateQueue(logger, q, stsClient, queueName, snsTopicArn)
	if err != nil {
		return err
	}

	// Subscribe this queue to the SNS topic
	if snsTopicArn != "" {
		err := broadcast.SetupQueueSubscription(awsSession, snsTopicArn, queueArn)
		if err != nil {
			return err
		}
	}

	go func() {
		for {
			msgs, err := q.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueUrl),
				MaxNumberOfMessages: aws.Int64(10),
				WaitTimeSeconds:     aws.Int64(20),
			})

			// Perform the first fetch in the foreground, to make sure that everything is working
			if err != nil {
				logger.Error("failed to fetch sqs", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			for _, m := range msgs.Messages {
				output <- &message{
					q:             q,
					receiptHandle: m.ReceiptHandle,
					queueUrl:      &queueUrl,
					data:          m.Body,
				}
			}
		}
	}()

	return nil
}

type publisher func(msg any) error

func newPublisher(logger *zap.Logger, awsSession *session.Session, queueName names.QueueName) (publisher, error) {
	q := sqs.New(awsSession)
	stsClient := sts.New(awsSession)

	queueUrl, _, err := getOrCreateQueue(logger, q, stsClient, string(queueName), "")
	if err != nil {
		return nil, err
	}

	publ := func(msg any) error {
		body, err := marshal(msg)
		if err != nil {
			return err
		}

		_, err = q.SendMessage(&sqs.SendMessageInput{
			QueueUrl:    &queueUrl,
			MessageBody: aws.String(string(body)),
		})
		if err != nil {
			return err
		}
		return nil
	}

	return publ, nil
}

func getQueueArn(q *sqs.SQS, queueUrl string) (string, error) {
	attributes, err := q.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueUrl),
		AttributeNames: []*string{
			aws.String("QueueArn"),
		},
	})
	if err != nil {
		return "", err
	}

	if arn, ok := attributes.Attributes["QueueArn"]; ok {
		return *arn, nil
	}

	return "", fmt.Errorf("could not get QueueArn")
}

func getOrCreateQueue(logger *zap.Logger, q *sqs.SQS, stsClient *sts.STS, queueName string, snsTopicArn string) (queueUrl string, queueArn string, err error) {
	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", "", err
	}
	r, err := q.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	// If queue exists
	if err == nil {
		arn, err := getQueueArn(q, *r.QueueUrl)
		if err != nil {
			return "", "", err
		}

		logger.Info("using existing SQS queue", zap.String("queueUrl", *r.QueueUrl))

		return *r.QueueUrl, arn, nil
	}

	// Queue does not exist
	//nolint:errorlint
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {

		pol := PolicyDocument{
			Version:   "2008-10-17",
			Statement: make([]StatementEntry, 0),
		}

		pol.Statement = append(pol.Statement, StatementEntry{
			Effect: "Allow",
			Principal: map[string]string{
				"AWS": "arn:aws:iam::" + *callerIdentity.Account + ":root",
			},
			Action:   []string{"SQS:*"},
			Resource: "arn:aws:sqs:eu-north-1:" + *callerIdentity.Account + ":" + queueName,
		})

		// If subbing from SNS, allow SNS to send messages
		if snsTopicArn != "" {
			pol.Statement = append(pol.Statement, StatementEntry{
				Effect: "Allow",
				Principal: map[string]string{
					"Service": "sns.amazonaws.com",
				},
				Action:   []string{"SQS:SendMessage"},
				Resource: "arn:aws:sqs:eu-north-1:" + *callerIdentity.Account + ":" + queueName,
				Condition: map[string]map[string]string{
					"ArnEquals": {
						"aws:SourceArn": snsTopicArn,
					},
				},
			})
		}

		policyJson, err := json.MarshalIndent(pol, "", "  ")
		if err != nil {
			return "", "", fmt.Errorf("could not build policy: %w", err)
		}

		// Create the Dead Letter Queue
		createdDeadLetterQueue, createQueueErr := q.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String(queueName + "_dead"),
			Attributes: map[string]*string{
				"KmsMasterKeyId": aws.String("alias/sns_and_sqs"),
				"Policy":         aws.String(string(policyJson)),
			},
		})
		if createQueueErr != nil {
			return "", "", createQueueErr
		}
		dlqArn, err := getQueueArn(q, *createdDeadLetterQueue.QueueUrl)
		if err != nil {
			return "", "", err
		}

		redrivePolicy, err := json.Marshal(map[string]any{
			"deadLetterTargetArn": dlqArn,
			"maxReceiveCount":     5,
		})
		if err != nil {
			return "", "", err
		}

		// Create queue
		created, createQueueErr := q.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String(queueName),
			Attributes: map[string]*string{
				"KmsMasterKeyId": aws.String("alias/sns_and_sqs"),
				"Policy":         aws.String(string(policyJson)),
				"RedrivePolicy":  aws.String(string(redrivePolicy)),
			},
		})
		if createQueueErr != nil {
			return "", "", createQueueErr
		}

		arn, err := getQueueArn(q, *created.QueueUrl)
		if err != nil {
			return "", "", err
		}

		logger.Info("created new SQS queues", zap.String("queueUrl", *created.QueueUrl), zap.String("dlqQueueUrl", *createdDeadLetterQueue.QueueUrl))

		// Created queue
		return *created.QueueUrl, arn, nil
	}

	// Different type of error
	return "", "", err
}

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect    string
	Principal map[string]string
	Action    []string
	Resource  string
	Condition map[string]map[string]string `json:"Condition,omitempty"`
}

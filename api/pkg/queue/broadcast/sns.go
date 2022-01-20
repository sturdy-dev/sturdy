package broadcast

import (
	"encoding/json"
	"fmt"
	"getsturdy.com/api/pkg/queue/names"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.uber.org/zap"
)

type Publisher func(msg interface{}) error

func NewPublisher(awsSession *session.Session, logger *zap.Logger, topicName names.BroadcastQueuePublisherName) (pub Publisher, snsArn string, err error) {
	s := sns.New(awsSession)

	topics, err := s.ListTopics(&sns.ListTopicsInput{})
	if err != nil {
		return nil, "", err
	}

	var hasTopic bool
	var topicArn string
	for _, topic := range topics.Topics {
		if strings.HasSuffix(*topic.TopicArn, ":"+string(topicName)) {
			hasTopic = true
			topicArn = *topic.TopicArn
		}
	}

	if err := assertKey(awsSession, logger); err != nil {
		return nil, "", fmt.Errorf("failed to setup key: %w", err)
	}

	// Create the topic
	if !hasTopic {
		logger.Info("setting up SNS topic", zap.String("name", string(topicName)))
		createdTopic, err := s.CreateTopic(&sns.CreateTopicInput{
			Name: aws.String(string(topicName)),
			Attributes: map[string]*string{
				"KmsMasterKeyId": aws.String("alias/sns_and_sqs"),
			},
			Tags: []*sns.Tag{{Key: aws.String("SturdyBroadcast"), Value: aws.String("true")}},
		})
		if err != nil {
			return nil, "", fmt.Errorf("failed to create topic: %w", err)
		}
		topicArn = *createdTopic.TopicArn
	}

	publ := func(msg interface{}) error {
		body, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		_, err = s.Publish(&sns.PublishInput{
			TopicArn: aws.String(topicArn),
			Message:  aws.String(string(body)),
		})
		if err != nil {
			return err
		}
		return nil
	}

	return publ, topicArn, nil
}

func SetupQueueSubscription(awsSession *session.Session, snsTopicArn string, sqsArn string) error {
	// AWS will deduplicate subscriptions, so it's safe to do this even if the queue is already subscribed
	s := sns.New(awsSession)
	_, err := s.Subscribe(&sns.SubscribeInput{
		TopicArn: aws.String(snsTopicArn),
		Endpoint: aws.String(sqsArn),
		Protocol: aws.String("sqs"),
		Attributes: map[string]*string{
			"RawMessageDelivery": aws.String("true"),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func assertKey(awsSession *session.Session, logger *zap.Logger) error {
	client := kms.New(awsSession)

	// Check if the key already exists
	aliases, err := client.ListAliases(&kms.ListAliasesInput{})
	if err != nil {
		return fmt.Errorf("failed to list aliases: %w", err)
	}
	for _, a := range aliases.Aliases {
		if *a.AliasName == "alias/sns_and_sqs" {
			return nil
		}
	}

	stsClient := sts.New(awsSession)
	getCallerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get called identity: %w", err)
	}

	policy := fmt.Sprintf(`{
    "Version": "2012-10-17",
    "Id": "key-consolepolicy-3",
    "Statement": [
        {
            "Sid": "Enable IAM User Permissions",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::%s:root"
            },
            "Action": "kms:*",
            "Resource": "*"
        },
        {
            "Sid": "Allow access for Key Administrators",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "%s"
                ]
            },
            "Action": [
                "kms:Create*",
                "kms:Describe*",
                "kms:Enable*",
                "kms:List*",
                "kms:Put*",
                "kms:Update*",
                "kms:Revoke*",
                "kms:Disable*",
                "kms:Get*",
                "kms:Delete*",
                "kms:TagResource",
                "kms:UntagResource",
                "kms:ScheduleKeyDeletion",
                "kms:CancelKeyDeletion"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "sns.amazonaws.com"
            },
            "Action": [
                "kms:GenerateDataKey",
                "kms:Decrypt",
                "kms:Get*"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "sqs.amazonaws.com"
            },
            "Action": [
                "kms:GenerateDataKey",
                "kms:Decrypt"
            ],
            "Resource": "*"
        }
    ]
}`, *getCallerIdentity.Account, *getCallerIdentity.Arn)

	createKey, err := client.CreateKey(&kms.CreateKeyInput{
		Policy: aws.String(policy),
	})

	if err != nil {
		return fmt.Errorf("failed to create kms key: %w", err)
	}

	_, err = client.CreateAlias(&kms.CreateAliasInput{
		AliasName:   aws.String("alias/sns_and_sqs"),
		TargetKeyId: createKey.KeyMetadata.KeyId,
	})
	if err != nil {
		return fmt.Errorf("failed to create kms alias: %w", err)
	}

	logger.Info("created new kms key")

	return nil
}

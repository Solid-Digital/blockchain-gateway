package aws

import (
	"encoding/json"
	"time"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/unchainio/pkg/errors"
)

const (
	QueueSleep          = 3 * time.Second
	MaxNumberOfMessages = 10
	VisibilityTimeout   = 10
	WaitTimeSeconds     = 20
)

func (c *Client) ReceiveMarketplaceNotification() ares.AWSMarketplaceNotificationMessage {
	msg := <-c.awsCh

	return msg
}

func (c *Client) subscribeToSQS(url string) error {

	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(url),
		MaxNumberOfMessages: aws.Int64(MaxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(VisibilityTimeout),
		WaitTimeSeconds:     aws.Int64(WaitTimeSeconds),
	}

	go func() {
		for {
			queueResponse, err := c.sqs.ReceiveMessage(receiveParams)
			if err != nil {
				c.awsCh <- ares.AWSMarketplaceNotificationMessage{
					Error: errors.Wrap(err, "aws marketplace queue error"),
				}
			}

			for _, msg := range queueResponse.Messages {
				var body ares.AWSMarketplaceNotificationMessageBody
				if msg.Body != nil {
					err := json.Unmarshal([]byte(*msg.Body), &body)
					if err != nil {
						c.awsCh <- ares.AWSMarketplaceNotificationMessage{
							Handle: msg.ReceiptHandle,
							Error:  errors.Wrap(err, "could not unmarshal marketplace notification message"),
						}
					}
					c.awsCh <- ares.AWSMarketplaceNotificationMessage{
						Body:   body,
						Handle: msg.ReceiptHandle,
					}
				}
			}
			time.Sleep(QueueSleep)
		}
	}()

	return nil
}

func (c *Client) DeleteSQSMessage(handle *string) error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.cfg.SQSURL),
		ReceiptHandle: handle,
	}
	_, err := c.sqs.DeleteMessage(deleteParams)
	if err != nil {
		return errors.Wrap(err, "could not delete SQS message")
	}

	return nil
}

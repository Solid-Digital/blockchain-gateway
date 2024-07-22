package aws

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/wire"
)

var ClientSet = wire.NewSet(New, wire.Bind(new(ares.AWSClient), new(Client)))

type Client struct {
	mpm   *marketplacemetering.MarketplaceMetering
	mes   *marketplaceentitlementservice.MarketplaceEntitlementService
	sqs   *sqs.SQS
	cfg   *Config
	awsCh chan ares.AWSMarketplaceNotificationMessage
}

func New(cfg *Config) (*Client, error) {

	creds := credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
	}))

	mpm := marketplacemetering.New(sess, &aws.Config{Region: aws.String(cfg.MarketplaceRegion)})

	mes := marketplaceentitlementservice.New(sess, &aws.Config{Region: aws.String(cfg.MarketplaceRegion)})

	s := sqs.New(sess, &aws.Config{MaxRetries: aws.Int(5), Region: aws.String(cfg.QueueRegion)})

	c := &Client{
		mpm:   mpm,
		mes:   mes,
		sqs:   s,
		cfg:   cfg,
		awsCh: make(chan ares.AWSMarketplaceNotificationMessage),
	}

	err := c.subscribeToSQS(cfg.SQSURL)
	if err != nil {
		return nil, err
	}

	return c, nil
}

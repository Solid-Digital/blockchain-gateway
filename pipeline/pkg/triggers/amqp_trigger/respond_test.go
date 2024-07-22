package amqp_trigger_test

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/amqp_trigger"
	"testing"
	"time"
)

func (s *TestSuite) TestTrigger_Respond() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	require.NoError(s.T(), err, "could not connect to test amqp - have you started docker compose?")

	c, err := conn.Channel()
	require.NoError(s.T(), err, "could not create channel for amqp test setup")

	err = c.ExchangeDeclare("test-exchange", "direct", true, false, false, false, nil)
	require.NoError(s.T(), err, "could not declare exchange on mq")

	q, err := c.QueueDeclare("queue-a", true, false, false, false, nil)
	require.NoError(s.T(), err, "could not declare queue")

	err = c.QueueBind(q.Name, "test-key", "test-exchange", false, nil)
	require.NoError(s.T(), err, "could not bind exchange to queue")

	err = c.Publish("test-exchange", "test-key", false, false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(s.helper.BytesFromFile("./testdata/example.json")),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	)

	time.Sleep(2 * time.Second)

	cases := map[string]struct {
		Stub          domain.Stub
		Config        []byte
		Success       bool
		ExpectedLength int
	}{
		"init trigger with valid config triggers as expected": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/config.toml"),
			true,
			0,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			trigger := amqp_trigger.NewTrigger()
			err := trigger.Init(tc.Stub, tc.Config)

			tag, _, _ := trigger.NextMessage()
			err = trigger.Respond(tag, nil, nil)

			if tc.Success {
				// time.Sleep(100 * time.Second)
				q, err := c.QueueInspect("queue-a")
				require.Equal(t, tc.ExpectedLength, q.Messages, "")
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}

	_, err = c.QueuePurge("queue-a", false)
	require.NoError(s.T(), err, "failed to purge messages of queue")
}

package amqp_trigger

import (
	"errors"
	"github.com/streadway/amqp"
)

func (t *Trigger) Respond(tag string, response map[string]interface{}, err error) error {
	d, ok := t.ResponseChannelMap.Load(tag)
	if !ok {
		return errors.New("could not find tag in response map")
	}

	delivery, ok := d.(amqp.Delivery)
	if !ok {
		return errors.New("could not cast to amqp delivery")
	}

	if err != nil {
		err = t.amqpChannel.Nack(delivery.DeliveryTag, false, false)
		if err != nil {
			return err
		}
		return nil
	}

	err = t.amqpChannel.Ack(delivery.DeliveryTag, false)
	if err != nil {
		return err
	}

	return nil
}

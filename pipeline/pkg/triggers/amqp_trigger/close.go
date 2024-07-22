package amqp_trigger

func (t *Trigger) Close() error {
	err := t.amqpChannel.Close()
	if err != nil {
		return err
	}
	err = t.amqpConn.Close()
	if err != nil {
		return err
	}
	return nil
}

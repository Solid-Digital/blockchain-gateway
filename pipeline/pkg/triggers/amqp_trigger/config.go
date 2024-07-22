package amqp_trigger

type Config struct {
	Username   string
	Password   string
	AmqpScheme string
	Domain     string
	Port       string

	QueueName    string
	ConsumerName string // if nil the server will give a name
}

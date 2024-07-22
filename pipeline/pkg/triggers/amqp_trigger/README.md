# AMQP Trigger

Start a simple consumer to get messages from an AMQP message queue.

## Config


Example config:

```toml
username = "guest"
password = "guest"
domain = "amqp://localhost" # make sure to include amqp:// or amqps://
port = "5672"

queueName = "queue-a"
consumerName = "trigger"
```

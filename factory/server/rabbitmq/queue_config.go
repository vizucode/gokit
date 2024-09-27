package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

func setupQueueConfig(ch *amqp.Channel, exchangeName, queueName string) (<-chan amqp.Delivery, error) {
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error in declaring the queue %s", err)
	}
	if err = ch.QueueBind(queue.Name, queue.Name, exchangeName, false, nil); err != nil {
		return nil, fmt.Errorf("error binding queue: %s", err)
	}

	return ch.Consume(
		queue.Name,
		queue.Name, // consumer or channel consumer
		false,      // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no waiting
		nil,        // arguments
	)
}

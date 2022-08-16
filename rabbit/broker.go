package rabbit

import (
	"fmt"

	"github.com/varopxndx/chat/config"

	"github.com/streadway/amqp"
)

// Broker struct
type Broker struct {
	URL    string
	conn   *amqp.Connection
	queues Queues
}

// Queues struct
type Queues struct {
	queueName string
	queue     amqp.Queue
}

// New creates a rabbit broker
func New(conf config.Broker) *Broker {
	conn, err := amqp.Dial(conf.URL)
	if err != nil {
		panic(fmt.Sprintf("error connecting to rabbitmq %v\n", err))
	}
	c, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("error creating rabbitmq channel %v\n", err))
	}
	defer c.Close()
	q, err := c.QueueDeclare(conf.QueueName, true, false, false, false, nil)
	if err != nil {
		panic(fmt.Sprintf("error configuring rabbitmq %v\n", err))
	}
	return &Broker{
		URL:    conf.URL,
		conn:   conn,
		queues: Queues{queueName: conf.QueueName, queue: q},
	}
}

// GetQueueName gets stock queue
func (b Broker) GetQueueName() string {
	return b.queues.queue.Name
}

// Close closes broker connection
func (b Broker) Close() error {
	return b.conn.Close()
}

func (b Broker) channel() (*amqp.Channel, error) {
	if b.conn == nil {
		c, err := amqp.Dial(b.URL)
		if err != nil {
			return nil, err
		}
		b.conn = c
	}
	return b.conn.Channel()
}

// Publish publishes message
func (b Broker) Publish(key string, msg []byte) error {
	channel, err := b.channel()
	if err != nil {
		return err
	}
	return channel.Publish(
		"",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
}

// Consume consumes message
func (b Broker) Consume(key string) (<-chan amqp.Delivery, error) {
	channel, err := b.channel()
	if err != nil {
		return nil, err
	}
	return channel.Consume(
		key,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

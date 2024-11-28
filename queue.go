package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Close()
	Consume()
	Publish()
}

type RabbitQueue struct {
	channel      *amqp.Channel
	conn         *amqp.Connection
	consumeQueue amqp.Queue
	publishQueue amqp.Queue
}

func NewMessageQueue(conStr string) (*RabbitQueue, error) {
	// Connect to MessageQueue Server
	conn, err := amqp.Dial(conStr)
	if err != nil {
		return nil, err
	}
	// Connect to specific channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// Declare the queues
	consumeQueue, err := ch.QueueDeclare(
		"rawForms", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	publishQueue, err := ch.QueueDeclare(
		"cleanForms", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	// Return new MessageQueue instance
	return &RabbitQueue{
		channel:      ch,
		conn:         conn,
		consumeQueue: consumeQueue,
		publishQueue: publishQueue,
	}, nil
}

func (mq *RabbitQueue) Close() {
	// Close the connection
	mq.conn.Close()
	mq.channel.Close()
}

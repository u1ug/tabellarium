package message_queue

import "github.com/streadway/amqp"

type Delivery struct {
	Exchange   string
	Key        string
	Mandatory  bool
	Immediate  bool
	Publishing amqp.Publishing
}

package message_queue

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

type Producer struct {
	Terminate chan bool
	Deliver   chan Delivery
	conn      *amqp.Connection
	ch        *amqp.Channel
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewProducer(buffering uint) *Producer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Producer{
		Terminate: make(chan bool),
		Deliver:   make(chan Delivery, buffering),
		ctx:       ctx,
		cancel:    cancel,
	}
}
func (p *Producer) Connect(username string, password string, hostname string, port int) error {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, hostname, port)
	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	p.conn = conn
	p.ch, err = p.conn.Channel()
	return err
}

func (p *Producer) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (*amqp.Queue, error) {
	q, err := p.ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	return &q, err
}

func (p *Producer) Listen() error {
	for {
		select {
		case delivery := <-p.Deliver:
			err := p.ch.Publish(delivery.Exchange, delivery.Key, delivery.Mandatory, delivery.Immediate, delivery.Publishing)
			if err != nil {
				fmt.Println(err) // for debug only, will be removed in future
			}
		case <-p.ctx.Done():
			return p.terminate()
		}
	}
}

func (p *Producer) terminate() error {
	var err error
	close(p.Deliver)
	close(p.Terminate)
	err = p.ch.Close()
	if err == nil {
		err = p.conn.Close()
	}
	return err
}

func (p *Producer) Close() error {
	p.cancel()
	return nil
}

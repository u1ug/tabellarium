package message_queue

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"golang.org/x/sync/semaphore"
)

type ConsumerHandler func(delivery amqp.Delivery) error

type Consumer struct {
	conn                  *amqp.Connection
	ch                    *amqp.Channel
	queue                 string
	maxConcurrentHandlers int64
	sem                   *semaphore.Weighted
	ctx                   context.Context
	cancel                context.CancelFunc
}

func NewConsumer(maxConcurrentHandlers int64) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		maxConcurrentHandlers: maxConcurrentHandlers,
		sem:                   semaphore.NewWeighted(maxConcurrentHandlers),
		ctx:                   ctx,
		cancel:                cancel,
	}
}

func (c *Consumer) Connect(username, password, hostname string, port int) error {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, hostname, port)
	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.ch, err = c.conn.Channel()
	return err
}

func (c *Consumer) Listen(queue string, handler ConsumerHandler) error {
	c.queue = queue

	messages, err := c.ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				return nil
			}

			if err := c.sem.Acquire(c.ctx, 1); err != nil {
				fmt.Println("Failed to acquire semaphore:", err)
				continue
			}

			go func(m amqp.Delivery) {
				defer c.sem.Release(1)

				err := handler(m)
				if err != nil {
					fmt.Println(err)
				}
			}(msg)

		case <-c.ctx.Done():
			return c.ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	for i := int64(0); i < c.maxConcurrentHandlers; i++ {
		err := c.sem.Acquire(c.ctx, 1)
		if err != nil {
			return err
		}
	}

	if err := c.ch.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	c.cancel()
	return nil
}

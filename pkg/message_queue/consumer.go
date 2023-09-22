package message_queue

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/semaphore"
	"tabellarium/pkg/logging"
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
	logger                *logging.Logger
}

func NewConsumer(maxConcurrentHandlers int64) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		maxConcurrentHandlers: maxConcurrentHandlers,
		sem:                   semaphore.NewWeighted(maxConcurrentHandlers),
		ctx:                   ctx,
		cancel:                cancel,
		logger:                logging.GetLogger(),
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
		c.logger.Errorln(err)
		return err
	}
	c.logger.Debugln("begin consumer listening loop")
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				c.logger.Debugln("no messages channel opened")
				return nil
			}

			if err := c.sem.Acquire(c.ctx, 1); err != nil {
				c.logger.Debugln(err)
				continue
			}

			go func(m amqp.Delivery) {
				defer c.sem.Release(1)
				err := handler(m)
				if err != nil {
					c.logger.Errorln(err)
				}
			}(msg)

		case <-c.ctx.Done():
			c.logger.Infoln("finishing consumer listening")
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

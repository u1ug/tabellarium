package message_queue

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"tabellarium/pkg/logging"
)

// Producer is the entity for message publishing.
type Producer struct {
	Deliver   chan Delivery // Used for message to queue transfer
	connected bool          // Use to prevent Listen() call without connection
	// AMQP entities
	conn *amqp.Connection
	ch   *amqp.Channel

	settings ProdSettings // Producer settings
	//Concurrent controlling entities
	ctx    context.Context
	cancel context.CancelFunc
	once   sync.Once

	logger *logging.Logger // Logrus logger
}

// ProdSettings are fields required for Producer init.
type ProdSettings struct {
	Logger    *logging.Logger
	Buffering uint // Deliver channel buffering.
	// AMQP connection details.
	MQUser     string
	MQPassword string
	MQAddress  string
	MQPort     int
	// AMQP queue parameters.
	Queue      string //Queue name
	Durable    bool   // If set to true, the queue will survive server restarts.
	AutoDelete bool   // If set to true, the queue will be deleted when there are no more consumers.
	Exclusive  bool   // If set to true, the queue can only be accessed by the current connection and will be deleted when that connection closes.
	NoWait     bool   // If set to true, the queue declaration will not wait for a reply from the server.
}

// NewProducer initializes new Producer
func NewProducer(s ProdSettings) *Producer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Producer{
		Deliver:   make(chan Delivery, s.Buffering),
		connected: false,
		settings:  s,
		ctx:       ctx,
		cancel:    cancel,
		logger:    s.Logger,
	}
}

func (p *Producer) Connect() error {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/", p.settings.MQUser, p.settings.MQPassword, p.settings.MQAddress, p.settings.MQPort)
	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(p.settings.Queue, p.settings.Durable, p.settings.AutoDelete, p.settings.Exclusive, p.settings.NoWait, nil)
	if err != nil {
		return err
	}
	p.logger.Debugln(q)
	p.conn = conn
	p.ch = ch
	p.connected = true
	return err
}

// Listen begin producer listening loop. Publishes messages via Deliver chan, can be stopped with p.Close().
func (p *Producer) Listen() error {
	if !p.connected || p.ch == nil || p.conn == nil {
		return errors.New("producer is not connected")
	}
	p.logger.Debugln("begin listening loop")
	for {
		select {
		case message := <-p.Deliver:
			pub := FormPublishing(message)
			err := p.ch.PublishWithContext(p.ctx, message.Exchange, message.Key, message.Mandatory, message.Immediate, pub)
			if err != nil {
				p.logger.Errorln(err)
			}
		case <-p.ctx.Done():
			p.logger.Println("exiting listening loop")
			return p.Close()
		}
	}
}

// Close handles producer listening finish.
func (p *Producer) Close() error {
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
	p.once.Do(func() {
		close(p.Deliver)
	})
	if err := p.ch.Close(); err != nil {
		return err
	}
	if err := p.conn.Close(); err != nil {
		return err
	}

	return nil
}

package service

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"tabellarium/internal/entities"
	"tabellarium/pkg/logging"
	"tabellarium/pkg/message_queue"
)

// Notificator defines the interface for notification operations.
type Notificator interface {
	Listen() error
	push(n entities.Notification) error
	Close() error
}

// NotificationService implements the Notificator interface.
type NotificationService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	consumer *message_queue.Consumer
	settings Settings
	logger   *logging.Logger
}

// Settings holds the configuration for the NotificationService.
type Settings struct {
	Queue       string
	MaxHandlers int64
	MQUsername  string
	MQPassword  string
	MQHostName  string
	MQPort      int
}

// NewNotificationService initializes a new NotificationService with the given settings.
func NewNotificationService(s Settings) *NotificationService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationService{
		ctx:      ctx,
		cancel:   cancel,
		consumer: message_queue.NewConsumer(s.MaxHandlers),
		settings: s,
		logger:   logging.GetLogger(),
	}
}

// Listen starts the listening loop for the NotificationService.
func (s NotificationService) Listen() error {
	s.logger.Infoln("listening for notifications")
	if err := s.consumer.Connect(s.settings.MQUsername, s.settings.MQPassword, s.settings.MQHostName, s.settings.MQPort); err != nil {
		s.logger.Fatalf("can not connect to the queue: %v\n", err)
		return err
	}
	return s.consumer.Listen(s.settings.Queue, func(delivery amqp.Delivery) error {
		s.logger.Debugf("received notification job: %v\n", delivery)
		return nil
	})
}

// Push sends a notification.
func (s NotificationService) push(n entities.Notification) error {
	panic("not implemented")
}

// Close shuts down the NotificationService.
func (s NotificationService) Close() error {
	if err := s.consumer.Close(); err != nil {
		return err
	}
	s.cancel()
	return nil
}

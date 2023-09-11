package service

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"tabellarium/pkg/message_queue"
)

// Notificator defines the interface for notification operations.
type Notificator interface {
	Listen() error
	Push(n Notification) error
	Close() error
}

// NotificationService implements the Notificator interface.
type NotificationService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	consumer *message_queue.Consumer
	settings Settings
}

// Settings holds the configuration for the NotificationService.
type Settings struct {
	ServiceHostname string
	Queue           string
	MaxHandlers     int64
	MQUsername      string
	MQPassword      string
	MQHostName      string
	MQPort          int
}

// NewNotificationService initializes a new NotificationService with the given settings.
func NewNotificationService(s Settings) *NotificationService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationService{
		ctx:      ctx,
		cancel:   cancel,
		consumer: message_queue.NewConsumer(s.MaxHandlers),
		settings: s,
	}
}

// Listen starts the listening loop for the NotificationService.
func (s NotificationService) Listen() error {
	fmt.Println("begin listening loop")
	if err := s.consumer.Connect(s.settings.MQUsername, s.settings.MQPassword, s.settings.MQHostName, s.settings.MQPort); err != nil {
		return err
	}
	return s.consumer.Listen(s.settings.Queue, func(delivery amqp.Delivery) error {
		fmt.Println(delivery)
		return nil
	})
}

// Push sends a notification.
func (s NotificationService) Push(n Notification) error {
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

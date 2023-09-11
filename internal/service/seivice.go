package service

import (
	"context"
	"tabellarium/pkg/message_queue"
)

type Notificator interface {
	Listen(hostname string) error
	Push(n Notification) error
	Close() error
}

type NotificationService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	consumer *message_queue.Consumer
}

func NewNotificationService(maxHandlers int64) *NotificationService {
	ctx, cancel := context.WithCancel(context.Background())
	cons := message_queue.NewConsumer(maxHandlers)
	n := &NotificationService{
		ctx:      ctx,
		cancel:   cancel,
		consumer: cons,
	}
	return n
}

func (s NotificationService) Listen(hostname string) error {
	panic(1)
}

func (s NotificationService) Push(n Notification) error {
	panic(1)
}

func (s NotificationService) Close() error {
	err := s.consumer.Close()
	if err != nil {
		return err
	}
	s.cancel()
	return err
}

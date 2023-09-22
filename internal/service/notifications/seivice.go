package notifications

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"tabellarium/internal/entities"
	"tabellarium/internal/infrastructure/expo_api"
	"tabellarium/internal/infrastructure/storage"
	"tabellarium/pkg/logging"
	"tabellarium/pkg/message_queue"
)

// NotificationService implements the Notificator interface.
type NotificationService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	consumer *message_queue.Consumer
	store    *storage.DeviceStorage
	settings Settings
	logger   *logging.Logger
}

// Settings holds the configuration for the NotificationService.
type Settings struct {
	// Redis configuration.
	RedisAddress  string
	RedisPort     uint
	RedisDB       int
	RedisPassword string
	RedisTTL      int
	// Rabbitmq settings
	Queue       string
	MaxHandlers int64 // Maximum amount of message handlers called
	MQUsername  string
	MQPassword  string
	MQHostName  string
	MQPort      int
}

// NewNotificationService initializes a new NotificationService with the given settings.
func NewNotificationService(s Settings) *NotificationService {
	store := storage.NewDeviceStorage(storage.Settings{
		Address:  s.RedisAddress,
		Port:     s.RedisPort,
		DB:       s.RedisDB,
		Password: s.RedisPassword,
		TTL:      s.RedisTTL,
	})
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationService{
		ctx:      ctx,
		cancel:   cancel,
		consumer: message_queue.NewConsumer(s.MaxHandlers),
		store:    store,
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
		notification, err := entities.ParseNotification(delivery.Body)
		if err != nil {
			s.logger.Errorln(err)
			return err
		}
		if err != nil {
			s.logger.Errorln(err)
			return err
		}
		receivers := s.store.FetchUsers(notification.To)
		notification.To = receivers
		err, resp := expo_api.SendNotification(notification)
		s.logger.Debugln(string(resp))
		if err != nil {
			s.logger.Errorln(err)
			return err
		}
		return nil
	})
}

// Close shuts down the NotificationService.
func (s NotificationService) Close() error {
	if err := s.consumer.Close(); err != nil {
		return err
	}
	s.cancel()
	return nil
}

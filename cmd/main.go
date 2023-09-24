package main

import (
	"sync"
	"tabellarium/internal/config"
	"tabellarium/internal/service/notifications"
	"tabellarium/internal/service/user_service"
)

func initNotifications(c *config.Config) error {
	serv := notifications.NewNotificationService(notifications.Settings{
		RedisAddress:  c.Redis.RedisAddress,
		RedisPort:     c.Redis.RedisPort,
		RedisDB:       c.Redis.RedisDB,
		RedisPassword: c.Redis.RedisPassword,
		RedisTTL:      c.Redis.RedisTTL,
		Queue:         c.RabbitMQ.Queue,
		MaxHandlers:   c.RabbitMQ.MaxHandlers,
		MQUsername:    c.RabbitMQ.MQUsername,
		MQPassword:    c.RabbitMQ.MQPassword,
		MQHostName:    c.RabbitMQ.MQHostName,
		MQPort:        c.RabbitMQ.MQPort,
	})
	return serv.Listen()
}

func initAPI(c *config.Config) error {
	serv := user_service.NewUserService(user_service.Settings{
		Address:       c.Service.Address,
		Port:          c.Service.Port,
		RedisAddress:  c.Redis.RedisAddress,
		RedisPort:     c.Redis.RedisPort,
		RedisDB:       c.Redis.RedisDB,
		RedisPassword: c.Redis.RedisPassword,
		RedisTTL:      c.Redis.RedisTTL,
	})
	return serv.Listen()
}

func main() {
	cfg, err := config.GetConfig("./config.json")
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := initNotifications(cfg)
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := initAPI(cfg)
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()
	wg.Wait()
}

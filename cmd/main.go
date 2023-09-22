package main

import (
	"sync"
	"tabellarium/internal/service/notifications"
	"tabellarium/internal/service/user_service"
)

func initNotifications() error {
	serv := notifications.NewNotificationService(notifications.Settings{
		RedisAddress:  "localhost",
		RedisPort:     6379,
		RedisDB:       0,
		RedisPassword: "",
		RedisTTL:      43200,
		Queue:         "Notifications",
		MaxHandlers:   10000,
		MQUsername:    "guest",
		MQPassword:    "guest",
		MQHostName:    "localhost",
		MQPort:        5672,
	})
	return serv.Listen()
}

func initAPI() error {
	serv := user_service.NewUserService(user_service.Settings{
		Address:       "localhost",
		Port:          8080,
		RedisAddress:  "localhost",
		RedisPort:     6379,
		RedisDB:       0,
		RedisPassword: "",
		RedisTTL:      43200,
	})
	return serv.Listen()
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := initNotifications()
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := initAPI()
		defer wg.Done()
		if err != nil {
			panic(err)
		}
	}()
	wg.Wait()
}

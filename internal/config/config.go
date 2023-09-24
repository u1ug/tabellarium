package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Service struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"service"`
	RabbitMQ struct {
		Queue       string `json:"queue"`
		MaxHandlers int64  `json:"maxHandlers"`
		MQUsername  string `json:"MQUsername"`
		MQPassword  string `json:"MQPassword"`
		MQHostName  string `json:"MQHostName"`
		MQPort      int    `json:"MQPort"`
	} `json:"rabbitMQ"`
	Redis struct {
		RedisAddress  string `json:"RedisAddress"`
		RedisPort     uint   `json:"RedisPort"`
		RedisDB       int    `json:"RedisDB"`
		RedisPassword string `json:"RedisPassword"`
		RedisTTL      int    `json:"RedisTTL"`
	} `json:"redis"`
}

func GetConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

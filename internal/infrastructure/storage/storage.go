package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type DeviceStorage struct {
	ctx context.Context
	rdb *redis.Client
	ttl int
}

type Settings struct {
	Address  string
	Port     uint
	DB       int
	Password string
	TTL      int
}

func NewDeviceStorage(s Settings) *DeviceStorage {
	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: s.Password, DB: s.DB})
	return &DeviceStorage{
		ctx: ctx,
		rdb: rdb,
		ttl: s.TTL,
	}
}

func (d *DeviceStorage) RegisterDevice(userID string, token string) error {
	offset := time.Duration(d.ttl) * time.Second
	currTime := time.Now()
	expirationTimestamp := currTime.Add(offset).Unix()
	err := d.rdb.HSet(d.ctx, userID, token, strconv.FormatInt(expirationTimestamp, 10)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (d *DeviceStorage) GetUserDevices(userID string) ([]string, error) {
	tokens, err := d.rdb.HGetAll(d.ctx, userID).Result()
	if err != nil {
		return nil, err
	}
	// Filter out expired IDs
	var validIDs []string
	currentTime := time.Now().Unix()
	for token, timestamp := range tokens {
		expirationTimestamp, _ := strconv.ParseInt(timestamp, 10, 64)
		if currentTime <= expirationTimestamp {
			validIDs = append(validIDs, token)
		}
	}
	return validIDs, nil
}

// FetchUsers returns an array with tokens of all user provided.
func (d *DeviceStorage) FetchUsers(users []string) []string {
	var tokens []string
	for _, userID := range users {
		devices, err := d.GetUserDevices(userID)
		if err != nil {
			continue
		}
		fmt.Println(devices)
		tokens = append(tokens, devices...)
	}
	return tokens
}

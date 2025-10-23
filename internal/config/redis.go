package config

import (
	"fmt"

	fiberRedis "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
}

func NewRedisClient(host string, port int, password string, db int) *Client {
	addr := fmt.Sprintf("%s:%d", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{RDB: rdb}
}

func NewFiberRedisStorage(host string, port int, password string, db int) *fiberRedis.Storage {
	storage := fiberRedis.New(fiberRedis.Config{
		Host:     host,
		Port:     port,
		Password: password,
		Database: db,
	})
	return storage
}

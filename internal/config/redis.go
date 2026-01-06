package config

import (
	"crypto/tls"
	"fmt"

	fiberRedis "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
}

func NewRedisClient(host string, port int, password string, db int, tlsEnabled bool) *Client {
	addr := fmt.Sprintf("%s:%d", host, port)

	var tlsConfig *tls.Config
	if tlsEnabled {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  password,
		DB:        db,
		TLSConfig: tlsConfig,
	})

	return &Client{RDB: rdb}
}

func NewFiberRedisStorage(host string, port int, password string, db int, tlsEnabled bool) *fiberRedis.Storage {
	var tlsConfig *tls.Config
	if tlsEnabled {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	storage := fiberRedis.New(fiberRedis.Config{
		Host:      host,
		Port:      port,
		Password:  password,
		Database:  db,
		TLSConfig: tlsConfig,
	})
	return storage
}

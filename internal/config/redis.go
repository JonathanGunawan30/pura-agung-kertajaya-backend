package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
}

func NewRedisClient(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	return &Client{RDB: rdb}
}

func (r *Client) BlacklistToken(ctx context.Context, jti string) error {
	return r.RDB.Set(ctx, "blacklist:"+jti, true, 24*time.Hour).Err()
}

func (r *Client) IsTokenBlacklisted(ctx context.Context, jti string) bool {
	exists, _ := r.RDB.Exists(ctx, "blacklist:"+jti).Result()
	return exists == 1
}

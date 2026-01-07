package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(addr, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{ //for configuration setup
		Addr:     addr,     //Redis server host + port
		Password: password, //Redis password (empty string if no password)
		DB:       db,       //Database number (0 default)
	})

	return &Client{rdb: rdb}
}
func (c *Client) SetToken(token string, userID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return c.rdb.Set(ctx, token, userID, time.Hour).Err()
}

func (c *Client) ValidateToken(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := c.rdb.Get(ctx, token).Result()
	if err == redis.Nil {
		return false, nil
	}
	return err == nil, err
}

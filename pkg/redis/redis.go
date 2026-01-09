package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

func New(
	addr string,
	password string,
	db int,
	ttl time.Duration,
) (*Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		rdb: rdb,
		ttl: ttl,
	}, nil
}

// ------------------------------------
// Token Store
// ------------------------------------
func (c *Client) SetToken(token string, userID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := c.tokenKey(token)

	return c.rdb.Set(ctx, key, userID, c.ttl).Err()
}

// ------------------------------------
// Token Validation
// ------------------------------------
func (c *Client) ValidateToken(token string) (uint, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := c.tokenKey(token)

	val, err := c.rdb.Get(ctx, key).Uint64()
	if err == redis.Nil {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}

	return uint(val), true, nil
}

// ------------------------------------
// Token Revoke (Logout)
// ------------------------------------
func (c *Client) RevokeToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return c.rdb.Del(ctx, c.tokenKey(token)).Err()
}

// ------------------------------------
// Helpers
// ------------------------------------
func (c *Client) tokenKey(token string) string {
	return fmt.Sprintf("auth:token:%s", token)
}

// ------------------------------------
// Close
// ------------------------------------
func (c *Client) Close() error {
	return c.rdb.Close()
}

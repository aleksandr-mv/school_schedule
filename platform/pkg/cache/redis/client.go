package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/cache"
)

type client struct {
	rdb     redis.Cmdable
	logger  Logger
	timeout time.Duration
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// NewClient создаёт обёртку над go-redis клиентом (ClusterClient или Client)
func NewClient(rdb redis.Cmdable, logger Logger, timeout time.Duration) cache.RedisClient {
	return &client{
		rdb:     rdb,
		logger:  logger,
		timeout: timeout,
	}
}

func (c *client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *client) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return data, err
}

func (c *client) Del(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Del(ctx, key).Err()
}

func (c *client) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Ping(ctx).Err()
}

func (c *client) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.HSet(ctx, key, values).Err()
}

func (c *client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.HGetAll(ctx, key).Result()
}

func (c *client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Expire(ctx, key, ttl).Err()
}

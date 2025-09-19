package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/cache"
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

func (c *client) Set(ctx context.Context, key string, value any) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Set(ctx, key, value, 0).Err()
}

func (c *client) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *client) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return data, err
}

func (c *client) HashSet(ctx context.Context, key string, values any) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.HSet(ctx, key, values).Err()
}

func (c *client) HGetAll(ctx context.Context, key string) ([]any, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	m, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res := make([]any, 0, len(m)*2)
	for k, v := range m {
		res = append(res, k, v)
	}
	return res, nil
}

func (c *client) Del(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Del(ctx, key).Err()
}

func (c *client) Exists(ctx context.Context, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Expire(ctx, key, expiration).Err()
}

func (c *client) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.Ping(ctx).Err()
}

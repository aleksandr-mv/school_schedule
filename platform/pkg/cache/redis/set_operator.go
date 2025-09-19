package redis

import (
	"context"
)

// SetOperator методы для работы с Redis Sets
// Выделены в отдельный файл для лучшей организации кода

func (c *client) SAdd(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.SAdd(ctx, key, value).Err()
}

func (c *client) SRem(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.SRem(ctx, key, value).Err()
}

func (c *client) SIsMember(ctx context.Context, key, value string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.SIsMember(ctx, key, value).Result()
}

func (c *client) SMembers(ctx context.Context, key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.rdb.SMembers(ctx, key).Result()
}

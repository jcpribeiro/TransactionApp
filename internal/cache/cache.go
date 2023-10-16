package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=cache_mock.go -package=$GOPACKAGE

type Cache interface {
	Get(ctx context.Context, key string, value interface{})
	Set(ctx context.Context, key string, value interface{}, expirationTime time.Duration) error
}

type cacheImpl struct {
	redis *redis.Client
	log   logrus.Logger
}

func NewCache(redis *redis.Client, log logrus.Logger) Cache {
	return &cacheImpl{
		redis: redis,
	}
}

func (c *cacheImpl) Get(ctx context.Context, key string, value interface{}) {
	cmd := c.redis.Get(ctx, key)
	if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
		c.log.Error(cmd.Err())
	}

	err := json.Unmarshal([]byte(cmd.Val()), &value)
	if err != nil {
		c.log.Error(fmt.Errorf("failed to unmarshal value: %w", err))
	}
}

func (c *cacheImpl) Set(ctx context.Context, key string, value interface{}, expirationTime time.Duration) error {
	b, err := marshalBinary(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	cmd := c.redis.Set(ctx, key, b, expirationTime)
	if cmd.Err() != nil {
		return fmt.Errorf("failed to set cache: %w", cmd.Err())
	}

	return nil
}

func marshalBinary(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

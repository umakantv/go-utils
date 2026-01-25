package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache implements Redis-based caching
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// newRedisCache creates a new Redis cache
func newRedisCache(config Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis with TTL
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, data, ttl).Err()
}

// Get retrieves a value from Redis
func (c *RedisCache) Get(key string) (interface{}, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	// Try to unmarshal as JSON first
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// If not JSON, return as string
		return val, nil
	}

	return result, nil
}

// Delete removes a key from Redis
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (c *RedisCache) Exists(key string) bool {
	count, err := c.client.Exists(c.ctx, key).Result()
	return err == nil && count > 0
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}
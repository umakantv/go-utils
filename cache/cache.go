package cache

import (
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Exists(key string) bool
	Close() error
}

// Config holds cache configuration
type Config struct {
	Type     string // "memory" or "redis"
	RedisAddr string // Redis server address (e.g., "localhost:6379")
	RedisPassword string // Redis password (optional)
	RedisDB   int    // Redis database number
}

// New creates a new cache instance based on the configuration
func New(config Config) (Cache, error) {
	switch config.Type {
	case "redis":
		return newRedisCache(config)
	case "memory":
		return newMemoryCache(), nil
	default:
		return newMemoryCache(), nil // default to memory
	}
}
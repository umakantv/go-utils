package cache

import (
	"sync"
	"time"
)

// item represents a cached item with expiration
type item struct {
	value      interface{}
	expiration int64 // unix timestamp
}

// MemoryCache implements in-memory caching
type MemoryCache struct {
	items map[string]*item
	mutex sync.RWMutex
}

// newMemoryCache creates a new in-memory cache
func newMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*item),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Set stores a value in the cache with TTL
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).Unix()
	}

	c.items[key] = &item{
		value:      value,
		expiration: expiration,
	}

	return nil
}

// Get retrieves a value from the cache
func (c *MemoryCache) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, ErrKeyNotFound
	}

	// Check if expired
	if item.expiration > 0 && time.Now().Unix() > item.expiration {
		// Item expired, delete it
		delete(c.items, key)
		return nil, ErrKeyNotFound
	}

	return item.value, nil
}

// Delete removes a key from the cache
func (c *MemoryCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
	return nil
}

// Exists checks if a key exists in the cache
func (c *MemoryCache) Exists(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	// Check if expired
	if item.expiration > 0 && time.Now().Unix() > item.expiration {
		return false
	}

	return true
}

// Close is a no-op for memory cache
func (c *MemoryCache) Close() error {
	return nil
}

// cleanup runs in a goroutine to remove expired items
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now().Unix()
		for key, item := range c.items {
			if item.expiration > 0 && now > item.expiration {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}
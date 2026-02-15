# Cache Module

The cache module provides a unified interface for caching with support for both in-memory and Redis backends.

## Dependencies

- `github.com/go-redis/redis/v8` for Redis support

## Configuration

Configure caching using the `Config` struct:

```go
type Config struct {
    Type         string // "memory" or "redis"
    RedisAddr    string // Redis server address (e.g., "localhost:6379")
    RedisPassword string // Redis password (optional)
    RedisDB      int    // Redis database number
}
```

## Creating a Cache

```go
// In-memory cache
memoryCache, err := cache.New(cache.Config{
    Type: "memory",
})

// Redis cache
redisCache, err := cache.New(cache.Config{
    Type: "redis",
    RedisAddr: "localhost:6379",
    RedisPassword: "mypassword", // optional
    RedisDB: 0,
})
```

## Basic Operations

### Set a value with TTL
```go
err := cache.Set("user:123", map[string]string{
    "name": "John",
    "email": "john@example.com",
}, 10*time.Minute)
```

### Get a value
```go
user, err := cache.Get("user:123")
if err == cache.ErrKeyNotFound {
    // Handle cache miss
}
```

### Check if key exists
```go
if cache.Exists("user:123") {
    // Key exists
}
```

### Delete a key
```go
err := cache.Delete("user:123")
```

## Usage Examples

### In-Memory Cache
```go
config := cache.Config{Type: "memory"}
cache, err := cache.New(config)
if err != nil {
    log.Fatal(err)
}
defer cache.Close()

// Store data
err = cache.Set("session:abc123", "user123", 30*time.Minute)

// Retrieve data
session, err := cache.Get("session:abc123")
```

### Redis Cache
```go
config := cache.Config{
    Type: "redis",
    RedisAddr: "localhost:6379",
    RedisPassword: "mypassword",
    RedisDB: 0,
}
cache, err := cache.New(config)
if err != nil {
    log.Fatal(err)
}
defer cache.Close()

// Store complex data
user := map[string]interface{}{
    "id": 123,
    "name": "John Doe",
    "roles": []string{"admin", "user"},
}
err = cache.Set("user:123", user, time.Hour)

// Retrieve and type assert
cachedUser, err := cache.Get("user:123")
if userData, ok := cachedUser.(map[string]interface{}); ok {
    fmt.Println("User name:", userData["name"])
}
```

### Integration with HTTP Server
```go
func getUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    userID := httpserver.GetRoutePath(ctx) // Assuming path param extraction

    // Try cache first
    cacheKey := "user:" + userID
    if cachedUser, err := cache.Get(cacheKey); err == nil {
        // Return cached data
        json.NewEncoder(w).Encode(cachedUser)
        return
    }

    // Fetch from database
    user := fetchUserFromDB(userID)

    // Cache for future requests
    cache.Set(cacheKey, user, 10*time.Minute)

    json.NewEncoder(w).Encode(user)
}
```

## Cache Interface

The `Cache` interface provides consistent operations across implementations:

```go
type Cache interface {
    Set(key string, value interface{}, ttl time.Duration) error
    Get(key string) (interface{}, error)
    Delete(key string) error
    Exists(key string) bool
    Close() error
}
```

## Data Serialization

- **In-Memory**: Stores Go values directly
- **Redis**: Serializes to JSON for storage, deserializes on retrieval
- **TTL**: Time-to-live supported for both backends
- **Cleanup**: In-memory cache automatically removes expired items

## Error Handling

- `ErrKeyNotFound`: Returned when key doesn't exist or has expired
- Connection errors are returned for Redis operations
- Serialization errors are propagated for complex types

## Best Practices

1. Use meaningful key names (e.g., "user:{id}", "session:{token}")
2. Set appropriate TTL values based on data freshness requirements
3. Handle cache misses gracefully by falling back to source data
4. Use cache for expensive operations (database queries, API calls)
5. Monitor cache hit rates and adjust TTL accordingly
6. Close cache connections when application shuts down
7. Consider cache size limits for memory cache in production
8. Use Redis for distributed caching across multiple instances

## Performance Considerations

- **In-Memory**: Fastest, but limited to single instance
- **Redis**: Slightly slower due to serialization, but distributable
- **TTL**: Helps prevent memory leaks and ensures data freshness
- **Cleanup**: In-memory cache runs periodic cleanup every 5 minutes
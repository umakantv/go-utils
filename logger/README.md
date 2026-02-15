# Logger Module

The logger module provides structured logging using [Uber Zap](https://github.com/uber-go/zap) with configurable output and caller information.

## Configuration

Configure the logger using `LoggerConfig`:

```go
type LoggerConfig struct {
    CallerKey  string // Field name for caller info (default: "caller")
    TimeKey    string // Field name for timestamp (default: "time")
    CallerSkip int    // Number of callers to skip
}
```

## Initialization

Initialize the logger once at application startup:

```go
config := LoggerConfig{
    CallerKey:  "file",
    TimeKey:    "timestamp",
    CallerSkip: 1,
}

logger.Init(config)
```

## Logging Levels

Available logging levels:

- `DebugLevel` (-1)
- `InfoLevel` (0) - default
- `WarnLevel` (1)
- `ErrorLevel` (2)
- `DPanicLevel` (3)
- `PanicLevel` (4)
- `FatalLevel` (5)

## Basic Logging

```go
logger.Info("Application started")
logger.Debug("Processing user", logger.String("user_id", "123"))
logger.Error("Database connection failed", logger.Error(err))
```

## Field Functions

Pre-built field functions for common types:

```go
logger.Info("User login",
    logger.String("username", "john"),
    logger.Int("attempts", 3),
    logger.Bool("success", true),
    logger.Duration("duration", time.Since(start)),
    logger.Any("metadata", map[string]interface{}{"ip": "192.168.1.1"}))
```

## Custom Fields

```go
logger.Info("Custom log",
    zap.String("custom_field", "value"),
    zap.Int64("timestamp", time.Now().Unix()))
```

## Logger Types

The module exports Zap types for advanced usage:

```go
import "go.uber.org/zap/zapcore"

level := logger.Level.InfoLevel
```

## Best Practices

1. Initialize logger at application startup
2. Use structured fields instead of string formatting
3. Choose appropriate log levels
4. Include relevant context in log messages
5. Use `Error` level for errors, `Info` for important events, `Debug` for development

## Integration with HTTP Server

The logger integrates seamlessly with the httpserver module for automatic request logging.
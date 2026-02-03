# User Service Example

A complete example of a REST API service demonstrating all the go-utils packages.

## Features

- **Database**: SQLite with local file storage
- **HTTP Server**: Standardized routing with authentication
- **Cache**: In-memory caching for performance
- **Logger**: Structured logging with custom format
- **Errors**: Standardized error responses

## API Endpoints

### Public Endpoints
- `GET /health` - Health check (no auth required)

### Protected Endpoints (Bearer token required)
- `GET /users` - List all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

## Authentication

All API endpoints (except health check) require Bearer token authentication:

```
Authorization: Bearer secret-token
```

## Request/Response Examples

### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Authorization: Bearer secret-token" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

### Get Users
```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer secret-token"
```

### Get User by ID
```bash
curl -X GET http://localhost:8080/users/1 \
  -H "Authorization: Bearer secret-token"
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Authorization: Bearer secret-token" \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "email": "jane@example.com"}'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/1 \
  -H "Authorization: Bearer secret-token"
```

## Running the Service

1. **Build and run (from project root):**
   ```bash
   cd examples/user_service
   go run main.go
   ```

2. **Or create a separate module:**
   ```bash
   cd examples/user_service
   go mod init user-service
   go mod edit -replace=github.com/umakantv/go-utils=../..
   go mod tidy
   go run main.go
   ```

2. **Check health:**
   ```bash
   curl http://localhost:8080/health
   ```

3. **Test API with authentication:**
   ```bash
   # Create a user
   curl -X POST http://localhost:8080/users \
     -H "Authorization: Bearer secret-token" \
     -H "Content-Type: application/json" \
     -d '{"name": "Test User", "email": "test@example.com"}'

   # List users
   curl -X GET http://localhost:8080/users \
     -H "Authorization: Bearer secret-token"
   ```

## Database

The service uses SQLite with a local file `./user_service.db`. The schema includes:

- `users` table with id, name, email, timestamps
- Index on email for performance

## Logging

All requests are logged with the format:
```
2023-12-01 10:30:45 - GetUser - GET - /users/{id} - client:user-service-client - User retrieved successfully
```

## Caching

- User list cached for 5 minutes
- Individual users cached for 10 minutes
- Cache automatically cleared on create/update/delete operations

## Error Responses

Standardized error responses using the errs package:

```json
{
  "Code": 404,
  "Message": "User not found"
}
```

## Architecture

```
main.go          - Service entry point
handlers/        - Request handlers
models/          - Data models
db/              - Database schema
```

This example demonstrates enterprise-ready patterns for building microservices with the go-utils packages.
# Mini HTTP Server in Go

A lightweight, optimized HTTP server implementation in Go, built from scratch without using the standard library's HTTP package. This project demonstrates core HTTP server concepts including request parsing, routing, and response handling.

## Features

### HTTP Protocol Implementation

- Custom HTTP/1.1 protocol implementation
- Request parsing and header handling
- Response generation with proper status codes
- Connection management with timeouts

### Routing System

- Path-based routing with regex pattern matching
- Support for URL parameters (e.g., `/users/:id`)
- HTTP method-based routing (GET, POST)
- Custom 404 Not Found handler

### API Endpoints

- `/v1/healthcheck` - Server health check endpoint
- `/v1/echo/:str` - Echo service that returns the URL parameter
- `/v1/user-agent` - Returns the User-Agent header from the request

### Error Handling

- Proper HTTP status codes (200, 400, 404, 500)
- Graceful error recovery with panic handling
- Detailed error logging

## Architecture

### Core Components

#### Server

- Handles TCP connections and converts them to HTTP requests
- Manages connection timeouts and error recovery
- Dispatches requests to the appropriate handlers

#### Router

- Maps URL patterns to handler functions
- Supports parameter extraction from URLs
- Handles method-based routing (GET, POST)

#### Request/Response

- Custom implementation of HTTP request parsing
- Response writer with header management
- Buffered I/O for efficient network operations

### Configuration

- Environment variable support
- Command-line flag integration
- Configurable port settings

## Usage

### Starting the Server

```bash
# Start with default port (4221)
./run.sh

# Start with custom port
./run.sh -port=8080

# Or using environment variable
PORT=8080 ./run.sh
```

### Making Requests

```bash
# Health check
curl http://localhost:4221/v1/healthcheck

# Echo service
curl http://localhost:4221/v1/echo/hello-world

# User agent
curl http://localhost:4221/v1/user-agent
```

## Development

### Project Structure

- `app/main.go` - Entry point and application setup
- `app/server.go` - Core HTTP server implementation
- `app/routes.go` - Routing system implementation
- `app/feature_routes.go` - API endpoint handlers
- `app/header.go` - HTTP header and request parsing
- `app/env.go` - Configuration and environment handling
- `app/errors.go` - Error handling utilities

### Adding New Routes

To add a new route, implement a handler function and register it in the `Routes()` function in `routes.go`:

```go
// 1. Create a handler function
func (app *application) NewEndpoint(w ResponseWriter, r *Request) {
    w.Write(StatusOK, []byte("This is a new endpoint"))
}

// 2. Register the route in Routes()
func (app *application) Routes() HttpHandler {
    router := NewRouter()

    // Existing routes
    router.GET("/v1/healthcheck", app.Healthcheck)

    // New route
    router.GET("/v1/new-endpoint", app.NewEndpoint)

    return router
}
```

## Performance Considerations

- Uses buffered I/O for efficient network operations
- Implements connection timeouts to prevent resource exhaustion
- Handles each connection in a separate goroutine for concurrency

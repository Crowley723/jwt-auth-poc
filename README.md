# Go HTTP SQLite Server Template

A minimal Go HTTP server template with SQLite database integration, structured logging, graceful shutdown, and a custom middleware system for rapid project initialization.

## Features

- Clean middleware architecture with dependency injection
- SQLite database integration with migrations
- Structured logging using Go's `slog` package
- Graceful shutdown with signal handling
- Hot reload development workflow
- Built-in JSON and text response utilities
- Debug support with Delve debugger
- Database connection management and query utilities

## Quick Start

```bash
# Install dependencies
make install

# Start development server with hot reload
make dev

# Server runs at http://localhost:8080
curl http://localhost:8080/health
```

## Development

```bash
make dev         # Hot reload development
make debug       # Debug mode with delve debugger on :2345
make test        # Run tests
make coverage    # Run tests with coverage
make generate    # Run go generate
```

## Architecture

The template uses a custom `AppContext` middleware system that provides:

- Request/response access
- Structured logger
- JSON/text response utilities
- Error handling helpers
- SQLite database connection and query methods

Add new routes in `api/server.go`:
```go
mux.HandleFunc("GET /api/users", middlewares.Wrap(handleUsersGET))
```

Create handlers in `api/handlers.go`:
```go
func handleUsersGET(ctx *middlewares.AppContext) {
    // Database operations available via ctx.DB
    // users, err := ctx.GetUsers()
    ctx.WriteJSON(http.StatusOK, map[string]string{"users": "data"})
}
```

## Database

The template includes SQLite integration with:

- Database connection management in `db/db.go`
- Migration system with SQL files in `db/sqlite_migrations/`
- Generated query methods in `db/sqlite_queries.go`
- Database access via `AppContext.DB`

Migrations are automatically run on server startup. Add new migrations as numbered SQL files in `db/sqlite_migrations/`.

## Possible Enhancements

- Docker containerization
- Support for other database types (mysql or postgresql)
- Bidirectional migrations (allow migrating down to earlier schema version)

## Use as Template

This template provides a foundation for HTTP services requiring database persistence, authentication, APIs, or web applications. The middleware system makes it easy to add JWT auth, CORS, rate limiting, and other common HTTP service features. The SQLite integration provides immediate data persistence without external database dependencies.
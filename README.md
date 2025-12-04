# JWT Authentication Proof of Concept

A JWT-based authentication system with access tokens and refresh tokens, built in Go with SQLite.

**Warning: This is a proof of concept and not production ready.**

## Features

- User registration with Argon2 password hashing
- User login issuing access tokens (JWT) and refresh tokens
- JWT access token generation using ECDSA P-256 signing
- Refresh token storage with SHA256 hashing
- Token refresh endpoint to exchange refresh tokens for new access tokens
- Protected endpoints requiring JWT authentication
- JWKS endpoint for public key distribution
- User CRUD operations
- Healthcheck endpoint

## Technology Stack

- Go 1.25+
- SQLite database
- go-jose/v4 for JWT with ECDSA support
- go-crypt for Argon2 password hashing
- slog for structured logging

## API Endpoints

### Authentication
- `POST /api/login` - Authenticate user, returns access token and refresh token
- `POST /api/refresh` - Exchange refresh token for new access token

### Protected Endpoints (require JWT)
- `GET /api/protected/data` - Returns protected user data
- `GET /api/protected/stats` - Returns user statistics

### User Management
- `GET /api/users` - List all users (limit 100)
- `POST /api/users` - Create a new user
- `GET /api/users/{id}` - Get user by ID
- `DELETE /api/users/{id}` - Delete user by ID

### System
- `GET /health` - Health check endpoint
- `GET /api/jwks.json` - JSON Web Key Set for public key distribution

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id TEXT NOT NULL,
    hash TEXT NOT NULL,
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);
```

## Security Implementation

### Token Types
- **Access Tokens**: JWT tokens with 24 hour expiry for API authentication
- **Refresh Tokens**: Random tokens with 30 day expiry for obtaining new access tokens

### Cryptography
- **ECDSA P-256**: JWT signing algorithm
- **Argon2**: Password hashing (RFC 9106 low memory profile)
- **SHA-256**: Refresh token hashing before storage
- **crypto/rand**: Cryptographically secure random token generation

### Security Practices
- Passwords hashed with Argon2 before storage
- Refresh tokens hashed before database storage
- JWT signature verification on protected endpoints
- Token expiry validation
- Structured error responses

## Getting Started

### Prerequisites
- Go 1.25 or higher

### Installation and Running

```bash
# Clone the repository
git clone https://github.com/Crowley723/jwt-auth-poc.git
cd jwt-auth-poc

# Run the server
go run ./main.go
```

The server starts on `http://localhost:8080`. On first run:
- Generates ECDSA key pair (stored in `./app/certs/`)
- Creates SQLite database (stored in `./app/data/app.db`)
- Runs database migrations

### Example API Usage

**Create a user:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe","password":"password123"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6...",
  "refresh_token_expiry": 1701648000
}
```

**Access protected endpoint:**
```bash
curl http://localhost:8080/api/protected/data \
  -H "Authorization: Bearer <access_token>"
```

**Refresh access token:**
```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

## License

See [LICENSE](./LICENSE)
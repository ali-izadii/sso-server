# SSO Server

A minimal Single Sign-On (SSO) server written in Go for educational purposes. This project demonstrates the core concepts of SSO authentication including user management, session handling, and token-based authentication.

## Features

- User registration and authentication
- JWT-based access tokens
- Refresh token rotation
- Session management
- Client application registration
- OAuth 2.0 authorization code flow (simplified)

## Data Model

### User
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| email | string | Unique user email |
| password_hash | string | Bcrypt hashed password |
| name | string | User display name |
| is_active | bool | Account status |
| created_at | timestamp | Creation time |
| updated_at | timestamp | Last update time |

### Client (Application)
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| name | string | Application name |
| secret | string | Client secret (hashed) |
| redirect_uris | []string | Allowed redirect URIs |
| is_active | bool | Client status |
| created_at | timestamp | Creation time |

### Session
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to User |
| refresh_token | string | Hashed refresh token |
| user_agent | string | Client user agent |
| ip_address | string | Client IP |
| expires_at | timestamp | Session expiration |
| created_at | timestamp | Creation time |

### AuthorizationCode
| Field | Type | Description |
|-------|------|-------------|
| code | string | Authorization code |
| client_id | UUID | Foreign key to Client |
| user_id | UUID | Foreign key to User |
| redirect_uri | string | Redirect URI used |
| scope | string | Requested scope |
| expires_at | timestamp | Code expiration (short-lived) |
| created_at | timestamp | Creation time |

## API Endpoints

### Authentication

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}

Response: 201 Created
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}

Response: 200 OK
{
  "access_token": "jwt_token",
  "refresh_token": "refresh_token",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

#### Refresh Token
```
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token"
}

Response: 200 OK
{
  "access_token": "new_jwt_token",
  "refresh_token": "new_refresh_token",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

#### Logout
```
POST /api/v1/auth/logout
Authorization: Bearer <access_token>

Response: 204 No Content
```

### User Management

#### Get Current User
```
GET /api/v1/users/me
Authorization: Bearer <access_token>

Response: 200 OK
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe"
}
```

#### Update User
```
PATCH /api/v1/users/me
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Jane Doe"
}

Response: 200 OK
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "Jane Doe"
}
```

### OAuth 2.0 (Simplified)

#### Authorization Endpoint
```
GET /oauth/authorize?client_id=<client_id>&redirect_uri=<uri>&response_type=code&scope=<scope>&state=<state>

Response: Redirect to login page or consent page
```

#### Token Endpoint
```
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&code=<code>&redirect_uri=<uri>&client_id=<client_id>&client_secret=<secret>

Response: 200 OK
{
  "access_token": "jwt_token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "refresh_token"
}
```

### Client Management (Admin)

#### Register Client
```
POST /api/v1/clients
Authorization: Bearer <admin_access_token>
Content-Type: application/json

{
  "name": "My Application",
  "redirect_uris": ["https://myapp.com/callback"]
}

Response: 201 Created
{
  "id": "uuid",
  "name": "My Application",
  "secret": "generated_secret",
  "redirect_uris": ["https://myapp.com/callback"]
}
```

### Health Check

#### Health
```
GET /health

Response: 200 OK
{
  "status": "healthy"
}
```

## Project Structure

```
sso-server/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── config/
│   ├── config.local.yaml     # Local development config
│   ├── config.dev.yaml       # Development environment config
│   └── config.prod.yaml      # Production config
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management (Viper)
│   ├── handler/
│   │   ├── auth.go           # Authentication handlers
│   │   ├── user.go           # User handlers
│   │   ├── oauth.go          # OAuth handlers
│   │   └── client.go         # Client handlers
│   ├── middleware/
│   │   └── auth.go           # JWT authentication middleware
│   ├── model/
│   │   ├── user.go           # User model
│   │   ├── session.go        # Session model
│   │   ├── client.go         # Client model
│   │   └── auth_code.go      # Authorization code model
│   ├── repository/
│   │   ├── user.go           # User repository
│   │   ├── session.go        # Session repository
│   │   ├── client.go         # Client repository
│   │   └── auth_code.go      # Authorization code repository
│   ├── service/
│   │   ├── auth.go           # Authentication service
│   │   ├── user.go           # User service
│   │   ├── token.go          # Token service
│   │   └── oauth.go          # OAuth service
│   └── database/
│       └── database.go       # Database connection
├── pkg/
│   └── validator/
│       └── validator.go      # Input validation helpers
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## Configuration

Configuration is managed using [Viper](https://github.com/spf13/viper) with YAML config files.

### Config Files

| File | Environment | Description |
|------|-------------|-------------|
| `config/config.local.yaml` | local | Local development with SQLite |
| `config/config.dev.yaml` | dev | Development with PostgreSQL |
| `config/config.prod.yaml` | prod | Production settings |

### Environment Selection

Set `APP_ENV` to select the config file:

```bash
APP_ENV=local   # loads config.local.yaml (default)
APP_ENV=dev     # loads config.dev.yaml
APP_ENV=prod    # loads config.prod.yaml
```

### Config Structure

```yaml
server:
  port: 8080
  host: localhost

database:
  driver: sqlite          # sqlite or postgres
  dsn: sso.db             # connection string

jwt:
  secret: your-secret-key # required
  expiry: 1h              # access token expiry
  refresh_expiry: 168h    # refresh token expiry (7 days)

oauth:
  auth_code_expiry: 10m   # authorization code expiry

log:
  level: debug            # debug, info, warn, error
  format: text            # text or json
```

### Environment Variable Override

Environment variables override config file values. Use underscore-separated uppercase names:

| Variable | Config Key |
|----------|------------|
| `SERVER_PORT` | `server.port` |
| `DATABASE_DSN` | `database.dsn` |
| `JWT_SECRET` | `jwt.secret` |

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Running the Server

```bash
# Run with local config (default)
go run cmd/server/main.go

# Run with specific environment
APP_ENV=dev go run cmd/server/main.go

# Override config with environment variables
JWT_SECRET=my-secret APP_ENV=prod go run cmd/server/main.go

# Build and run
go build -o sso-server cmd/server/main.go
./sso-server
```

### Running Tests

```bash
go test ./...
```

## Security Considerations

This is an educational project. For production use, consider:

- Use HTTPS only
- Implement rate limiting
- Add CSRF protection
- Use secure cookie settings
- Implement proper password policies
- Add audit logging
- Use a production-grade database
- Implement proper key rotation
- Add multi-factor authentication

## License

MIT License - Feel free to use for learning purposes.

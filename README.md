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
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
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

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| DATABASE_URL | Database connection string | sqlite://sso.db |
| JWT_SECRET | Secret key for JWT signing | (required) |
| JWT_EXPIRY | Access token expiry duration | 1h |
| REFRESH_TOKEN_EXPIRY | Refresh token expiry duration | 7d |
| AUTH_CODE_EXPIRY | Authorization code expiry | 10m |

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Running the Server

```bash
# Set required environment variables
export JWT_SECRET="your-secret-key"

# Run the server
go run cmd/server/main.go

# Or build and run
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

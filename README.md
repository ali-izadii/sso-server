# SSO Server - Single Sign-On Authentication Service

A lightweight, production-ready Single Sign-On (SSO) server built with Go, PostgreSQL, and Redis following clean architecture principles.

## 🚀 Features

### MVP Features (Phase 1)
- **User Authentication**: Registration, login, logout
- **OAuth 2.0 Flow**: Authorization code flow for client applications
- **JWT Token Management**: Access and refresh token generation/validation
- **Session Management**: Secure session handling with Redis
- **Client Application Management**: Register and manage OAuth clients
- **Basic User Profile**: User information endpoints

### Future Features (Phase 2+)
- **SAML 2.0 Support**: Enterprise SSO integration
- **Multi-Factor Authentication (MFA)**: TOTP, SMS, Email verification
- **Social Login**: Google, GitHub, Microsoft OAuth
- **Role-Based Access Control (RBAC)**: User roles and permissions
- **Admin Dashboard**: Web UI for user and application management
- **Audit Logging**: Security and access logs
- **Rate Limiting**: API protection and abuse prevention
- **Email Templates**: Customizable notification emails

## 🏗️ Architecture

### Project Structure
```
sso-server/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── domain/                  # Domain entities and interfaces
│   │   ├── models/              # Data models
│   │   └── repositories/        # Repository interfaces
│   ├── infrastructure/          # External dependencies
│   │   ├── database/            # PostgreSQL connection
│   │   ├── redis/               # Redis connection
│   │   └── migrations/          # Database migrations
│   ├── services/                # Business logic
│   │   ├── auth/                # Authentication service
│   │   ├── oauth/               # OAuth 2.0 service
│   │   └── user/                # User management service
│   ├── handlers/                # HTTP handlers (controllers)
│   │   ├── auth/                # Auth endpoints
│   │   ├── oauth/               # OAuth endpoints
│   │   └── user/                # User endpoints
│   ├── middleware/              # HTTP middleware
│   └── utils/                   # Utility functions
├── pkg/                         # Public packages
│   ├── jwt/                     # JWT utilities
│   ├── password/                # Password hashing
│   └── validator/               # Input validation
├── migrations/                  # Database migration files
├── docker-compose.yml           # Local development setup
├── Dockerfile                   # Container image
├── .env.example                 # Environment variables template
└── README.md
```

### Technology Stack
- **Backend**: Go 1.21+ with Gin framework
- **Database**: PostgreSQL 15+ for persistent data
- **Cache**: Redis 7+ for sessions and temporary data
- **Authentication**: JWT tokens with RSA/HMAC signing
- **Migration**: golang-migrate for database versioning
- **Containerization**: Docker and Docker Compose

## 📊 Data Schema

### Core Entities

#### 1. Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE
);
```

#### 2. Applications Table (OAuth Clients)
```sql
CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    client_id VARCHAR(255) UNIQUE NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    redirect_uris TEXT[] NOT NULL,  -- Array of allowed redirect URIs
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### 3. Authorization Codes Table
```sql
CREATE TABLE authorization_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    redirect_uri VARCHAR(500) NOT NULL,
    scopes VARCHAR(500) DEFAULT 'openid profile email',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### 4. Access Tokens Table
```sql
CREATE TABLE access_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(500) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    scopes VARCHAR(500) DEFAULT 'openid profile email',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### 5. Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(500) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    access_token_id UUID REFERENCES access_tokens(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### 6. User Sessions Table
```sql
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Database Indexes
```sql
-- Performance indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_active ON users(is_active);
CREATE INDEX idx_applications_client_id ON applications(client_id);
CREATE INDEX idx_authorization_codes_code ON authorization_codes(code);
CREATE INDEX idx_authorization_codes_expires ON authorization_codes(expires_at);
CREATE INDEX idx_access_tokens_token ON access_tokens(token);
CREATE INDEX idx_access_tokens_user_app ON access_tokens(user_id, application_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX idx_user_sessions_user_expires ON user_sessions(user_id, expires_at);
```

### Key Relationships
- **Users ↔ Applications**: Many-to-many through authorization codes and tokens
- **Users → Sessions**: One-to-many (user can have multiple active sessions)
- **Users → Tokens**: One-to-many (user can have tokens for multiple applications)
- **Applications → Tokens**: One-to-many (application can have tokens for multiple users)
- **Access Tokens ↔ Refresh Tokens**: One-to-one relationship for token pairs

### Data Flow
1. **User Registration/Login** → `users` table
2. **OAuth Authorization** → `authorization_codes` table (temporary)
3. **Token Exchange** → `access_tokens` + `refresh_tokens` tables
4. **Session Management** → `user_sessions` table
5. **Application Management** → `applications` table

## 🔌 API Endpoints

### Authentication Endpoints
```
POST   /auth/register           # User registration
POST   /auth/login              # User login
POST   /auth/logout             # User logout
POST   /auth/refresh            # Refresh access token
GET    /auth/verify             # Verify token validity
```

### OAuth 2.0 Endpoints
```
GET    /oauth/authorize         # OAuth authorization endpoint
POST   /oauth/token             # Token exchange endpoint
GET    /oauth/userinfo          # User information endpoint
POST   /oauth/revoke            # Token revocation
```

### User Management
```
GET    /users/profile           # Get user profile
PUT    /users/profile           # Update user profile
POST   /users/change-password   # Change password
```

### Application Management
```
GET    /applications            # List user's applications
POST   /applications            # Register new application
PUT    /applications/:id        # Update application
DELETE /applications/:id        # Delete application
```

## 🛠️ Setup & Installation

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd sso-server
   ```

2. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start dependencies with Docker**
   ```bash
   docker-compose up -d postgres redis
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

### Docker Setup
```bash
docker-compose up -d
```

## ⚙️ Configuration

### Environment Variables
```bash
# Server Configuration
PORT=8080
ENV=development

# Database Configuration  
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=sso_db
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# OAuth Configuration
OAUTH_CODE_EXPIRY=10m
```

## 🔐 Security Features

- **Password Security**: bcrypt hashing with salt
- **JWT Security**: RSA/HMAC signing with configurable expiry
- **Session Security**: Redis-based session management
- **CORS Protection**: Configurable cross-origin policies
- **Rate Limiting**: Per-IP and per-user rate limits
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries

## 📈 Monitoring & Observability

### Health Checks
```
GET /health              # Application health
GET /health/db           # Database connectivity
GET /health/redis        # Redis connectivity
```

### Metrics (Future)
- Prometheus metrics endpoint
- Request duration and count
- Authentication success/failure rates
- Token generation and validation metrics

## 🧪 Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run with coverage
make test-coverage

# Load testing
make load-test
```

## 📝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Support

- **Issues**: GitHub Issues for bug reports and feature requests
- **Documentation**: Wiki pages for detailed guides
- **Community**: Discussions for questions and support

---

**Note**: This is an MVP implementation. For production use, consider additional security hardening, monitoring, and compliance requirements based on your specific needs.

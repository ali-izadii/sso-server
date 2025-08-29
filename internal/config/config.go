package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	OAuth      OAuthConfig
	Email      EmailConfig
	RateLimit  RateLimitConfig
	CORS       CORSConfig
	Security   SecurityConfig
	Logging    LoggingConfig
	Monitoring MonitoringConfig
	Cache      CacheConfig
}

type ServerConfig struct {
	Port         string
	Env          string
	Debug        bool
	HotReload    bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Algorithm          string
}

type OAuthConfig struct {
	AuthCodeExpiry time.Duration
}

type EmailConfig struct {
	SMTPHost    string
	SMTPPort    int
	Username    string
	Password    string
	FromEmail   string
	FromName    string
	TemplateDir string
	UseTLS      bool
}

type RateLimitConfig struct {
	Enabled  bool
	Requests int
	Window   time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposeHeaders  []string
	MaxAge         time.Duration
}

type SecurityConfig struct {
	HeadersEnabled     bool
	HSTSMaxAge         int
	ContentTypeNoSniff bool
	FrameDeny          bool
	XSSProtection      bool
}

type LoggingConfig struct {
	Level    string
	Format   string
	Output   string
	Filename string
}

type MonitoringConfig struct {
	MetricsEnabled     bool
	MetricsPath        string
	HealthCheckEnabled bool
	HealthCheckPath    string
}

type CacheConfig struct {
	TTL       time.Duration
	MaxMemory string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Determine which .env file to load based on ENV
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	// Load appropriate .env file
	switch env {
	case "production":
		_ = godotenv.Load(".env.production", ".env")
	case "docker":
		_ = godotenv.Load(".env.docker", ".env")
	default:
		_ = godotenv.Load(".env.local", ".env")
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Env:          env,
			Debug:        getBoolEnv("DEBUG_ENABLED", env == "development"),
			HotReload:    getBoolEnv("HOT_RELOAD", env == "development"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", "30s"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "sso_db"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", "3600s"),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", generateWarningSecret()),
			AccessTokenExpiry:  getDurationEnv("JWT_ACCESS_EXPIRY", "15m"),
			RefreshTokenExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", "168h"),
			Algorithm:          getEnv("JWT_ALGORITHM", "HS256"),
		},
		OAuth: OAuthConfig{
			AuthCodeExpiry: getDurationEnv("OAUTH_CODE_EXPIRY", "10m"),
		},
		Email: EmailConfig{
			SMTPHost:    getEnv("SMTP_HOST", "localhost"),
			SMTPPort:    getEnvAsInt("SMTP_PORT", 587),
			Username:    getEnv("SMTP_USERNAME", ""),
			Password:    getEnv("SMTP_PASSWORD", ""),
			FromEmail:   getEnv("FROM_EMAIL", "noreply@localhost"),
			FromName:    getEnv("FROM_NAME", "SSO Server"),
			TemplateDir: getEnv("EMAIL_TEMPLATE_DIR", "templates/email"),
			UseTLS:      getBoolEnv("SMTP_USE_TLS", true),
		},
		RateLimit: RateLimitConfig{
			Enabled:  getBoolEnv("RATE_LIMIT_ENABLED", env == "production"),
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", 60),
			Window:   getDurationEnv("RATE_LIMIT_WINDOW", "60s"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getSliceEnv("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getSliceEnv("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
			ExposeHeaders:  getSliceEnv("CORS_EXPOSE_HEADERS", []string{}),
			MaxAge:         getDurationEnv("CORS_MAX_AGE", "12h"),
		},
		Security: SecurityConfig{
			HeadersEnabled:     getBoolEnv("SECURITY_HEADERS_ENABLED", env == "production"),
			HSTSMaxAge:         getEnvAsInt("HSTS_MAX_AGE", 31536000),
			ContentTypeNoSniff: getBoolEnv("CONTENT_TYPE_NOSNIFF", true),
			FrameDeny:          getBoolEnv("FRAME_DENY", true),
			XSSProtection:      getBoolEnv("XSS_PROTECTION", true),
		},
		Logging: LoggingConfig{
			Level:    getEnv("LOG_LEVEL", getDefaultLogLevel(env)),
			Format:   getEnv("LOG_FORMAT", getDefaultLogFormat(env)),
			Output:   getEnv("LOG_OUTPUT", "stdout"),
			Filename: getEnv("LOG_FILENAME", ""),
		},
		Monitoring: MonitoringConfig{
			MetricsEnabled:     getBoolEnv("METRICS_ENABLED", env == "production"),
			MetricsPath:        getEnv("METRICS_PATH", "/metrics"),
			HealthCheckEnabled: getBoolEnv("HEALTH_CHECK_ENABLED", true),
			HealthCheckPath:    getEnv("HEALTH_CHECK_PATH", "/health"),
		},
		Cache: CacheConfig{
			TTL:       getDurationEnv("CACHE_TTL", "3600s"),
			MaxMemory: getEnv("CACHE_MAX_MEMORY", "256mb"),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s, using default", key)
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Warning: Invalid boolean value for %s, using default", key)
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Warning: Invalid duration value for %s, using default", key)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getDefaultLogLevel(env string) string {
	if env == "production" {
		return "info"
	}
	return "debug"
}

func getDefaultLogFormat(env string) string {
	if env == "production" {
		return "json"
	}
	return "text"
}

func generateWarningSecret() string {
	log.Println("WARNING: Using default JWT secret. This is NOT secure for production!")
	return "insecure-default-secret-change-me-in-production"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Env == "production" {
		if c.JWT.Secret == generateWarningSecret() || len(c.JWT.Secret) < 32 {
			log.Fatal("JWT secret must be set and at least 32 characters long in production")
		}
		if c.Database.SSLMode == "disable" {
			log.Println("WARNING: SSL is disabled for database connection in production")
		}
	}
	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

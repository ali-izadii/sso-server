package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sso-server/internal/infrastracture/database"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"sso-server/internal/config"
	"sso-server/internal/handlers/user"
	userService "sso-server/internal/services/user"
)

func main() {
	myFigure := figure.NewColorFigure("SSO Server", "", "green", true)
	myFigure.Print()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := initRedis(cfg)
	defer redisClient.Close()

	// Initialize repositories
	userRepo := database.NewUserRepository(db)

	// Initialize services
	userSvc := userService.NewService(userRepo)

	// Initialize handlers
	userHandler := user.NewHandler(userSvc)

	// Initialize router
	router := setupRouter(cfg)

	// Register routes
	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "sso-server",
		})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting SSO server on port %s", cfg.Server.Port)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDatabase initializes PostgreSQL connection
func initDatabase(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Println("Database connection established")
	return db, nil
}

// initRedis initializes Redis connection
func initRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	log.Println("Redis connection established")
	return client
}

// setupRouter configures the Gin router
func setupRouter(cfg *config.Config) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware (basic implementation)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	return router
}

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vannyezha/sis-inv/internal/config"
	"github.com/vannyezha/sis-inv/internal/database"
	"github.com/vannyezha/sis-inv/internal/handlers"
	"github.com/vannyezha/sis-inv/internal/middleware"
	"github.com/vannyezha/sis-inv/internal/utils"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		// Try loading from current directory
		_ = godotenv.Load(".env")
	}

	// Load config
	cfg := config.Load()

	// Initialize JWT
	utils.InitJWT(cfg.JWTSecret)

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Seed default data
	if err := database.SeedDefaultAdmin(db); err != nil {
		log.Fatalf("❌ Failed to seed data: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	systemHandler := handlers.NewSystemHandler(db)

	// Setup router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// System (public)
		v1.GET("/system/health", systemHandler.HealthCheck)

		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.Me)
			auth.PUT("/password", middleware.AuthMiddleware(), authHandler.ChangePassword)
		}

		// User Management (Admin only)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("ADMIN"))
		{
			users.GET("", userHandler.List)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 SIS-INV API running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
		os.Exit(1)
	}
}

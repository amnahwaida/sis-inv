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

	// Static files for uploads
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		_ = os.Mkdir("uploads", 0755)
	}
	r.Static("/uploads", "./uploads")

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// System (public)
		v1.GET("/system/health", systemHandler.HealthCheck)

		// Upload (Auth only)
		uploadHandler := handlers.NewUploadHandler()
		v1.POST("/upload", middleware.AuthMiddleware(), uploadHandler.UploadFile)

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

		// Students
		studentHandler := handlers.NewStudentHandler(db)
		students := v1.Group("/students")
		students.Use(middleware.AuthMiddleware())
		{
			students.GET("/search", studentHandler.Search)
			students.GET("/classes", studentHandler.GetUniqueClasses)
			
			// Admin only CRUD
			adminStudents := students.Group("")
			adminStudents.Use(middleware.RoleMiddleware("ADMIN"))
			{
				adminStudents.GET("", studentHandler.List)
				adminStudents.POST("", studentHandler.Create)
				adminStudents.PUT("/:id", studentHandler.Update)
				adminStudents.DELETE("/:id", studentHandler.Delete)
				adminStudents.GET("/export", studentHandler.ExportExcel)
				adminStudents.POST("/import", studentHandler.ImportExcel)
			}
		}

		// Items Management
		itemHandler := handlers.NewItemHandler(db)
		items := v1.Group("/items")
		items.Use(middleware.AuthMiddleware())
		{
			// All logged in users can view items
			items.GET("", itemHandler.List)
			items.GET("/:id", itemHandler.Get)
			items.GET("/code/:code", itemHandler.GetByCode)
			items.GET("/:id/history", itemHandler.GetHistory)
			
			// Only Admins can modify items
			adminItems := items.Group("")
			adminItems.Use(middleware.RoleMiddleware("ADMIN"))
			{
				adminItems.POST("", itemHandler.Create)
				adminItems.PUT("/:id", itemHandler.Update)
				adminItems.DELETE("/:id", itemHandler.Delete)
				adminItems.POST("/:id/qr", itemHandler.GenerateQRCode)
				adminItems.POST("/import", itemHandler.ImportExcel)
			}
		}

		// Transactions Management
		trxHandler := handlers.NewTransactionHandler(db)
		transactions := v1.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware())
		{
			// All logged in staff can borrow/return and view their own
			transactions.POST("/borrow", trxHandler.Borrow)
			transactions.POST("/return", trxHandler.Return)
			transactions.GET("/my", trxHandler.MyBorrowings)
			transactions.GET("/my/history", trxHandler.MyBorrowingsHistory)
			transactions.GET("/student/:nis", trxHandler.GetStudentHistory)
			transactions.GET("/student/:nis/export", trxHandler.StudentHistoryExport)
		}

		// Dashboard
		dashHandler := handlers.NewDashboardHandler(db)
		dashboard := v1.Group("/dashboard")
		// In a real app we might cache this or separate per role, but here we just require auth
		dashboard.Use(middleware.AuthMiddleware())
		{
			dashboard.GET("/summary", dashHandler.GetSummary)
		}

		// Categories Management
		catHandler := handlers.NewCategoryHandler(db)
		categories := v1.Group("/categories")
		categories.Use(middleware.AuthMiddleware())
		{
			categories.GET("", catHandler.List)
			// Only Admins can modify categories
			adminCats := categories.Group("")
			adminCats.Use(middleware.RoleMiddleware("ADMIN"))
			{
				adminCats.POST("", catHandler.Create)
				adminCats.PUT("/:id", catHandler.Update)
				adminCats.DELETE("/:id", catHandler.Delete)
			}
		}

		// Locations Management
		locHandler := handlers.NewLocationHandler(db)
		locations := v1.Group("/locations")
		locations.Use(middleware.AuthMiddleware())
		{
			locations.GET("", locHandler.List)
			// Only Admins can modify locations
			adminLocs := locations.Group("")
			adminLocs.Use(middleware.RoleMiddleware("ADMIN"))
			{
				adminLocs.POST("", locHandler.Create)
				adminLocs.PUT("/:id", locHandler.Update)
				adminLocs.DELETE("/:id", locHandler.Delete)
			}
		}
		// Reports & Audit
		reportHandler := handlers.NewReportHandler(db)
		auditHandler := handlers.NewAuditHandler(db)
		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware())
		{
			reports.GET("/export/items", reportHandler.ExportItems)
			reports.GET("/export/transactions", reportHandler.ExportTransactions)
			reports.GET("/export/audit", reportHandler.ExportAuditLogs)
			reports.GET("/export/maintenance", reportHandler.ExportMaintenanceLogs)
			reports.GET("/active-borrowings", reportHandler.ActiveBorrowings)
			reports.GET("/overdue", reportHandler.OverdueReport)
			reports.GET("/history", reportHandler.TransactionHistory)
			reports.GET("/audit", reportHandler.AuditLogs)

			// Stock Opname / Audit Sessions
			audit := reports.Group("/audit-sessions")
			{
				audit.GET("", auditHandler.ListSessions)
				audit.POST("", auditHandler.StartSession)
				audit.GET("/:id", auditHandler.GetSessionDetail)
				audit.GET("/:id/export", reportHandler.ExportAuditSession)
				audit.POST("/:id/scan", auditHandler.ScanItem)
				audit.POST("/:id/close", auditHandler.CloseSession)
			}
		}

		// Maintenance / Log Perbaikan
		maintHandler := handlers.NewMaintenanceHandler(db)
		maintenance := v1.Group("/maintenance")
		maintenance.Use(middleware.AuthMiddleware())
		{
			maintenance.GET("", maintHandler.List)
			maintenance.POST("", maintHandler.Create)
			maintenance.PUT("/:id", maintHandler.UpdateStatus)
			maintenance.DELETE("/:id", maintHandler.Delete)
		}

		// Settings (Admin only for PUT, All Auth for GET)
		settingHandler := handlers.NewSettingHandler(db)
		settings := v1.Group("/settings")
		settings.Use(middleware.AuthMiddleware())
		{
			settings.GET("", settingHandler.List)
			settings.PUT("", middleware.RoleMiddleware("ADMIN"), settingHandler.Update)
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

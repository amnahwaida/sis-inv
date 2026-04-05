package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/vannyezha/sis-inv/internal/utils"
)

func main() {
	// Load .env
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" { dbUser = "sekolah_admin" }
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" { dbPass = "sekolah_secret_2024" }
	dbName := os.Getenv("DB_NAME")
	if dbName == "" { dbName = "inventaris_db" }
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" { dbHost = "localhost" }
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" { dbPort = "5432" }

	adminUser := os.Getenv("ADMIN_DEFAULT_USERNAME")
	if adminUser == "" { adminUser = "admin" }
	adminPass := os.Getenv("ADMIN_DEFAULT_PASSWORD")
	if adminPass == "" { adminPass = "password123" }

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Hash the password
	hash, err := utils.HashPassword(adminPass)
	if err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}

	// Update user
	res, err := db.Exec(ctx, 
		"UPDATE users SET role = 'ADMIN', password_hash = $1, is_active = true WHERE username = $2", 
		hash, adminUser)
	
	if err != nil {
		log.Fatalf("❌ Failed to update admin user: %v", err)
	}

	if res.RowsAffected() == 0 {
		fmt.Printf("⚠️  User '%s' not found. Creating new admin user...\n", adminUser)
		_, err = db.Exec(ctx, 
			"INSERT INTO users (username, full_name, role, password_hash, is_active) VALUES ($1, $2, 'ADMIN', $3, true)",
			adminUser, "System Admin", hash)
		if err != nil {
			log.Fatalf("❌ Failed to create admin user: %v", err)
		}
		fmt.Printf("✅ Admin user '%s' created as defined in .env.\n", adminUser)
	} else {
		fmt.Printf("✅ Admin user '%s' reset to role ADMIN and credentials from .env.\n", adminUser)
	}

	fmt.Println("\n[CREDENTIALS USED]")
	fmt.Println("Username:", adminUser)
	fmt.Println("Password:", "(hidden - see .env file)")
}
